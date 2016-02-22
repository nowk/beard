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
	html := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/greeting}}</h1>`
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
	tmpl.Partial(func(path string) (io.Reader, error) {
		var p []byte
		switch path {
		case "b":
			p = []byte(` {{>c}}`)
		case "c":
			p = []byte(`{{d}}`)

		default:
			t.Errorf("invalid partial %s", path)
		}

		return bytes.NewReader(p), nil
	})

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

type mFileCloser struct {
	io.ReadSeeker

	CloseFunc func() error
}

func (m mFileCloser) Close() error {
	return m.CloseFunc()
}

func TestTemplateClosesPartials(t *testing.T) {
	html := `<h1>{{a}}{{>b}}{{e}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"d": "World",
		"e": "!",
	}

	var exp = struct {
		body   string
		closed int
	}{
		body:   `<h1>Hello World!</h1>`,
		closed: 2,
	}

	var closed = 0

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}
	tmpl.Partial(func(path string) (io.Reader, error) {
		var p []byte
		switch path {
		case "b":
			p = []byte(` {{>c}}`)
		case "c":
			p = []byte(`{{d}}`)

		default:
			t.Errorf("invalid partial %s", path)
		}

		return mFileCloser{
			ReadSeeker: bytes.NewReader(p),
			CloseFunc: func() error {
				closed++

				return nil
			},
		}, nil
	})

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp.body)).
		And(errorIs(nil))

	if exp.closed != closed {
		t.Errorf(
			"expected to have closed %d partials, got %d", exp.closed, closed)
	}
}

