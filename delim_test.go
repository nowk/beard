package beard

import (
	"reflect"
	"testing"
)

func Test_ldelim(t *testing.T) {
	for _, v := range []struct {
		giv string
		del string
		exp string
		ma  matchLevel
	}{
		{"hello {{c}}", "{{", "hello {{", exMatch},
		{"hello {{{c}}}", "{{", "hello {{", exMatch},
		{"hello {", "{{", "hello {", paMatch},
		{"hello {", "{{{", "hello {", paMatch},
		{"hello {{", "{{{", "hello {{", paMatch},
		{"hello {c}", "{{", "hello {c}", noMatch},
		{"hello {c}", "{{{", "hello {c}", noMatch},
		{"hello {{c}}", "{{{", "hello {{c}}", noMatch},
	} {
		var exp = struct {
			byt []byte
			ma  matchLevel
		}{
			[]byte(v.exp), v.ma,
		}

		byt, ma := (&Ldelim{[]byte(v.del)}).Match([]byte(v.giv))
		if exp.ma != ma {
			t.Errorf("expected a level %d match, got %d: %s %s",
				exp.ma, ma, v.giv, v.del)
		}

		if !reflect.DeepEqual(exp.byt, byt) {
			t.Errorf("expected %s, got %s", string(exp.byt), string(byt))
		}
	}
}
func Test_rdelim(t *testing.T) {
	for _, v := range []struct {
		giv string
		del string
		exp string
		ma  matchLevel
	}{
		{"c}}</h1>", "}}", "c}}", exMatch},
		{"c}}}</h1>", "}}", "c}}", exMatch},
		{"c}", "}}", "c}", paMatch},
		{"c}", "}}}", "c}", paMatch},
		{"c}}", "}}}", "c}}", paMatch},
		{"c}</h1>", "}}", "c}</h1>", noMatch},
		{"c}</h1>", "}}}", "c}</h1>", noMatch},
		{"c}}</h1>", "}}}", "c}}</h1>", noMatch},
	} {
		var exp = struct {
			byt []byte
			ma  matchLevel
		}{
			[]byte(v.exp), v.ma,
		}

		byt, ma := (&Rdelim{[]byte(v.del)}).Match([]byte(v.giv))
		if exp.ma != ma {
			t.Errorf("expected a level %d match, got %d: %s %s",
				exp.ma, ma, v.giv, v.del)
		}

		if !reflect.DeepEqual(exp.byt, byt) {
			t.Errorf("expected %s, got %s", string(exp.byt), string(byt))
		}
	}
}
