package beard

import (
	"errors"
	"io"
)

func Render(fi File, d map[string]interface{}, fn PartialFunc) io.Reader {
	te := &Template{
		File: fi,
		Data: &Data{Value: d},
	}
	te.Partial(fn)

	return te
}

// RenderInLayout allows a file to be rendered within a layout. Rendering is
// handled by way of a partial, the partial syntax uses the keyword yield
// eg. {{>yield}}
func RenderInLayout(
	la, fi File, d map[string]interface{}, fn PartialFunc) io.Reader {

	te := Render(fi, d, fn).(*Template)

	layo := &Template{
		File: la,
		Data: te.Data,
	}
	layo.Partial(yieldFunc(te))

	return layo
}

func yieldFunc(r io.Reader) PartialFunc {
	return func(path string) (io.Reader, error) {
		if path == "yield" {
			return r, nil
		}

		return nil, errInvalidYieldTag
	}
}

var errInvalidYieldTag = errors.New("invalid yield tag")
