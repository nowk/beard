package beard

type block struct {
	tag    []byte
	cursor int
	data   *Data

	iterd int
}

func newBlock(tag []byte, c int, data *Data) *block {
	return &block{
		tag:    tag,
		cursor: c,
		data:   data,
	}
}

func (b *block) Data() *Data {
	if b.data.IsSlice() {
		// get data for current iteration context
		return b.data.Index(b.iterd)
	}

	return b.data
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
