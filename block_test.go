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

func Test_blockAsgetValueOnArray(t *testing.T) {
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

func Test_blockAsgetValueWithAsAsRootPath(t *testing.T) {
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

func Test_blockAsgetValueWhereAsAsRootPathEqualsTag(t *testing.T) {
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

func Test_blockAsKeyValuegetKeyOnMap(t *testing.T) {
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

func Test_blockAsKeyValuegetValueOnMap(t *testing.T) {
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

func Test_blockAsKeyValuegetKeyOnStruct(t *testing.T) {
	type s struct {
		a string
		c string
		e string
	}

	for _, v := range []interface{}{
		s{},
		&s{},
	} {
		bl := newBlock("char", 0, &Data{Value: v})
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
}

func Test_blockAsKeyValuegetValueOnStruct(t *testing.T) {
	type s struct {
		a string
		c string
		e string
	}

	for _, v := range []interface{}{
		s{"b", "d", "f"},
		&s{"b", "d", "f"},
	} {
		bl := newBlock("char", 0, &Data{Value: v})
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
}

func Test_blockAsKeyValueOnArrayKeyReturnsIndex(t *testing.T) {
	bl := newBlock("", 0, &Data{Value: []interface{}{
		"a", "b", "c",
	}})
	bl.As("i", "v")

	{
		var (
			exp = "0"
			got = string(bl.Data().Get("i").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
	{
		var (
			exp = "a"
			got = string(bl.Data().Get("v").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "1"
			got = string(bl.Data().Get("i").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
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
			exp = "2"
			got = string(bl.Data().Get("i").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
	{
		var (
			exp = "c"
			got = string(bl.Data().Get("v").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockAsKeyValueOnArrayOfMapsKeyReturnsIndex(t *testing.T) {
	bl := newBlock("", 0, &Data{Value: []interface{}{
		map[string]interface{}{
			"value": "a",
		},
		map[string]interface{}{
			"value": "b",
		},
		map[string]interface{}{
			"value": "c",
		},
	}})
	bl.As("i", "v")

	{
		var (
			exp = "0"
			got = string(bl.Data().Get("i").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
	{
		var (
			exp = "a"
			got = string(bl.Data().Get("v.value").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "1"
			got = string(bl.Data().Get("i").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
	{
		var (
			exp = "b"
			got = string(bl.Data().Get("v.value").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.Increment()

	{
		var (
			exp = "2"
			got = string(bl.Data().Get("i").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
	{
		var (
			exp = "c"
			got = string(bl.Data().Get("v.value").Bytes())
		)
		if exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}
