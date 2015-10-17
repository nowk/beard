package beard

import (
	"bytes"
	"io"
	"testing"
)

func TestTemplateBufTruncdOut(t *testing.T) {
	html := `<h1>Hello {{c}}</h1>`
	data := map[string]interface{}{
		"c": "world!",
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	var cases = []struct {
		n      int
		nread  int
		buf    string
		truncd string
		out    string
		err    error
	}{
		{5, 5, "", "", "<h1>H", nil},
		{6, 0, "ello {", "", "", nil},
		{3, 3, "c}", "o ", "ell", nil},
		{0, 0, "c}", "o ", "", nil},
		{3, 3, "", "orld!", "o w", nil},
		{3, 3, "", "d!", "orl", nil},
		{15, 7, "", "", "d!</h1>", nil},
		{3, 0, "", "", "", io.EOF},
	}

	i := 0
	for {
		if i > len(cases)-1 {
			break
		}
		var exp = cases[i]

		buf := make([]byte, exp.n)

		n, err := tmpl.Read(buf)
		if err != exp.err {
			t.Errorf("expected %s  error, got %s", exp.err, err)
		}
		if got := n; exp.nread != got {
			t.Errorf("expected to read %d bytes, read %d", exp.nread, got)
		}
		if got := string(tmpl.buf); exp.buf != got {
			t.Errorf("expected buf %s, got %s", exp.buf, got)
		}
		if got := string(tmpl.truncd); exp.truncd != got {
			t.Errorf("expected truncd %s, got %s", exp.truncd, got)
		}
		if got := string(buf[:n]); exp.out != got {
			t.Errorf("expected out %s, got %s", exp.out, got)
		}

		i++
	}

	if exp := len(cases); exp != i {
		t.Errorf("expected %d cases, executed %d", exp, i)
	}
}

func TestTemplateBasicVariables(t *testing.T) {
	html := `<h1>{{a}} {{b}}{{c}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"b": "World",
		"c": "!",
	}

	var exp = `<h1>Hello World!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateArrayBlock(t *testing.T) {
	html := `{{#words}}({{.}})({{.}}){{/words}}`
	data := map[string]interface{}{
		"words": []interface{}{
			"a", "b", "c",
		},
	}

	var exp = `(a)(a)(b)(b)(c)(c)`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateSameArrayInArray(t *testing.T) {
	html := `{{#words}}({{.}}){{#words}}({{.}}){{/words}}{{/words}}`
	data := map[string]interface{}{
		"words": []interface{}{
			"a", "b", "c",
		},
	}

	var exp = `(a)(a)(b)(c)(b)(a)(b)(c)(c)(a)(b)(c)`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateVarPath(t *testing.T) {
	html := `<h1>{{a}} {{b.c}}{{d}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"b": map[string]interface{}{
			"c": "World",
		},
		"d": "!",
	}

	var exp = `<h1>Hello World!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateArrayOfObjects(t *testing.T) {
	html := `{{#words}}({{wo.rd}})({{wo.rd}}){{/words}}`
	data := map[string]interface{}{
		"words": []map[string]interface{}{
			map[string]interface{}{
				"wo": map[string]interface{}{
					"rd": "a",
				},
			},
			map[string]interface{}{
				"wo": map[string]interface{}{
					"rd": "b",
				},
			},
			map[string]interface{}{
				"wo": map[string]interface{}{
					"rd": "c",
				},
			},
		},
	}

	var exp = `(a)(a)(b)(b)(c)(c)`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateArrayInPath(t *testing.T) {
	html := `{{#lots.o.words}}({{.}})({{.}}){{/lots.o.words}}`
	data := map[string]interface{}{
		"lots": map[string]interface{}{
			"o": map[string]interface{}{
				"words": []interface{}{
					"a", "b", "c",
				},
			},
		},
	}

	var exp = `(a)(a)(b)(b)(c)(c)`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateObjectBlock(t *testing.T) {
	html := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/gretting}}</h1>`
	data := map[string]interface{}{
		"greeting": map[string]interface{}{
			"a": "Hello",
			"b": "World",
			"c": map[string]interface{}{
				"d": "!",
			},
		},
	}

	var exp = `<h1>Hello World!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateOutsideOfBlockVar(t *testing.T) {
	html := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/greeting}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"greeting": map[string]interface{}{
			"b": "World",
			"c": map[string]interface{}{
				"d": "!",
			},
		},
	}

	var exp = `<h1>Hello World!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateOutsideOfBlockVarUsesClosestVar(t *testing.T) {
	html := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/greeting}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"greeting": map[string]interface{}{
			"a": "Hola",
			"b": "World",
			"c": map[string]interface{}{
				"d": "!",
			},
		},
	}

	var exp = `<h1>Hola World!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateEscapesStrings(t *testing.T) {
	html := `<code>{{code}}</code>`
	data := map[string]interface{}{
		"code": "<h1>Hello World!</h1>",
	}

	var exp = `<code>&lt;h1&gt;Hello World!&lt;/h1&gt;</code>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateDontEscapesStrings(t *testing.T) {
	html := `<code>{{&code}}</code>`
	data := map[string]interface{}{
		"code": "<h1>Hello World!</h1>",
	}

	var exp = `<code><h1>Hello World!</h1></code>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateNotEscapeDelimDoesNotAttemptToMatch(t *testing.T) {
	html := `<code>{{&code}}</code>`
	data := map[string]interface{}{
		"code": "{{c}}",
	}

	var exp = `<code>{{c}}</code>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplatePartial(t *testing.T) {
	html := `<h1>{{a}}{{>b}}{{e}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"d": "World",
		"e": "!",
	}

	var exp = `<h1>Hello World!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}
	tmpl.Partial(func(path string) (File, error) {
		var p []byte
		if path == "b" {
			p = []byte(` {{>c}}`)
		} else {
			p = []byte(`{{d}}`)
		}

		return bytes.NewReader(p), nil
	})

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateErrorsUnclosedBlock(t *testing.T) {
	t.Skip()
}

var a = func(tmpl *Template) StepFunc {
	return func(t testing.TB, ctx Context) {
		var buffer = make([]byte, 0, 32)
		buf := bytes.NewBuffer(buffer)
		n, err := io.Copy(buf, tmpl)

		ctx.Set("tmpl", tmpl)
		ctx.Set("buf", buf)
		ctx.Set("n", n)
		ctx.Set("err", err)
	}
}

var bodyEquals = func(exp string) StepFunc {
	return func(t testing.TB, ctx Context) {
		var (
			buf = ctx.Get("buf").(*bytes.Buffer)
			n   = ctx.Get("n").(int64)
		)

		if exp := int64(len(exp)); exp != n {
			t.Errorf("expected %d bytes read, got %d", exp, n)
		}
		if got := buf.String(); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

var errorIs = func(exp error) StepFunc {
	return func(t testing.TB, ctx Context) {
		var err = ctx.Get("err")
		if err != exp {
			t.Errorf("expected no error, got %s", err)
		}
		if exp == nil {
			return
		}

		if got := err.(error); exp != got {
			t.Errorf("expected %s error, got %s", exp, got)
		}
	}
}
