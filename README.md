# tmplx

A fast, pure-Go template engine with Mustache/Jinja2-like syntax. No external dependencies.

## Features

- **Familiar syntax**: `{{ var }}`, `{{ var | filter }}`, `{{# if }}`, `{{# each }}`, `{{> partial }}`
- **80+ built-in filters**: string manipulation, math, collections, encoding, type conversion, regex
- **Control flow**: if/else/unless, each loops, with blocks, nested blocks
- **Partials**: reusable template fragments
- **Raw blocks**: `{{{ raw }}}` for unprocessed content
- **Comments**: `{{! comment }}` stripped from output
- **Pure Go**: zero external dependencies, single binary
- **CLI tool**: render templates from files or inline, lint templates
- **Library API**: embeddable in any Go project

## Quick Start

### Install

```bash
go install github.com/EdgarOrtegaRamirez/tmplx/cmd/tmplx@latest
```

### CLI Usage

```bash
# Render inline template
tmplx eval 'Hello {{ name | upper }}!' --json '{"name": "world"}'
# Output: Hello WORLD!

# Render from file
tmplx render greeting.html --json '{"name": "Alice"}'

# List available filters
tmplx filters

# Lint a template for syntax errors
tmplx lint my-template.html
```

### Library Usage

```go
package main

import (
    "fmt"
    "github.com/EdgarOrtegaRamirez/tmplx"
)

func main() {
    e := tmplx.New()
    e.SetVariable("name", "World")
    e.SetVariable("items", []interface{}{"a", "b", "c"})
    
    result, err := e.Render("Hello {{ name }}! Items: {{ items | join(', ') }}")
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
    // Output: Hello World! Items: a, b, c
}
```

## Template Syntax

### Variables

```mustache
{{ name }}                  <!-- Simple variable -->
{{ user.name }}             <!-- Dotted path -->
{{ value | upper }}         <!-- With filter -->
{{ value | upper | truncate(10) }}  <!-- Chained filters -->
```

### Control Flow

```mustache
{{# if active }}
  Active
{{ else }}
  Inactive
{{/ if }}

{{# unless hidden }}
  Visible when not hidden
{{/ unless }}

{{# each items }}
  {{ item.name }} - {{ item.value }}
{{/ each }}

{{# with user }}
  {{ name }} ({{ email }})
{{/ with }}
```

### Partials

```mustache
{{> header }}
{{> sidebar }}
```

Register partials:
```go
e.AddPartial("header", "<h1>{{ title }}</h1>")
```

### Raw Blocks

```mustache
{{{ This is NOT processed as a template }}}
```

### Comments

```mustache
{{! This comment is stripped from output }}
```

## Built-in Filters

### String Filters
`upper`, `lower`, `capitalize`, `title`, `snake_case`, `camel_case`, `kebab_case`, `pascal_case`, `trim`, `ltrim`, `rtrim`, `strip`, `lstrip`, `rstrip`, `truncate(n)`, `truncate_words(n)`, `replace(old, new)`, `replace_first(old, new)`, `replace_last(old, new)`, `repeat(n)`, `pad_left(n, char)`, `pad_right(n, char)`, `center(n)`, `wrap(n)`, `reverse`, `contains(s)`, `starts_with(s)`, `ends_with(s)`, `split(sep)`, `words`, `join(sep)`, `html_escape`, `html_unescape`, `url_encode`, `url_decode`, `base64_encode`, `base64_decode`, `word_count`, `char_count`, `line_count`

### Regex Filters
`regex_match(pattern)`, `regex_find(pattern)`, `regex_find_all(pattern)`, `regex_replace(pattern, replacement)`, `regex_split(pattern)`

### Math Filters
`abs`, `round(n)`, `ceil`, `floor`, `min`, `max`, `sum`, `average`, `clamp(min, max)`, `percentage`, `format_number(n)`

### Collection Filters
`length`, `len`, `size`, `count`, `first`, `last`, `slice(start, end)`, `reverse_collection`, `flatten`, `compact`, `unique`, `uniq`, `sort`, `sort_by(key)`, `map(key)`, `select(key)`, `reject(key)`, `group_by(key)`, `join(sep)`

### Type Filters
`type_of`, `to_string`, `string`, `str`, `int`, `float`, `number`, `bool`, `is_nil`, `is_none`, `is_empty`, `is_blank`

### Utility Filters
`default(value)`, `not`, `ternary(true_val, false_val)`, `now`, `timestamp`, `to_json`, `to_yaml`, `to_csv`, `to_lines`, `keys`, `values`, `wrap(n)`

## Architecture

```
Source Template
     │
     ▼
┌─────────┐
│  Lexer   │  Tokenize input into tokens ({{, }}, identifiers, text, etc.)
└────┬────┘
     │
     ▼
┌─────────┐
│  Parser  │  Build AST (Abstract Syntax Tree) from tokens
└────┬────┘
     │
     ▼
┌───────────┐
│ Evaluator │  Walk AST nodes, resolve variables, apply filters
└────┬──────┘
     │
     ▼
Rendered Output
```

### Components

- **Lexer** (`lexer.go`): Tokenizes template syntax into a stream of typed tokens
- **Parser** (`parser.go`): Recursive descent parser that builds an AST from tokens
- **AST** (`ast.go`): Node types for all template constructs
- **Evaluator** (`evaluator.go`): Tree-walking interpreter that renders the AST
- **Filters** (`filters.go`): 80+ built-in filter functions
- **Engine** (`engine.go`): Public API that orchestrates the pipeline

## Performance

tmplx is a tree-walking interpreter (not bytecode-compiled), so it's best suited for:
- Development and prototyping
- Small to medium templates
- Embedded systems where simplicity matters
- Applications where Go compilation speed matters more than template execution speed

## License

MIT
