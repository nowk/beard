package beard

import (
	"testing"
)

func Test_blockisFinished(t *testing.T) {
	bl := newBlock(nil, 0, &Data{value: []string{
		"a", "b",
	}})

	if bl.isFinished() {
		t.Errorf("expected block to not be finished")
	}

	bl.increment()
	bl.increment()

	if !bl.isFinished() {
		t.Errorf("expected block to be finished")
	}
}

func Test_blockgetValueDotOnSlice(t *testing.T) {
	bl := newBlock(nil, 0, &Data{value: []interface{}{
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

	bl.increment()

	{
		var exp = "b"

		d := bl.Data().Get(".")
		if got := string(d.Bytes()); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}

func Test_blockgetValuePathOnSlice(t *testing.T) {
	bl := newBlock(nil, 0, &Data{value: []interface{}{
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

	bl.increment()

	{
		var exp = "World"

		d := bl.Data().Get("a.b")
		if got := string(d.Bytes()); exp != got {
			t.Errorf("expected %s error, got %s", exp, got)
		}
	}
}
