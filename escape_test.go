package beard

import (
	"testing"
)

func Test_escapeBytes(t *testing.T) {
	for _, v := range []struct {
		giv string
		exp string
	}{
		{"<h1>", "&lt;h1&gt;"},
		{"</h1>", "&lt;/h1&gt;"},
		{"<h1>{{c}}</h1>", "&lt;h1&gt;&#123;&#123;c&#125;&#125;&lt;/h1&gt;"},
	} {
		b := escapeBytes([]byte(v.giv))
		if got := string(b); v.exp != got {
			t.Errorf("expected %s, got %s", v.exp, got)
		}
	}
}
