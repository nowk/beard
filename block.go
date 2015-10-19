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

// Increment increments the iteration index. Non-slice types must call
// Increment() after it's been rendered.
func (b *block) Increment() int {
	b.iterd++

	return b.iterd
}

// Finished checks to see if a block has been completely iterated through.
func (b *block) Finished() bool {
	if b.Skip() {
		return true
	}
	if b.Empty() {
		return true
	}

	return !(b.iterd < b.data.Len())
}
