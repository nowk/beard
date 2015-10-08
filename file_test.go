package beard

import (
	"bytes"
	"io"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestRenderable(t *testing.T) {
	f := afero.MemFileCreate("-")
	f.WriteString(`<h1>hello {{c}}</h1>`)
	f.Seek(0, 0)

	fi := &Renderable{File: f}

	{
		b := make([]byte, 11)

		var exp = []byte(`<h1>hello {`)

		n, err := fi.Read(b)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
		if n != 0 {
			t.Errorf("expected 0 bytes read, got %d", n)
		}
		if got := fi.buf; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %s, got %s", string(exp), string(got))
		}
	}

	{
		b := make([]byte, 9)

		var exp = []byte(`<h1>hello`)

		n, err := fi.Read(b)
		if err != nil {
			t.Errorf("expected no error, got %s", err)
		}
		if n != 9 {
			t.Errorf("expected 9 bytes read, got %d", n)
		}
		if got := b[:n]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %s, got %s", string(exp), string(got))
		}
		if got := fi.buf; !reflect.DeepEqual([]byte("c}}</h1>"), got) {
			t.Errorf("expected %s, got %s", "c}}</h1>", string(got))
		}
	}

	{
		b := make([]byte, 9)

		var exp = []byte(` c</h1>`)

		n, err := fi.Read(b)
		if err != io.EOF {
			t.Errorf("expected EOF error, got %s", err)
		}
		if n != 7 {
			t.Errorf("expected 7 bytes read, got %d", n)
		}
		if got := b[:n]; !reflect.DeepEqual(exp, got) {
			t.Errorf("expected %s, got %s", string(exp), string(got))
		}
	}
}

func TestRenderable2(t *testing.T) {
	f := afero.MemFileCreate("-")
	f.WriteString(`<h1>hello {{c}}</h1>`)
	f.Seek(0, 0)

	fi := &Renderable{File: f}

	b := bytes.NewBuffer(nil)

	var exp = `<h1>hello c</h1>`

	n, err := io.Copy(b, fi)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	if n != 16 {
		t.Errorf("expected 10 bytes read, got %d", n)
	}

	if got := b.String(); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected %s, got %s", string(exp), string(got))
	}
}
