package beard

import (
	"testing"
)

func Test_blockisFinished(t *testing.T) {
	bl := newBlock("", 0, &Data{Value: []string{
		"a", "b",
	}})

	if bl.Finished() {
		t.Errorf("expected block to not be finished")
	}

	bl.Increment()
	bl.Increment()

	if !bl.Finished() {
		t.Errorf("expected block to be finished")
	}
}

func Test_blockgetValueDotOnSlice(t *testing.T) {
	bl := newBlock("", 0, &Data{Value: []interface{}{
		"a",
		"b",
	}})

	{
		var exp = "a"

		d := bl.Data().Get(".")
		if got := string(d.Bytes()); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var exp = "b"

		d := bl.Data().Get(".")
		if got := string(d.Bytes()); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockgetValueOnMap(t *testing.T) {
	bl := newBlock("", 0, &Data{Value: map[string]interface{}{
		"a": "b",
		"c": "d",
	}})

	{
		var (
			exp = "b"
			got = string(bl.Data().Get("a").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
	{
		var (
			exp = "d"
			got = string(bl.Data().Get("c").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockgetValuePathOnSlice(t *testing.T) {
	bl := newBlock("", 0, &Data{Value: []interface{}{
		map[string]interface{}{
			"a": map[string]interface{}{
				"b": "Hello",
			},
		},
		map[string]interface{}{
			"a": map[string]interface{}{
				"b": "World",
			},
		},
	}})

	{
		var exp = "Hello"

		d := bl.Data().Get("a.b")
		if got := string(d.Bytes()); exp != got {
			t.Errorf("expected %s error, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var exp = "World"

		d := bl.Data().Get("a.b")
		if got := string(d.Bytes()); exp != got {
			t.Errorf("expected %s error, got %s", exp, got)
		}
	}
}

func Test_blockgetValueAsOnSimpleArray(t *testing.T) {
	bl := newBlock("chars", 0, &Data{Value: []string{
		"a", "b",
	}})
	bl.As("char")

	{
		var (
			exp = "a"
			got = string(bl.Data().Get("char").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "b"
			got = string(bl.Data().Get("char").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockgetValueAsOnArrayOfObject(t *testing.T) {
	bl := newBlock("chars", 0, &Data{Value: []interface{}{
		map[string]interface{}{
			"value": "a",
		},
		map[string]interface{}{
			"value": "b",
		},
	}})
	bl.As("char")

	{
		var (
			exp = "a"
			got = string(bl.Data().Get("char.value").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "b"
			got = string(bl.Data().Get("char.value").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockgetValueAsKeyIsEqualToTag(t *testing.T) {
	bl := newBlock("char", 0, &Data{Value: []interface{}{
		map[string]interface{}{
			"value": "a",
		},
		map[string]interface{}{
			"value": "b",
		},
	}})
	bl.As("char")

	{
		var (
			exp = "a"
			got = string(bl.Data().Get("char.value").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "b"
			got = string(bl.Data().Get("char.value").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockAsGetKeyOnMap(t *testing.T) {
	bl := newBlock("char", 0, &Data{Value: map[string]interface{}{
		"a": "b",
		"c": "d",
		"e": "f",
	}})
	bl.As("k", "v")

	{
		var (
			exp = "a"
			got = string(bl.Data().Get("k").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "c"
			got = string(bl.Data().Get("k").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "e"
			got = string(bl.Data().Get("k").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockAsGetValueOnMap(t *testing.T) {
	bl := newBlock("char", 0, &Data{Value: map[string]interface{}{
		"a": "b",
		"c": "d",
		"e": "f",
	}})
	bl.As("k", "v")

	{
		var (
			exp = "b"
			got = string(bl.Data().Get("v").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "d"
			got = string(bl.Data().Get("v").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "f"
			got = string(bl.Data().Get("v").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockAsGetKeyOnMStruct(t *testing.T) {
	type s struct {
		a string
		c string
		e string
	}

	bl := newBlock("char", 0, &Data{Value: s{}})
	bl.As("k", "v")

	{
		var (
			exp = "a"
			got = string(bl.Data().Get("k").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "c"
			got = string(bl.Data().Get("k").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "e"
			got = string(bl.Data().Get("k").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}
