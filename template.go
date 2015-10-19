package beard

import (
	"bytes"
	"errors"
	"io"
)

// File is a simplifed interface of io.ReadSeeker
type File interface {
	Read([]byte) (int, error)
	Seek(int64, int) (int64, error)
}

// PartialFunc represents the func signature used when requesting a partial to
// be rendered. The interface{} to be returned must either assert to a File or
// *Template.
type PartialFunc func(string) (io.Reader, error)

type Template struct {
	// File is the template file to be rendered. File must be explicitly closed
	// by the user
	File File

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

	// partialFunc is a user definable func to handle how partials should be
	// handled. Partials will inherit this function, unless the PartialFunc
	// returns a *Template
	partialFunc PartialFunc

	// partial holds the reference to the current partial requiring rendering.
	// This will be nil'd when the partial has been completely rendered.
	partial io.Reader
}

var _ io.Reader = &Template{}

func (t *Template) Read(p []byte) (int, error) {
	lenp := len(p)
	writ := 0

	if t.partial != nil {
		n, err := t.readPartial(p)
		if err == nil {
			return n, nil
		}
		if err != io.EOF {
			return n, err
		}

		writ = n
	}

	// flush trucncated out to write
	if lent := len(t.truncd); lent > 0 {
		n, tr := t.flush(p)
		t.truncd = tr
		// return if we've written out to the length of p. NOTE this should
		// never write more than lenp
		if n >= lenp {
			return n, nil
		}

		writ = n
	}

	// share the given []byte argument so we don't have to allocate a temp
	// buffer on each Read call.

	n, err := t.File.Read(p[writ:])
	if err != nil {
		if err != io.EOF {
			return n, err
		}
		if !t.eof {
			t.eof = (err == io.EOF)
		}
	}

	// alloc buf with a cap of n
	// we do need a separate allocated block here.
	if t.buf == nil {
		t.buf = make([]byte, 0, n)
	}
	t.buf = append(t.buf, p[writ:writ+n]...)

	// trim p, so we can start from it's last written point
	p = p[:writ]

	switch b, ma := t.delim().Match(t.buf); ma {
	case paMatch:
		t.buf = b

	case exMatch:
		var (
			lenb   = len(b)
			lentag = lenb - len(t.delim().Value())

			del = b[lentag:]
			tag = b[:lentag]
		)

		t.buf = t.buf[lenb:]
		t.cursor += lenb
		t.swapDelim() // swap delim early, blocks will return early

		var val []byte

		// when we find a matching rdelim, {{..}} has been closed and we can now
		// parse for the var value
		if bytes.Equal(del, rdelim.Value()) {
			val, err = t.handleVar(tag)
			if err != nil {
				return writ, err
			}
			if val == nil {
				return writ, nil
			}
		} else {
			val = tag
		}

		// if we are in a current block and there is no data to render dont
		// render any of the inner block content
		_, bl := t.currentBlock()
		if bl != nil && bl.Skip() {
			return writ, nil
		}

		// combine truncated with current value and write
		val = append(t.truncd, val...)
		n := len(val)

		// amount to be written must fit in with in the available space that has
		// yet to be written on p. If we have more to write, truncate for next
		// Read
		if availn := lenp - writ; n > availn {
			n = availn
			t.truncd = val[n:]
		} else {
			t.truncd = t.truncd[:0]
		}

		p = append(p, val[:n]...)
		writ += n

	default:
		// if we have a buf, flush it. NOTE: buf at this point will always fit
		// into p
		if n := len(t.buf); n > 0 {
			p = append(p, t.buf[:n]...)
			writ += n

			t.buf = t.buf[:0]
			t.cursor += writ
		}
	}

	n = writ
	if t.eof && len(t.buf) == 0 && len(t.truncd) == 0 {
		// check for unclosed blocks
		var err = io.EOF
		if len(t.blocks) > 0 {
			err = errUnclosedBlocks
		}

		return n, err
	}

	return n, nil
}

func (t *Template) readPartial(p []byte) (int, error) {
	n, err := t.partial.Read(p)
	if err == nil {
		return n, nil
	}

	// we are in charge of explicitly closing, if we get a closer
	defer closePartial(t.partial)

	if err != io.EOF {
		return n, err
	}

	t.partial = nil

	return n, io.EOF

}

// delim returns the "current" delim in the Template, it defaults to ldelim.
func (t *Template) delim() Delim {
	if t.del == nil {
		return ldelim
	}

	return t.del
}

// swapDelim swaps the delim on a Template.
func (t *Template) swapDelim() {
	if bytes.Equal(t.delim().Value(), ldelim.Value()) {
		t.del = rdelim
	} else {
		t.del = ldelim
	}
}

