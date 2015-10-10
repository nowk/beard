package beard

import (
	"bytes"
	"io"
	"testing"
)

func TestRenderableBufTruncdOut(t *testing.T) {
	file := bytes.NewBufferString(`<h1>Hello {{c}}</h1>`)

	rend := &Renderable{
		File: file,
		Data: map[string]interface{}{
			"c": "world!",
		},
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

func TestRenderableReader(t *testing.T) {
	file := bytes.NewBufferString(`<h1>Hello {{word}}{{d}}</h1>`)

	rend := &Renderable{
		File: file,
		Data: map[string]interface{}{
			"word": "World!",
		},
	}

	var exp = `<h1>Hello World!</h1>`

	buf := bytes.NewBuffer(nil)
	n, err := io.Copy(buf, rend)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if exp := int64(21); exp != n {
		t.Errorf("expected %d bytes read, got %d", exp, n)
	}
	if got := buf.String(); exp != got {
		t.Errorf("expected %s, got %s", exp, got)
	}
}
