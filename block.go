package beard

import (
	"reflect"
)

type block struct {
	name   []byte
	cursor int

	data *blockData

	iterd int
}

func newBlock(name []byte, c int, data interface{}) *block {
	bl := &block{
		name:   name,
		cursor: c,

		data: &blockData{
			reflect.ValueOf(data),
		},
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

func (b *block) getvof(k string) interface{} {
	if b.data.isSlice() {
		v := b.data.Index(b.iterd).Interface()
		// . returns the value itself
		if k == "." {
			return v
		}

		return getvof(k, v)
	}

	return getvof(k, b.data.Value)
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
