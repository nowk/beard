# beard

Mustache inspired templates in Go.

## Install

	go get github.com/nowk/beard

## Usage


	data := map[string]interface{}{
		"a": "Hello",
		"b": "World",
		"c": "!",
	}

	tmpl := beard.Render(
		bytes.NewReader([]byte(`<h1>{{a}} {{b}}{{c}}</h1>`)), 
		data,
		nil, // we expect no partials
	)

	w := bytes.NewBuffer(nil)
	
	_, err := io.Copy(w, tmpl)
	if err != nil {
		// handle
	}

*When passing in a `os.File` or other object that must be `Close()`, it is the users job to ensure those descriptors are closed.*

---

#### Variables

Variables can be a single node in the data tree, or a path to an inner value.

Template:

	{{name}}
	{{title}}
	{{car.model}}

Data:

	map[string]interface{}{
		"name":  "Batman",
		"title": "The bat/man",
		"car":   map[string]interface{}{
			"model": "The Bat Mobile",
		},
	}

*Data may also contain structs.*

Output:

	Batman
	The bat/man
	The Bat Mobile

---

Variables are escaped by default. You can use `&` to unescape a variable.

Template:

	{{battlecry}}
	{{&battlecry}}

Data:

	map[string]interface{}{
		"battlecry": "I'm <strong>bat and man!</strong>",
	}
	
Output:

	I'm &lt;strong&gt;bat and man!&lt;/strong&gt;
	I'm <strong>bat and man!</strong>

---

#### Blocks

Template:

	{{#profile}}
		{{name}}
		{{title}}
		{{car.model}}
	{{/profile}}
	
Data:

	map[string]interface{}{
		"profile" map[string]interface{}{
			"name":  "Batman",
			"title": "The bat/man",
			"car":   map[string]interface{}{
				"model": "The Bat Mobile",
			},
		},
	}
	
Output:

	Batman
	The bat/man
	The Bat Mobile

---

Non-empty lists.

Template:

	{{#enemies}}
		{{.}}
	{{/enemies}}

Data:

	map[string]interface{}{
		"enemies": []string{
			"joker",
			"penguin",
			"the governator",
		},
	}

Output:

	joker
	penguin
	the governator

---

Lists can also be a lists of objects.

Template:

	{{#enemies}}
		{{name}}
	{{/enemies}}

Data:

	map[string]interface{}{
		"enemies": []struct{
			name string
		} {
			{name: "joker"},
			{name: "penguin"},
			{name: "the governator"},
		},
	}

Output:

	joker
	penguin
	the governator

---

Variables not found within the block will look outside the block's data scope  in an attempt to find a matching path.

Template:
	
	{{#car}}
		{{name}}
		{{title}}
		{{model}}
	{{/car}}

Data:

	map[string]interface{}{
		"name":  "Batman",
		"title": "The bat/man",
		"car":   map[string]interface{}{
			"model": "The Bat Mobile",
		},
	}
	
Output:

	Batman
	The bat/man
	The Bat Mobile

---

When looking up variables, it will always use the closest match.

Template:
	
	{{model}}
	{{#car}}
		{{model}}
	{{/car}}

Data:

	map[string]interface{}{
		"model": "Cindy Crawford",
		"car":   map[string]interface{}{
			"model": "The Bat Mobile",
		},
	}
	
Output:

	Cindy Crawford
	The Bat Mobile

#### Inverted blocks

Template:

	{{#enemies}}
		{{.}}
	{{/enemies}}
	{{^enemies}}
		We all love this dude!
	{{/enemies}}

*Blocks that are empty, nil or undefined will not render their inner content.*

Data:

	map[string]interface{}{
		"enemies": []interface{}{},
	}

	// or
	
	map[string]interface{}{
		"enemies": nil,
	}

	// or
	
	map[string]interface{}{}

Output:

	We all love this dude!

---

#### Partials

Partials require the user to define a `PartialFunc` to return the partial file. 

`PartialFunc` gives the user the ability to define their own lookup logic to find partial files.

	type PartialFunc func(string) (io.Reader, error)

The `string` argument is the value in the partial definition

Template:

	{{>shared/file}}

`string` argument:

	shared/file

---

	tmpl := beard.Render(
		bytes.NewReader([]byte(`<h1>{{a}} {{b}}{{c}}</h1>`)), 
		data,
		func(path string) (io.Reader, error) {
			// use path to look up partial file your way and 
			// return a io.Reader
		},
	)

*Partials will inherit all data from their parent template.*

*The defined `PartialFunc` is inherited throughout the templates partial chain and will be used when rendering partials within other partials.*

*If the `io.Reader` returned implements `io.ReadCloser`, beard will close those descriptors.*

---

Template:

	{{>a}}

Partial: 

	{{b}}

Data:

	map[string]interface{}{
		"b": "c",
	}

Output:

	c

---

Partials within blocks will honor the block's variable look up rules.

Template:

	{{#profile}}
		{{>partialFile}}
	{{/profile}}

Partial:

	{{model}}
	{{name}}
	{{title}}

Data:

	map[string]interface{}{
		"model":   "Cindy Crawford",
		"profile": map[string]interface{}{
			"name":  "Batman",
			"title": "The bat/man",
		},
	}

Output:

	Cindy Crawford
	Batman
	The bat/man

---

#### Layouts
Beard provides a  utility method to aid in rendering a template within a layout. 

	tmpl := beard.RenderInLayout(
		bytes.NewReader([]byte(`<body>{{>yield}}</body>`)),
		bytes.NewReader([]byte(`<h1>Hello {{c}}!</h1>`)),
		data,
		func(path string) (io.Reader, error) {
			//
		},
	)

*The `PartialFunc` defined here is applied to the inner template and not to the "layout".*

*The "layout" file will have full access to the `data`.*

*If any of the readers need to be `Close()` it is the user's job to ensure those descriptors are closed.*

---

Layouts, under the hood, leverage the existing partial syntax and use a special partial definition.

	{{>yield}}

*This is specific to `RenderInLayout` only.*

## TODO

- [ ] `func` support
- [ ] named blocks 

		{{#words each word}}
			{{word}}
		{{/words}}

		{{#word as word}}
			{{word.inEnglish}}
		{{/word}}

- [ ] a simple way to handle condition logic


## License

MIT
