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

	// tru holds bytes that were truncated due to the len of the reader bytes
	tru []byte

	// eof marks the File has reached EOF.
	eof bool
}

var _ io.Reader = &Renderable{}

func (r *Renderable) Read(p []byte) (int, error) {
	lenp := len(p)
	buf := make([]byte, lenp)

	p = p[:0]

	n, err := r.File.Read(buf)
	if err != nil && err != io.EOF {
		// TODO must still read what was read to p
		return 0, err
	}
	if !r.eof {
		r.eof = (err == io.EOF)
	}

	r.buf = append(r.buf, buf[:n]...)
	del := r.delim()

	// log.Printf("buf:%s", string(r.buf))
	// log.Printf("del:%s", string(del))

	b, ma := matchdel(r.buf, del)
	if ma == paMatch {
		r.buf = b
	}
	if ma == exMatch {
		v := append(r.tru, b[:len(b)-len(del)]...)

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

		// truncate if v is longer than our given reader bytes
		if len(v) > lenp {
			p = append(p, v[:lenp]...)
			r.tru = v[lenp:]
		} else {
			p = append(p, v...)
			r.tru = r.tru[:0]
		}

		r.buf = r.buf[len(b):] // trim buffer

		// log.Printf("mat:%s", string(v))
		log.Printf("buf:%s", string(r.buf))

		swapDelim(r)
	}
	if ma == noMatch {
		if len(r.buf) > lenp {
			p = append(p, r.buf[:lenp]...)
			r.buf = r.buf[lenp:]
		} else {
			p = append(p, r.buf...)
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

// swapDelim swaps the delim on a Renderable.
func swapDelim(r *Renderable) {
	if bytes.Equal(r.delim(), ldelim) {
		r.del = rdelim
	} else {
		r.del = ldelim
	}
}
