package beard

import (
	"bytes"
	"io"
	"log"
)

type Renderable struct {
	// File is the template file to be rendered
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
	p = p[:0]

	if len(r.truncd) > 0 {
		n := flush(&r.truncd, &p, lenp)

		// see if there is any buffer left to write to
		if lenp = lenp - n; lenp == 0 {
			return n, nil
		}
	}
	buf := make([]byte, lenp)

	n, err := r.File.Read(buf)
	if err != nil && err != io.EOF {
		return n, err
	}
	if !r.eof {
		r.eof = (err == io.EOF)
	}

	r.buf = append(r.buf, buf[:n]...)
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
			log.Printf("mat:%s", string(v))

			dat, ok := r.Data[string(v)]
			if ok {
				v = []byte(dat.(string))
			} else {
				v = v[:0]
			}
		}

		// combine truncated with current value to be written
		v = append(r.truncd, v...)

		// truncate if v is longer than our given reader bytes
		if len(v) > lenp {
			p = append(p, v[:lenp]...)
			r.truncd = v[lenp:]
		} else {
			p = append(p, v...)
			r.truncd = r.truncd[:0]
		}

		r.buf = r.buf[len(b):] // trim buffer

		// log.Printf("mat:%s", string(v))
		log.Printf("buf:%s", string(r.buf))

		swapDelim(r)

	default:
		flush(&r.buf, &p, lenp)
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

// flush b to out up to n
func flush(b, out *[]byte, n int) int {
	var (
		_b   = *b
		_out = *out
	)
	if len(_b) > n {
		*out = append(_out, _b[:n]...)
		*b = _b[n:]
	} else {
		*out = append(_out, _b...)
		*b = _b[:0]
	}

	return len(*out)
}

// swapDelim swaps the delim on a Renderable.
func swapDelim(r *Renderable) {
	if bytes.Equal(r.delim(), ldelim) {
		r.del = rdelim
	} else {
		r.del = ldelim
	}
}
