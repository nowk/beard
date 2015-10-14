package beard

import (
	"testing"
)

func Test_blockisFinished(t *testing.T) {
	bl := newBlock("", 0, []string{
		"a", "b",
	})

	if bl.isFinished() {
		t.Errorf("expected block to not be finished")
	}

	bl.increment()
	bl.increment()

	if !bl.isFinished() {
		t.Errorf("expected block to be finished")
	}
}

func Test_blockgetvofDotOnSlice(t *testing.T) {
	bl := newBlock("", 0, []interface{}{
		"a",
		"b",
	})

	{
		var exp = "a"

		d := bl.getvof(".")
		if got := d.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}

	bl.increment()

	{
		var exp = "b"

		d := bl.getvof(".")
		if got := d.(string); exp != got {
			t.Errorf("expected %s, got %s", exp, got)
		}
	}
}
