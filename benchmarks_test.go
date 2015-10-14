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

func BenchmarkBasicBlock(b *testing.B) {
	tmpl := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/greeting}}</h1>`
	data := map[string]interface{}{
		"greeting": map[string]interface{}{
			"a": "Hello",
			"b": "World",
			"c": map[string]interface{}{
				"d": "!",
			},
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

func BenchmarkBlockWithOutsideVar(b *testing.B) {
	tmpl := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/greeting}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"greeting": map[string]interface{}{
			"b": "World",
			"c": map[string]interface{}{
				"d": "!",
			},
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

// BenchmarkBasicVar                1000000              2389 ns/op             160 B/op          8 allocs/op
// BenchmarkArray                    500000              3004 ns/op             149 B/op          9 allocs/op
// BenchmarkArrayInArray              20000            130202 ns/op             722 B/op         33 allocs/op
// BenchmarkBasicBlock               200000              9093 ns/op             816 B/op         38 allocs/op
// BenchmarkBlockWithOutsideVar      200000              9864 ns/op             864 B/op         40 allocs/op
