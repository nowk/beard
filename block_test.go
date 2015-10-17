package beard

import (
	"testing"
)

func Test_blockisFinished(t *testing.T) {
	bl := newBlock("", 0, &Data{Value: []string{
		"a", "b",
	}})

	if bl.IsFinished() {
		t.Errorf("expected block to not be finished")
	}

	bl.Increment()
	bl.Increment()

	if !bl.IsFinished() {
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
