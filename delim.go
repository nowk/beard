package beard

import (
	"bytes"
)

type matchLevel int

const (
	_ matchLevel = iota

	noMatch
	paMatch
	exMatch
)

type Delim interface {
	Match([]byte) ([]byte, matchLevel)

	// Value should return the actual Delimeter value in []byte, eg {{ or }}
	Value() []byte
}

var (
	ldelim = &Ldelim{Delim: []byte("{{")}
	rdelim = &Rdelim{Delim: []byte("}}")}
)

type Ldelim struct {
	Delim []byte
}

var _ Delim = &Ldelim{}

func (d *Ldelim) Match(b []byte) ([]byte, matchLevel) {
	return matchDelim(b, d.Delim)
}

func (d *Ldelim) Value() []byte {
	return d.Delim
}

type Rdelim struct {
	Delim []byte
}

var _ Delim = &Rdelim{}

func (d *Rdelim) Match(b []byte) ([]byte, matchLevel) {
	return matchDelim(b, d.Delim)
}

func (d *Rdelim) Value() []byte {
	return d.Delim
}

func matchDelim(b, del []byte) ([]byte, matchLevel) {
	lenb := len(b)
	lend := len(del)

	// find the delim in full
	if i := bytes.Index(b, del); i != -1 {
		return b[:i+lend], exMatch
	}

	// find a partial match
	z := lend - 1
	for ; z > 0; z-- {
		i := bytes.Index(b, del[:z])
		if i == -1 {
			continue
		}
		// match must be at the end of the byte array
		if i+z == lenb {
			return b, paMatch
		}
	}

	return b, noMatch
}
