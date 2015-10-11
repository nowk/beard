package beard

import (
	"reflect"
)

type block struct {
	name   string
	cursor int

	data *blockData

	iterd int
}

func newBlock(name string, c int, data interface{}) *block {
	var bd *blockData
	var ok bool
	bd, ok = data.(*blockData)
	if !ok {
		bd = &blockData{
			reflect.ValueOf(data),
		}
	}

	bl := &block{
		name:   name,
		cursor: c,

		data: bd,
	}

	return bl
}

// increment increments and returns the current iterd
// All block types must explicitly call increment after they have been read
// through.
func (b *block) increment() int {
	b.iterd++

	return b.iterd
}

// isFinished checks to see if the block has finished and should be exited
// This assumes that increment has been explicitly called after each a block has
// been read through.
func (b *block) isFinished() bool {
	return !(b.iterd < b.data.Len())
}

func (b *block) getData(k string) interface{} {
	if k == "." {
		if b.data.isSlice() {
			return b.data.Index(b.iterd).Interface()
		}
	}

	return nil
}

type blockData struct {
	reflect.Value
}

func (d *blockData) Len() int {
	if d.isSlice() {
		return d.Value.Len()
	}

	return 0
}

func (d *blockData) isSlice() bool {
	return d.Kind() == reflect.Slice
}
