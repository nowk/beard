package beard

import (
	"bytes"
	"io"
	"log"

	"github.com/spf13/afero"
)

type File struct {
	afero.File

	Pwd          string
	TemplateRoot string
}

var _ io.Reader = &File{}

type Renderable struct {
	afero.File

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
		log.Printf("mat:%s", string(v))

		// truncate if v is longer than our given reader bytes
		if len(v) > lenp {
			p = append(p, v[:lenp]...)
			r.tru = v[lenp:]
		} else {
			p = append(p, v...)
			r.tru = r.tru[:0]
		}

		r.buf = r.buf[len(b):] // trim buffer

		log.Printf("mat:%s", string(v))
		log.Printf("buf:%s", string(r.buf))

		swapDelim(r)
	}

	// attempt to flush the rest of the buffer
	if r.eof {
		if len(r.buf) > lenp {
			p = append(p, r.buf[:lenp]...)
			r.buf = r.buf[lenp:]
		} else {
			p = append(p, r.buf...)
			r.buf = r.buf[:0]

			return len(p), io.EOF
		}
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
