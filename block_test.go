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

func Test_blockgetValueDotOnSlice(t *testing.T) {
	bl := newBlock("", 0, []interface{}{
		"a", "b",
	})

	var exp = []string{
		"a", "b",
	}

	d := bl.getData(".")
	if got := d.(string); exp[bl.iterd] != got {
		t.Errorf("expected %s, got %s", exp[bl.iterd], got)
	}
	bl.increment()

	d = bl.getData(".")
	if got := d.(string); exp[bl.iterd] != got {
		t.Errorf("expected %s, got %s", exp[bl.iterd], got)
	}
	bl.increment()

	if !bl.isFinished() {
		t.Errorf("expected block to be finished")
	}
}
