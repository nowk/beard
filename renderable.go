package beard

import (
	"bytes"
	"io"
)

type Renderable struct {
	// File is the template file to be rendered. File must be explicitly closed
	// by the user
	File io.Reader

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
		v := b[:len(b)-len(del)]

		// when we find a matching rdelim, {{..}} has been closed and we can now
		// parse for the var value
		if bytes.Equal(del, rdelim) {
			dat, ok := r.Data[string(v)]
			if ok {
				v = []byte(dat.(string))
			} else {
				v = v[:0]
			}
		}

		// combine truncated with current value and write
		if val := append(r.truncd, v...); len(val) > lenp {
			n := write(val, &p, writ, lenp)

			r.truncd = val[n:]
		} else {
			_ = write(val, &p, writ, len(val))

			r.truncd = r.truncd[:0]
		}

		// trim buffer of our matched bytes
		r.buf = r.buf[len(b):]
		swapDelim(r)

	default:
		// if we have a buf, flush it
		// a buf at this point will always be within the length of p
		if lenbuf := len(r.buf); lenbuf > 0 {
			_ = write(r.buf, &p, writ, writ+lenbuf)

			r.buf = r.buf[:0]
		}
	}

	if r.eof && len(r.buf) == 0 {
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

// write writes b to out, starting at i to n
func write(b []byte, out *[]byte, i, n int) int {
	_out := *out

	z := 0
	for ; i < n; i++ {
		*out = _out[:i+1]
		(*out)[i] = b[z]

		z++
	}

	return z
}

// swapDelim swaps the delim on a Renderable.
func swapDelim(r *Renderable) {
	if bytes.Equal(r.delim(), ldelim) {
		r.del = rdelim
	} else {
		r.del = ldelim
	}
}
