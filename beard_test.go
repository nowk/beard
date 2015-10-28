package beard

import (
	"bytes"
	"io"
	"testing"
)

func TestRenderInLayout(t *testing.T) {
	layo := `<body>{{>yield}}</body>`
	html := `<h1>{{a}} {{>b}}{{d}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"c": "World",
		"d": "!",
	}

	var exp = `<body><h1>Hello World!</h1></body>`

	tmpl := RenderInLayout(
		bytes.NewReader([]byte(layo)),
		bytes.NewReader([]byte(html)),
		data,
		func(path string) (io.Reader, error) {
			if path == "b" {
				return bytes.NewReader([]byte(`{{c}}`)), nil
			}

			t.Errorf("invalid path %s", path)
			return nil, nil
		},
	).(*Template)

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestRenderInLayoutYieldError(t *testing.T) {
	layo := `<body>{{>badyieldtag}}</body>`
	html := `<h1>{{a}} {{>b}}{{d}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"c": "World",
		"d": "!",
	}

	var exp = `<body>`

	tmpl := RenderInLayout(
		bytes.NewReader([]byte(layo)),
		bytes.NewReader([]byte(html)),
		data,
		func(path string) (io.Reader, error) {
			if path == "b" {
				return bytes.NewReader([]byte(`{{c}}`)), nil
			}

			t.Errorf("invalid path %s", path)
			return nil, nil
		},
	).(*Template)

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(errInvalidYieldTag))
}
