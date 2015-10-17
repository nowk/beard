package beard

import (
	"bytes"
	"io"
)

type Renderable struct {
	// File is the template file to be rendered. File must be explicitly closed
	// by the user
	File io.ReadSeeker

	// Data is the data set to be used when compiling variables into the html
	Data *Data

	// del represents the current delimiter
	del Delim

	// buf holds bytes that end with a partially matched delimiter. It will
	// attempt an exact match on the next read
	buf []byte

	// truncd holds bytes that were truncated due to length of p
	truncd []byte

	// eof marks the File has reached EOF.
	eof bool

	// cursor is the location at which the reader is at
	cursor int

	// blocks are added in a FILO order. The last block in the list would be the
	// current block
	blocks []*block
}

var _ io.Reader = &Renderable{}

func (r *Renderable) Read(p []byte) (int, error) {
	lenp := len(p)
	writ := 0

	// flush trucncated out to write
	if lent := len(r.truncd); lent > 0 {
		n, tr := r.flush(p)
		r.truncd = tr
		// return if we've written out to the length of p. NOTE this should
		// never write more than lenp
		if n >= lenp {
			return n, nil
		}

		writ = n
	}

	// share the given []byte argument so we don't have to allocate a temp
	// buffer on each Read call.

	n, err := r.File.Read(p[writ:])
	if err != nil {
		if err != io.EOF {
			return n, err
		}
		if !r.eof {
			r.eof = (err == io.EOF)
		}
	}

	// alloc buf with a cap of n
	// we do need a separate allocated block here.
	if r.buf == nil {
		r.buf = make([]byte, 0, n)
	}
	r.buf = append(r.buf, p[writ:writ+n]...)

	// trim p, so we can start from it's last written point
	p = p[:writ]

	switch b, ma := r.delim().Match(r.buf); ma {
	case paMatch:
		r.buf = b

	case exMatch:
		var (
			lenb   = len(b)
			lentag = lenb - len(r.delim().Value())

			del = b[lentag:]
			tag = b[:lentag]
		)

		r.buf = r.buf[lenb:]
		r.cursor += lenb
		r.swapDelim() // swap delim early, blocks will return early

		var val []byte

		// when we find a matching rdelim, {{..}} has been closed and we can now
		// parse for the var value
		if bytes.Equal(del, rdelim.Value()) {
			val, err = r.handleVar(tag)
			if err != nil {
				return writ, err
			}
			if val == nil {
				return writ, nil
			}
		} else {
			val = tag
		}

		// combine truncated with current value and write
		val = append(r.truncd, val...)

		n := len(val)
		if n > lenp {
			n = lenp - writ
			r.truncd = val[n:]
		} else {
			n -= writ
			r.truncd = r.truncd[:0]
		}

		p = append(p, val[:n]...)
		writ += n

	default:
		// if we have a buf, flush it. NOTE: buf at this point will always fit
		// into p
		if n := len(r.buf); n > 0 {
			p = append(p, r.buf[:n]...)
			writ += n

			r.buf = r.buf[:0]
			r.cursor += writ
		}
	}

	n = writ
	if r.eof && len(r.buf) == 0 && len(r.truncd) == 0 {
		return n, io.EOF
	}

	return n, nil
}

// delim returns the "current" delim in the Renderable, it defaults to ldelim.
func (r *Renderable) delim() Delim {
	if r.del == nil {
		return ldelim
	}

	return r.del
}

// swapDelim swaps the delim on a Renderable.
func (r *Renderable) swapDelim() {
	if bytes.Equal(r.delim().Value(), ldelim.Value()) {
		r.del = rdelim
	} else {
		r.del = ldelim
	}
}

func (r *Renderable) handleVar(v []byte) ([]byte, error) {
	tag := string(bytes.TrimSpace(v))
	if len(tag) == 0 {
		// TODO handle if tag is empty
	}

	esc := true

	switch tag[0] {
	case '#':
		bl := r.newBlock(tag, r.cursor)
		// bl := r.newBlock(tag, r.cursor-(len(v)+len(rdelim)+len(ldelim)))
		if bl == nil {
			// TODO not sure how this can actually happen...?
		}

		return v[:0], nil

	case '/':
		_, bl := r.currentBlock()
		if bl == nil {
			// TODO error: invalid block
		}
		if bl.tag != tag {
			// TODO error: non-matching block
		}
		if bl.increment(); bl.isFinished() {
			r.popBlock()

			return nil, nil
		}

		// reset the buffer and move the cursor to the block's cursor
		// location
		r.buf = r.buf[:0]
		r.cursor = bl.cursor

		// set the File's cusror to be read at on the next Read
		_, err := r.File.Seek(int64(r.cursor), 0)

		return nil, err

	case '&':
		tag = tag[1:]
		esc = false
	}
	// TODO how to handle/detect unclosed blocks

	val := r.getValue(tag)
	if esc {
		val = escapeBytes(val)
	}

	return val, nil
}

func (r *Renderable) newBlock(tag string, c int) *block {
	bl := r.findBlock(tag, c)
	if bl != nil {
		return bl
	}

	d := r.Data.Get(tag[1:])
	if d == nil {
		// TODO handle
	}
	bl = newBlock(tag, c, d)

	// lazy alloc
	if r.blocks == nil {
		r.blocks = make([]*block, 0, 32)
	}

	r.blocks = append(r.blocks, bl)

	return bl
}

// findBlock finds a block by it's name and cursor.
// The addition of the cursor provides a method of assigning uniqueness to a
// block, allowing blocks to nest the same block and have a fresh reference to
// the underlying data.
//
// The name provided should not containa any block prefixes,
// eg. #words -> words.
func (r *Renderable) findBlock(tag string, c int) *block {
	z := len(r.blocks) - 1
	if z < 0 {
		return nil
	}

	// look up block in reverse
	for ; z > -1; z-- {
		bl := r.blocks[z]
		if bl.cursor == c && bl.tag == tag {
			return bl
		}
	}

	return nil
}

// currentBlock returns the last block (and it's index) on the block list,
// which represents the current block.
func (r *Renderable) currentBlock() (int, *block) {
	z := len(r.blocks) - 1
	if z < 0 {
		return -1, nil
	}

	return z, r.blocks[z]
}

// popBlock pops off the last block in the blocks list
func (r *Renderable) popBlock() *block {
	i, bl := r.currentBlock()
	if i < 0 {
		return nil
	}
	if bl == nil {
		return nil
	}

	r.blocks = r.blocks[:i]

	return bl
}

// getValue looks up the value within Data map. It will iterrate *up* the blocks
// before looking at the root Data field itself.
func (r *Renderable) getValue(k string) []byte {
	z := len(r.blocks)
	for ; z > 0; z-- {
		bl := r.blocks[z-1]

		if v := bl.Data().Get(k); v != nil {
			return v.Bytes()
		}
	}

	// . never looks up outside of a block
	if k == "." {
		return []byte{}
	}
	if v := r.Data.Get(k); v != nil {
		return v.Bytes()
	}

	return nil
}

// flush writes truncated out to p. It writes up to the lesser of the two
// lengths, p vs truncd. It returns any remaing bytes that couldn't be written
// due to length constraints.
func (r *Renderable) flush(p []byte) (int, []byte) {
	lent := len(r.truncd)
	z := len(p)
	if lent < z {
		z = lent
	}

	i := 0
	for ; i < z; i++ {
		p[i] = r.truncd[i]
	}
	if i >= lent {
		return i, r.truncd[:0]
	}

	return i, r.truncd[i:]
}
