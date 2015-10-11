package beard

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkBasicVar(b *testing.B) {
	tmpl := `<h1>Hello {{c}}</h1>`
	data := map[string]interface{}{
		"c": "world!",
	}

	rend := &Renderable{
		File: bytes.NewReader([]byte(tmpl)),
		Data: data,
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		rend.File.Seek(0, 0)
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, rend)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArray(b *testing.B) {
	tmpl := `{{#words}}({{.}}){{/words}}`
	data := map[string]interface{}{
		"words": []string{
			"a", "b", "c",
		},
	}

	rend := &Renderable{
		File: bytes.NewReader([]byte(tmpl)),
		Data: data,
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		rend.File.Seek(0, 0)
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, rend)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArrayInArray(b *testing.B) {
	tmpl := `{{#words}}({{.}}){{#words}}({{.}}){{#words}}({{.}}){{/words}}{{/words}}{{/words}}`
	data := map[string]interface{}{
		"words": []string{
			"a", "b", "c",
		},
	}

	rend := &Renderable{
		File: bytes.NewReader([]byte(tmpl)),
		Data: data,
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		rend.File.Seek(0, 0)
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, rend)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBasicVar        1000000              1419 ns/op              64 B/op          4 allocs/op
// BenchmarkArray            500000              2590 ns/op             130 B/op          9 allocs/op
// BenchmarkArrayInArray      10000            244304 ns/op             542 B/op         27 allocs/op