func TestTemplatePartialFuncIsNil(t *testing.T) {
	html := `<h1>{{a}}{{>b}}{{e}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"d": "World",
		"e": "!",
	}

	var exp = `<h1>Hello`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(errInvalidPartialFunc))
}

func TestTemplatePartialInBlock(t *testing.T) {
	html := `<h1>{{#a}}{{>b}}{{f}}{{/a}}</h1>`
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"c": "Hello",
			"e": "World",
			"f": "!",
		},
	}

	var exp = `<h1>Hello World!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}
	tmpl.Partial(func(path string) (io.Reader, error) {
		var p []byte
		switch path {
		case "b":
			p = []byte(`{{c}}{{>d}}`)
		case "d":
			p = []byte(` {{e}}`)

		default:
			t.Errorf("invalid partial %s", path)
		}

		return bytes.NewReader(p), nil
	})

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateVarsContainSpaces(t *testing.T) {
	html := `<h1>{{a }}{{ > b }}{{ e}}</h1>`
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
	tmpl.Partial(func(path string) (io.Reader, error) {
		var p []byte
		switch path {
		case "b":
			p = []byte(` {{> c}}`)
		case "c":
			p = []byte(`{{d}}`)

		default:
			t.Error("invalid partial %s", path)
		}

		return bytes.NewReader(p), nil
	})

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateEmptyNilUnexistBlockData(t *testing.T) {
	html := `<h1>{{#words}}({{.}}){{/words}}</h1>`

	for _, v := range []interface{}{
		map[string]interface{}{
			"words": []interface{}{},
		},
		map[string]interface{}{
			"words": nil,
		},
		map[string]interface{}{
			"badwords": []string{"foo"},
		},
	} {
		var exp = `<h1></h1>`

		tmpl := &Template{
			File: bytes.NewReader([]byte(html)),
			Data: &Data{Value: v},
		}

		Asser{t}.
			Given(a(tmpl)).
			Then(bodyEquals(exp)).
			And(errorIs(nil))
	}
}

func TestTemplateInvertedBlock(t *testing.T) {
	html := `<h1>{{#words}}({{.}}){{/words}}{{^words}}Hola Mundo!{{/words}}</h1>`

	for _, v := range []struct {
		data map[string]interface{}
		exp  string
	}{
		{map[string]interface{}{"words": []string{"a", "b", "c"}}, `<h1>(a)(b)(c)</h1>`},
		{map[string]interface{}{}, `<h1>Hola Mundo!</h1>`},
	} {
		tmpl := &Template{
			File: bytes.NewReader([]byte(html)),
			Data: &Data{Value: v.data},
		}

		Asser{t}.
			Given(a(tmpl)).
			Then(bodyEquals(v.exp)).
			And(errorIs(nil))
	}
}

func TestTemplateInvertedBlockDoesNotTraverseUp(t *testing.T) {
	html := `<h1>{{#many.words}}({{.}}){{/many.words}}{{^many.words}}Hola Mundo!{{/many.words}}</h1>`
	data := map[string]interface{}{
		"words": []string{
			"a", "b", "c",
		},
		"many": map[string]interface{}{
			"words": []string{},
		},
	}

	var exp = `<h1>Hola Mundo!</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateBlockAs(t *testing.T) {
	{
		var html = `{{#words as word}}{{word}}{{/words}}`

		data := map[string]interface{}{
			"words": []string{
				"a", "b", "c",
			},
		}

		tmpl := &Template{
			File: bytes.NewReader([]byte(html)),
			Data: &Data{Value: data},
		}

		var exp = "abc"

		Asser{t}.
			Given(a(tmpl)).
			Then(bodyEquals(exp)).
			And(errorIs(nil))
	}

	{
		var html = `{{#words as word}}{{word.value}}{{/words}}`

		data := map[string]interface{}{
			"words": []map[string]interface{}{
				map[string]interface{}{
					"value": "a",
				},
				map[string]interface{}{
					"value": "b",
				},
				map[string]interface{}{
					"value": "c",
				},
			},
		}

		tmpl := &Template{
			File: bytes.NewReader([]byte(html)),
			Data: &Data{Value: data},
		}

		var exp = "abc"

		Asser{t}.
			Given(a(tmpl)).
			Then(bodyEquals(exp)).
			And(errorIs(nil))
	}
}

func TestTemplateBlockAsKeyValue(t *testing.T) {
	var html = `{{#words as k, v}}{{k}}{{/words}}`

	data := map[string]interface{}{
		"words": map[string]interface{}{
			"a": "b",
			"c": "d",
			"e": "f",
		},
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	var exp = "ace"

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestTemplateErrorsUnclosedBlock(t *testing.T) {
	html := `<h1>{{#words}}({{.}})</h1>`
	data := map[string]interface{}{
		"words": []string{
			"a", "b", "c",
		},
	}

	var exp = `<h1>(a)</h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(errUnclosedBlocks))
}

func TestTemplateErrorEmptyTag(t *testing.T) {
	html := `<h1>{{}}</h1>`

	var exp = `<h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: nil,
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(errEmptyTag))
}

func TestTemplateErrorNilBlock(t *testing.T) {
	html := `<h1>{{#words}}{{.}}{{/words}}</h1>`
	data := map[string]interface{}{
		"words": []string{"a"},
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	b := make([]byte, 19)

	tmpl.Read(b)
	tmpl.Read(b)
	tmpl.Read(b)
	tmpl.Read(b)
	tmpl.Read(b)

	tmpl.blocks[0] = nil

	_, err := tmpl.Read(b)
	if err != errNilBlock {
		t.Errorf("expected no error, got %s", err)
	}
}

func TestTemplateErrorMismatchBlock(t *testing.T) {
	html := `<h1>{{#words}}{{.}}{{/words}}</h1>`
	data := map[string]interface{}{
		"words": []string{"a"},
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	b := make([]byte, 19)

	tmpl.Read(b)
	tmpl.Read(b)
	tmpl.Read(b)
	tmpl.Read(b)
	tmpl.Read(b)

	tmpl.blocks[0].tag = "/sentences"

	_, err := tmpl.Read(b)
	if err != errBlockMismatch {
		t.Errorf("expected no error, got %s", err)
	}
}

func TestTemplateBlockNoData(t *testing.T) {
	html := `<h1>{{#words}}{{.}}{{/words}}</h1>`
	data := map[string]interface{}{}

	var exp = `<h1></h1>`

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	Asser{t}.
		Given(a(tmpl)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
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
