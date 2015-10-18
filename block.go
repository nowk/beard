package beard

type block struct {
	tag    string
	cursor int
	data   *Data

	iterd int
}

func newBlock(tag string, c int, data *Data) *block {
	return &block{
		tag:    tag,
		cursor: c,
		data:   data,
	}
}

func (b *block) Data() *Data {
	if b.Skip() {
		return nil
	}
	if b.data.IsSlice() {
		// get data for current iteration context
		return b.data.Index(b.iterd)
	}

	return b.data
}

func (b *block) Skip() bool {
	if b.Inverted() {
		return !b.Empty()
	}

	return b.Empty()
}

func (b *block) Inverted() bool {
	return b.tag != "" && b.tag[0] == '^'
}

func (b *block) Empty() bool {
	return b.data == nil || b.data.Len() == 0
}

// Increment increments and returns the current iterd
// All block types must explicitly call increment after they have been read
// through.
func (b *block) Increment() int {
	b.iterd++

	return b.iterd
}

// IsFinished checks to see if the block has finished and should be exited
// This assumes that increment has been explicitly called after each a block has
// been read through.
func (b *block) IsFinished() bool {
	if b.Skip() {
		return true
	}
	if b.Empty() {
		return true
	}

	return !(b.iterd < b.data.Len())
}
