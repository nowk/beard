package beard

import (
	"bytes"
	"io"
	"testing"
)

func TestRenderableBufTruncdOut(t *testing.T) {
	tmpl := `<h1>Hello {{c}}</h1>`
	data := map[string]interface{}{
		"c": "world!",
	}

	rend := &Renderable{
		File: bytes.NewReader([]byte(tmpl)),
		Data: data,
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

		n, err := rend.Read(buf)
		if err != exp.err {
			t.Errorf("expected %s  error, got %s", exp.err, err)
		}
		if got := n; exp.nread != got {
			t.Errorf("expected to read %d bytes, read %d", exp.nread, got)
		}
		if got := string(rend.buf); exp.buf != got {
			t.Errorf("expected buf %s, got %s", exp.buf, got)
		}
		if got := string(rend.truncd); exp.truncd != got {
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

func TestRenderableBasicVariables(t *testing.T) {
	tmpl := `<h1>{{a}} {{b}}{{c}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"b": "World",
		"c": "!",
	}

	var exp = `<h1>Hello World!</h1>`

	rend := &Renderable{
		File: bytes.NewReader([]byte(tmpl)),
		Data: data,
	}

	Asser{t}.
		Given(a(rend)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestRenderableArrayBlock(t *testing.T) {
	tmpl := `{{#words}}({{.}})({{.}}){{/words}}`
	data := map[string]interface{}{
		"words": []interface{}{
			"a", "b", "c",
		},
	}

	var exp = `(a)(a)(b)(b)(c)(c)`

	rend := &Renderable{
		File: bytes.NewReader([]byte(tmpl)),
		Data: data,
	}

	Asser{t}.
		Given(a(rend)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

func TestRenderableSameArrayInArray(t *testing.T) {
	tmpl := `{{#words}}({{.}}){{#words}}({{.}}){{/words}}{{/words}}`
	data := map[string]interface{}{
		"words": []interface{}{
			"a", "b", "c",
		},
	}

	var exp = `(a)(a)(b)(c)(b)(a)(b)(c)(c)(a)(b)(c)`

	rend := &Renderable{
		File: bytes.NewReader([]byte(tmpl)),
		Data: data,
	}

	Asser{t}.
		Given(a(rend)).
		Then(bodyEquals(exp)).
		And(errorIs(nil))
}

var a = func(rend *Renderable) StepFunc {
	return func(t testing.TB, ctx Context) {
		buf := bytes.NewBuffer(nil)
		n, err := io.Copy(buf, rend)

		ctx.Set("rend", rend)
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
