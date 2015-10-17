package beard

import (
	"bytes"
	"io"
	"testing"
)

func BenchmarkBasicVar(b *testing.B) {
	html := `<h1>Hello {{c}}</h1>`
	data := map[string]interface{}{
		"c": "world!",
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		tmpl.File.Seek(0, 0)
		tmpl.blocks = tmpl.blocks[:0]
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArray(b *testing.B) {
	html := `{{#words}}({{.}}){{/words}}`
	data := map[string]interface{}{
		"words": []string{
			"a", "b", "c",
		},
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		tmpl.File.Seek(0, 0)
		tmpl.blocks = tmpl.blocks[:0]
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArrayInArray(b *testing.B) {
	html := `{{#words}}({{.}}){{#words}}({{.}}){{#words}}({{.}}){{/words}}{{/words}}{{/words}}`
	data := map[string]interface{}{
		"words": []string{
			"a", "b", "c",
		},
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		tmpl.File.Seek(0, 0)
		tmpl.blocks = tmpl.blocks[:0]
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBasicBlock(b *testing.B) {
	html := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/greeting}}</h1>`
	data := map[string]interface{}{
		"greeting": map[string]interface{}{
			"a": "Hello",
			"b": "World",
			"c": map[string]interface{}{
				"d": "!",
			},
		},
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		tmpl.File.Seek(0, 0)
		tmpl.blocks = tmpl.blocks[:0]
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBlockWithOutsideVar(b *testing.B) {
	html := `<h1>{{#greeting}}{{a}} {{b}}{{c.d}}{{/greeting}}</h1>`
	data := map[string]interface{}{
		"a": "Hello",
		"greeting": map[string]interface{}{
			"b": "World",
			"c": map[string]interface{}{
				"d": "!",
			},
		},
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		tmpl.File.Seek(0, 0)
		tmpl.blocks = tmpl.blocks[:0]
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEscape(b *testing.B) {
	html := `<code>{{c}}{{c}}{{c}}</code>`
	data := map[string]interface{}{
		"c": "<h1>Hello World!</h1>",
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		tmpl.File.Seek(0, 0)
		tmpl.blocks = tmpl.blocks[:0]
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBasicVar                 500000              2954 ns/op             192 B/op          9 allocs/op
// BenchmarkArray                    300000              4698 ns/op             344 B/op         17 allocs/op
// BenchmarkArrayInArray             200000             10862 ns/op            1016 B/op         47 allocs/op
// BenchmarkBasicBlock               200000              9483 ns/op             864 B/op         40 allocs/op
// BenchmarkBlockWithOutsideVar      200000              9945 ns/op             912 B/op         42 allocs/op
// BenchmarkEscape                   200000              9766 ns/op             864 B/op         29 allocs/op
