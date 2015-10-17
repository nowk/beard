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
	tmpl.Partial(func(path string) (File, error) {
		var p []byte
		switch path {
		case "a":
			p = []byte(`{{b}}`)
		case "c":
			p = []byte(` {{>d}}`)
		case "d":
			p = []byte(`{{e}}`)
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
		buf.Reset()

		b.StartTimer()

		_, err := io.Copy(buf, tmpl)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkBasicVar                 500000              2978 ns/op             192 B/op          9 allocs/op
// BenchmarkArray                    300000              4084 ns/op             344 B/op         17 allocs/op
// BenchmarkArrayInArray             100000             11193 ns/op            1016 B/op         47 allocs/op
// BenchmarkBasicBlock               200000              9768 ns/op             864 B/op         40 allocs/op
// BenchmarkBlockWithOutsideVar      200000             10172 ns/op             912 B/op         42 allocs/op
// BenchmarkEscape                   200000              9208 ns/op             864 B/op         29 allocs/op
// BenchmarkPartialInPartial         100000             11259 ns/op            1128 B/op         41 allocs/op
