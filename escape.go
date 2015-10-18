package beard

var escapeList = map[byte][]byte{
	// NOTE from https://golang.org/src/html/escape.go#L189
	'&':  []byte("&amp;"),
	'\'': []byte("&#39;"),
	'<':  []byte("&lt;"),
	'>':  []byte("&gt;"),
	'"':  []byte("&#34;"),

	// NOTE these are the base for the default beard delimiters
	'{': []byte("&#123;"),
	'}': []byte("&#125;"),
}

func escapeBytes(b []byte) []byte {
	lenb := len(b)
	i := 0
	for ; i < lenb; i++ {
		char := b[i]
		esc, ok := escapeList[char]
		if !ok {
			continue
		}

		lenesc := len(esc)

		// we'll normally get a default cap of 8, lets use as much of that
		// if we can before we make another allocation
		if capb := cap(b); capb < lenb+lenesc {
			c := capb * 2

			// check to see if the 2x cap is greater than the cap plus 2x the
			// longest escape, most escapes will come in pairs
			if min := 2*6 + capb; c < min {
				c = min
			}

			tmp := make([]byte, 0, c)
			tmp = append(tmp, b...)

			b = tmp
		}

		// replace the char with the first part of the escaped
		b[i] = esc[0]

		// grow, copy and insert
		for j := 1; j < lenesc; j++ {
			b = b[:len(b)+1]

			copy(b[i+j+1:], b[i+j:])
			b[i+j] = esc[j]
		}

		// recalculate to keep looping properly
		lenb = len(b)
		i = i + lenesc - 1
	}

	return b
}
