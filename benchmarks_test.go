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
		tmpl.cursor = 0
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
		tmpl.cursor = 0
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
		tmpl.cursor = 0
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
		tmpl.cursor = 0
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
		tmpl.cursor = 0
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
		tmpl.cursor = 0
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// NOTE this does not take into account file lookups or things of that nature
// those will vary based on your usage.
func BenchmarkPartialInPartial(b *testing.B) {
	html := `{{>a}}{{>c}}{{f}}`
	data := map[string]interface{}{
		"b": "Hello",
		"e": "World",
		"f": "!",
	}

	tmpl := &Template{
		File: bytes.NewReader([]byte(html)),
		Data: &Data{Value: data},
	}
	tmpl.Partial(func(path string) (io.Reader, error) {
		var p []byte
		switch path {
		case "a":
			p = []byte(`{{b}}`)
		case "c":
			p = []byte(` {{>d}}`)
		case "d":
			p = []byte(`{{e}}`)

		default:
			b.Fatalf("invalid partial %s", path)
		}

		return bytes.NewReader(p), nil
	})

	buf := bytes.NewBuffer(nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		tmpl.File.Seek(0, 0)
		tmpl.blocks = tmpl.blocks[:0]
		tmpl.cursor = 0
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBasicVar                 500000              2579 ns/op             192 B/op          9 allocs/op
// BenchmarkArray                    200000              9138 ns/op             600 B/op         35 allocs/op
// BenchmarkArrayInArray              10000            103828 ns/op            7817 B/op        434 allocs/op
// BenchmarkBasicBlock               200000             11042 ns/op             864 B/op         40 allocs/op
// BenchmarkBlockWithOutsideVar      200000              9416 ns/op             912 B/op         42 allocs/op
// BenchmarkEscape                   200000              9318 ns/op             864 B/op         29 allocs/op
// BenchmarkPartialInPartial         200000             11860 ns/op            1176 B/op         41 allocs/op
