package beard

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestDataGetBasicMap(t *testing.T) {
	data := map[string]interface{}{
		"a": "Hello",
		"b": map[string]string{
			"c": "World",
		},
		"d": map[string]interface{}{
			"e": map[string]interface{}{
				"f": "!",
			},
		},
	}

	d := Data{Value: data}

	for _, v := range []struct {
		giv, exp string
	}{
		{"a", "Hello"},
		{"b.c", "World"},
		{"d.e.f", "!"},
	} {
		b := d.Get(v.giv)
		if got := string(b.Bytes()); v.exp != got {
			t.Errorf("expected %q, got %q", v.exp, got)
		}
	}
}

func TestDataGetNotKeyable(t *testing.T) {
	data := map[string]interface{}{
		"a": "b",
		"c": []interface{}{
			"d",
		},
	}

	d := Data{Value: data}

	for _, v := range []string{
		"a.b",
		"c.d",
	} {
		b := d.Get(v)
		if b != nil {
			t.Errorf("expected nil, got %s", b)
		}
	}
}

func TestDataGetStructFields(t *testing.T) {
	data := map[string]interface{}{
		"a": "Hello",
		"b": struct {
			C string
		}{
			C: "World",
		},
		"d": &struct {
			E struct {
				f string
			}
		}{
			E: struct {
				f string
			}{
				f: "!",
			},
		},
	}

	d := Data{Value: data}

	for _, v := range []struct {
		giv, exp string
	}{
		{"a", "Hello"},
		{"b.C", "World"},
	} {
		b := d.Get(v.giv)
		if got := string(b.Bytes()); v.exp != got {
			t.Errorf("expected %s, got %s", v.exp, got)
		}
	}

	{
		var exp = "!"

		b := d.Get("d.E.f")
		if got := string(b.Bytes()); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func TestDataGetUnknownPath(t *testing.T) {
	var data = map[string]interface{}{
		"a": "Hello",
		"b": struct {
			c string
		}{
			c: "World",
		},
	}

	d := Data{Value: data}

	for _, v := range []string{
		"a.b",
		"b.d",
		"b.c.d",
	} {
		b := d.Get(v)
		if b != nil {
			t.Errorf("expected nil, got %s", b)
		}
	}
}

func TestDataGetDotReturnsData(t *testing.T) {
	var data = map[string]interface{}{
		"a": "Hello",
	}

	var exp = &Data{Value: data}

	d := Data{Value: data}

	if got := d.Get("."); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected itself, got %s", got)
	}
}

func TestDataVariousFormats(t *testing.T) {
	var data = map[string]interface{}{
		"a": 123,
		"b": int64(456),
		"c": 1.00,
		"d": 1.01,
		"e": true,
		"f": []byte("Hello"),
	}

	d := Data{Value: data}

	{
		var exp = []byte("123")

		if got := d.Get("a").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	{
		var exp = []byte("456")

		if got := d.Get("b").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	{
		var exp = []byte("1")

		if got := d.Get("c").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	{
		var exp = []byte("1.01")

		if got := d.Get("d").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	{
		var exp = []byte("true")

		if got := d.Get("e").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	{
		var exp = []byte("Hello")

		if got := d.Get("f").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func TestDataByteRepresentations(t *testing.T) {
	type a struct {
		B string
	}

	data := map[string]interface{}{
		"a": a{"Hello"},
		"b": &a{"Hello"},
		"c": []string{"Hello", "World"},
	}

	d := Data{Value: data}

	{
		var exp = []byte(fmt.Sprintf("%s", a{"Hello"}))

		if got := d.Get("a").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	{
		var exp = []byte(fmt.Sprintf("%s", &a{"Hello"}))

		if got := d.Get("b").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	{
		var exp = []byte(fmt.Sprintf("%s", []string{"Hello", "World"}))

		if got := d.Get("c").Bytes(); !bytes.Equal(exp, got) {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}
