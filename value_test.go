package beard

import (
	"reflect"
	"testing"
)

func Test_getvofBasicMap(t *testing.T) {
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

	for _, v := range []struct {
		giv, exp string
	}{
		{"a", "Hello"},
		{"b.c", "World"},
		{"d.e.f", "!"},
	} {
		b := getvof(v.giv, data)
		if got := b.(string); v.exp != got {
			t.Errorf("expected %q, got %q", v.exp, got)
		}
	}
}

func Test_getvofNotKeyable(t *testing.T) {
	data := map[string]interface{}{
		"a": "b",
		"c": []interface{}{
			"d",
		},
	}

	for _, v := range []string{
		"a.b",
		"c.d",
	} {
		b := getvof(v, data)
		if b != nil {
			t.Errorf("expected nil, got %s", b)
		}
	}
}

func Test_getvofStructFields(t *testing.T) {
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

	for _, v := range []struct {
		giv, exp string
	}{
		{"a", "Hello"},
		{"b.C", "World"},
	} {
		b := getvof(v.giv, data)
		if got := b.(string); v.exp != got {
			t.Errorf("expected %s, got %s", v.exp, got)
		}
	}

	{
		var exp = "!"

		b := getvof("d.E.f", data)
		if got := b.(reflect.Value).String(); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_getvofUnknownPath(t *testing.T) {
	var data = map[string]interface{}{
		"a": "Hello",
		"b": struct {
			c string
		}{
			c: "World",
		},
	}

	for _, v := range []string{
		"a.b",
		"b.d",
		"b.c.d",
	} {
		b := getvof(v, data)
		if b != nil {
			t.Errorf("expected nil, got %s", b)
		}
	}
}
