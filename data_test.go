package beard

import (
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

	d := Data{value: data}

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

	d := Data{value: data}

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

	d := Data{value: data}

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

	d := Data{value: data}

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

	var exp = &Data{value: data}

	d := Data{value: data}

	if got := d.Get("."); !reflect.DeepEqual(exp, got) {
		t.Errorf("expected itself, got %s", got)
	}
}
