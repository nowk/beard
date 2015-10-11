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
	Data map[string]interface{}

	// del represents the current delimiter
	del []byte

	// buf holds bytes that end with a partially matched delimiter. It will
	// attempt an exact match on the next read
	buf []byte

	// truncd holds bytes that were truncated due to length of p
	truncd []byte

	// eof marks the File has reached EOF.
	eof bool

	// cursor is the location at which the reader is at
	cursor int

	// block represents the current block in the file.
	// blocks are added in a FILO order. The last block in the list would be the
	// current block
	blocks []*block
}

var _ io.Reader = &Renderable{}

func (r *Renderable) Read(p []byte) (int, error) {
	lenp := len(p)
	writ := 0

	// flush trucncated out to write
	if len(r.truncd) > 0 {
		if writ = flush(r.truncd, p, lenp); len(r.truncd) > writ {
			r.truncd = r.truncd[writ:]
		} else {
			r.truncd = r.truncd[:0]
		}

		// return if we've written out to the length of p
		if n := lenp - writ; n == 0 {
			return writ, nil
		}
	}

	// share the given []byte argument so we don't have to allocate a temp
	// buffer on each Read call.

	n, err := r.File.Read(p[writ:])
	if err != nil && err != io.EOF {
		return n, err
	}
	if !r.eof {
		r.eof = (err == io.EOF)
	}

	// combine buffered with what was just read
	r.buf = append(r.buf, p[writ:writ+n]...)

	// if nothing was written reset p, else trim
	if writ == 0 {
		p = p[:0]
	} else {
		p = p[:writ+n]
	}

	del := r.delim()

	b, ma := matchdel(r.buf, del)
	switch ma {
	case paMatch:
		r.buf = b

	case exMatch:
		lenb := len(b)

		r.cursor += lenb
		r.buf = r.buf[lenb:]

		swapDelim(r) // swap delim early, blocks will return early

		v := b[:lenb-len(del)]

		// when we find a matching rdelim, {{..}} has been closed and we can now
		// parse for the var value
		if bytes.Equal(del, rdelim) {
			k := string(bytes.TrimSpace(v))

			// TODO handle if k is empty

			switch k[0] {
			case '#':
				v = v[:0]

				bl := r.newBlock(k[1:], r.cursor-(lenb+len(ldelim)))
				if bl == nil {
					// TODO handle
				}

				// log.Printf("> [%d] %s", bl.cursor, bl.name)
				// log.Printf("# of blocks: %d", len(r.blocks))

			case '/':
				v = v[:0]

				_, bl := r.currentBlock()
				if bl == nil {
					// TODO handle
				}
				if bl.name != k[1:] {
					// TODO handle
				}
				// log.Printf("[%d] %s", bl.cursor, bl.name)
				// log.Printf("-- %s", string(r.buf))
				if bl.increment(); bl.isFinished() {
					r.popBlock()

					return len(p), nil
				}

				// reset the buffer and move the cursor to the block's cursor
				// location
				r.buf = r.buf[:0]
				r.cursor = bl.cursor

				// set the File's cusror to be read at on the next Read
				_, err := r.File.Seek(int64(r.cursor), 0)
				if err != nil {
					return len(p), err
				}

				return len(p), nil

			default:
				v = r.getValue(k)
			}
		}

		// combine truncated with current value and write
		v = append(r.truncd, v...)
		z := 0

		if lenv := len(v); lenv > lenp {
			z = lenp - writ

			r.truncd = v[z:]
		} else {
			z = lenv - writ

			r.truncd = r.truncd[:0]
		}

		p = append(p[:writ], v[:z]...)

	default:
		// TODO if we've read all and haven't found the rdelim, error.
		// also put a threshold. vars paths should not really be overly too long.
		if bytes.Equal(del, rdelim) {
			break
		}

		// if we have a buf, flush it
		// a buf at this point will always fit into p
		if lenbuf := len(r.buf); lenbuf > 0 {
			p = append(p[:writ], r.buf[:lenbuf]...)

			r.buf = r.buf[:0]
			r.cursor += len(p) // TODO will need to check this
		}
	}

	if r.eof && len(r.buf) == 0 {
		// log.Printf("EOF")
		return len(p), io.EOF
	}

	return len(p), nil
}

// delim returns the "current" delim in the Renderable, it defaults to ldelim.
func (r *Renderable) delim() []byte {
	if r.del == nil {
		return ldelim
	}

	return r.del
}

func (r *Renderable) newBlock(name string, c int) *block {
	// _, bl := r.currentBlock()
	bl := r.findBlock(name, c)
	if bl != nil {
		// TODO check that the returned block and the name match
		return bl
	}

	d, ok := r.Data[name]
	if !ok {
		// TODO handle
	}

	bl = newBlock(name, c, d)
	r.blocks = append(r.blocks, bl)

	return bl
}

// findBlock finds a block by it's name and cursor.
// The addition of the cursor provides a method of assigning uniqueness to a
// block, allowing blocks to nest the same block and have a fresh reference to
// the underlying data.
//
// The name provided should not containa any block prefixes,
// eg. #words -> words
func (r *Renderable) findBlock(name string, c int) *block {
	for _, v := range r.blocks {
		if v.name == name && v.cursor == c {
			return v
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
		if d := bl.getData(k); d != nil {
			return valueByte(d)
		}
	}

	// . never looks up outside of a block
	if k == "." {
		return []byte{}
	}

	d, ok := r.Data[k]
	if !ok {
		return []byte{}
	}

	return valueByte(d)
}

// flush writes b to out, up to max
func flush(b, out []byte, max int) int {
	if lenb := len(b); lenb < max {
		max = lenb
	}

	i := 0
	for ; i < max; i++ {
		out[i] = b[i]
	}

	return i
}

// swapDelim swaps the delim on a Renderable.
func swapDelim(r *Renderable) {
	if bytes.Equal(r.delim(), ldelim) {
		r.del = rdelim
	} else {
		r.del = ldelim
	}
}

func valueByte(d interface{}) []byte {
	switch v := d.(type) {
	case string:
		return []byte(v)
	}

	return []byte{}
}
