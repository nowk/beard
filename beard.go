package beard

import (
	"io"
)

// Render provides s a convenience func to create
func Render(
	fi File, d map[string]interface{}, opts ...func(*Template)) io.Reader {

	te := &Template{
		File: fi,
		Data: &Data{Value: d},
	}

	for _, v := range opts {
		v(te)
	}

	return te
}

// RenderInLayout allows a file to be rendered within a layout. Rendering is
// handled by way of a partial, the partial syntax uses the keyword yield
// eg. {{>yield}}
func RenderInLayout(
	la, fi File, d map[string]interface{}, opts ...func(*Template)) io.Reader {

	te := Render(fi, d, opts...).(*Template)

	return &Template{
		File: la,
		Data: te.Data,

		PartialFunc: layoutFunc(te),
	}
}

var layoutFunc = func(v interface{}) PartialFunc {
	return func(path string) (interface{}, error) {
		if path == "yield" {
			return v, nil
		}

		// TODO should we error here?
		return nil, nil
	}
}