func (t *Template) handleVar(v []byte) ([]byte, error) {
	tag := string(cleanSpaces(v))
	if len(tag) == 0 {
		// TODO handle if tag is empty
	}

	esc := true

	switch tag[0] {
	case '#', '^':
		bl := t.newBlock(tag, t.cursor)
		// bl := t.newBlock(tag, t.cursor-(len(v)+len(rdelim)+len(ldelim)))
		if bl == nil {
			// TODO not sure how this can actually happen...?
		}

		return v[:0], nil

	case '/':
		_, bl := t.currentBlock()
		if bl == nil {
			// TODO error: invalid block
		}
		if bl.tag != tag {
			// TODO error: non-matching block
		}
		if bl.Increment(); bl.IsFinished() {
			t.popBlock()

			return nil, nil
		}

		// reset the buffer and move the cursor to the block's cursor
		// location
		t.buf = t.buf[:0]
		t.cursor = bl.cursor

		// set the File's cusror to be read at on the next Read
		_, err := t.File.Seek(int64(t.cursor), 0)

		return nil, err

	case '&':
		tag = tag[1:]
		esc = false

	case '>':
		r, err := t.newPartial(tag[1:])
		if err != nil {
			return nil, err
		}
		t.partial = r

		return nil, nil
	}
	// TODO how to handle/detect unclosed blocks

	val := t.getValue(tag)
	if esc {
		val = escapeBytes(val)
	}

	return val, nil
}

func (t *Template) newBlock(tag string, c int) *block {
	bl := t.findBlock(tag, c)
	if bl != nil {
		return bl
	}

	d := t.Data.Get(tag[1:])
	if d == nil {
		// TODO handle
	}
	bl = newBlock(tag, c, d)

	// lazy alloc
	if t.blocks == nil {
		t.blocks = make([]*block, 0, 32)
	}

	t.blocks = append(t.blocks, bl)

	return bl
}

// findBlock finds a block by it's name and cursot.
// The addition of the cursor provides a method of assigning uniqueness to a
// block, allowing blocks to nest the same block and have a fresh reference to
// the underlying data.
//
// The name provided should not containa any block prefixes,
// eg. #words -> words.
func (t *Template) findBlock(tag string, c int) *block {
	z := len(t.blocks) - 1
	if z < 0 {
		return nil
	}

	// look up block in reverse
	for ; z > -1; z-- {
		bl := t.blocks[z]
		if bl.cursor == c && bl.tag == tag {
			return bl
		}
	}

	return nil
}

// currentBlock returns the last block (and it's index) on the block list,
// which represents the current block.
func (t *Template) currentBlock() (int, *block) {
	z := len(t.blocks) - 1
	if z < 0 {
		return -1, nil
	}

	return z, t.blocks[z]
}

// popBlock pops off the last block in the blocks list
func (t *Template) popBlock() *block {
	i, bl := t.currentBlock()
	if i < 0 {
		return nil
	}
	if bl == nil {
		return nil
	}

	t.blocks = t.blocks[:i]

	return bl
}

func (t *Template) newPartial(path string) (io.Reader, error) {
	if t.partialFunc == nil {
		return nil, errInvalidPartialFunc
	}

	r, err := t.partialFunc(path)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, nil
	}
	// if we get a File, make it a template
	if f, ok := r.(File); ok {
		te := &Template{
			File: f,
			Data: t.Data,
		}
		te.Partial(t.partialFunc)

		r = te
	}

	return r, nil
}

// getValue looks up the value within Data map. It will iterrate *up* the blocks
// before looking at the root Data field itself.
func (t *Template) getValue(k string) []byte {
	z := len(t.blocks)
	for ; z > 0; z-- {
		bl := t.blocks[z-1]
		if bl.Skip() {
			return nil
		}
		if bl.Inverted() {
			continue
		}
		if v := bl.Data().Get(k); v != nil {
			return v.Bytes()
		}
	}

	// . never looks up outside of a block
	if k == "." {
		return []byte{}
	}
	if v := t.Data.Get(k); v != nil {
		return v.Bytes()
	}

	return nil
}

// flush writes truncated out to p. It writes up to the lesser of the two
// lengths, p vs truncd. It returns any remaing bytes that couldn't be written
// due to length constraints.
func (t *Template) flush(p []byte) (int, []byte) {
	lent := len(t.truncd)
	z := len(p)
	if lent < z {
		z = lent
	}

	i := 0
	for ; i < z; i++ {
		p[i] = t.truncd[i]
	}
	if i >= lent {
		return i, t.truncd[:0]
	}

	return i, t.truncd[i:]
}

// Partial sets the partialFunc
func (t *Template) Partial(fn PartialFunc) {
	t.partialFunc = fn
}

// cleanSpaces removes all spaces by shifting over the spaces allow us to return
// a space clean version without allocating
func cleanSpaces(b []byte) []byte {
	lenb := len(b)

	j := 0 // track i minus spaces
	i := 0
	for ; i < lenb; i++ {
		c := b[i]
		if c == ' ' {
			continue
		}

		b[j] = c

		j++
	}

	return b[:j]
}

func closePartial(par interface{}) {
	var c io.ReadCloser

	switch v := par.(type) {
	case *Template:
		f, ok := v.File.(io.ReadCloser)
		if ok {
			c = f
		}
	case io.ReadCloser:
		c = v
	}

	if c == nil {
		return
	}

	c.Close()
}

var (
	errInvalidPartialFunc = errors.New("partial func is undefined")
	errUnclosedBlocks     = errors.New("unclosed blocks")
)
