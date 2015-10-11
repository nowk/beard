package beard

import (
	"bytes"
)

var (
	ldelim = []byte("{{")
	rdelim = []byte("}}")
)

type matchLevel int

const (
	_ matchLevel = iota

	noMatch
	paMatch
	exMatch
)

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
