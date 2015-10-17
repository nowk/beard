package beard

import (
	"bytes"
	"testing"
)

func TestRenderInLayout(t *testing.T) {
	layo := `<body>{{>yield}}</body>`
	html := `<h1>{{a}} {{b}}{{c}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"b": "World",
		"c": "!",
	}

	var exp = `<body><h1>Hello World!</h1></body>`

	tmpl := RenderInLayout(
		mFile{bytes.NewReader([]byte(layo))},
		mFile{bytes.NewReader([]byte(html))},
		data,
	).(*Template)

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))

}
