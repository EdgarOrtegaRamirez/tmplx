# AGENTS.md

## Project: tmplx

A pure-Go template engine with Mustache/Jinja2-like syntax.

## Architecture

- **Lexer** (`lexer.go`): Tokenizer that converts template strings into token streams
- **Parser** (`parser.go`): Recursive descent parser building AST from tokens
- **AST** (`ast.go`): Abstract syntax tree node definitions
- **Evaluator** (`evaluator.go`): Tree-walking evaluator that renders AST to string output
- **Filters** (`filters.go`): 80+ built-in filter functions registered in `FilterRegistry`
- **Engine** (`engine.go`): Public API: `New()`, `Render()`, `SetVariable()`, etc.
- **CLI** (`cmd/tmplx/main.go`): Command-line interface for rendering, linting, and listing filters

## Building

```bash
go build -o tmplx ./cmd/tmplx
```

## Testing

```bash
go test -v ./...
```

## Key Design Decisions

1. **Pure Go, zero dependencies** - no external packages, easy to vendor
2. **Tree-walking evaluator** - simpler than bytecode, good for most use cases
3. **Filter chaining** - filters can be piped: `{{ name | upper | truncate(10) }}`
4. **Mustache-like syntax** - familiar to most developers
5. **Type coercion** - filters auto-convert types (string/int/float/bool)
6. **Partials** - registered by name, not file system lookup (security)

## Token Types

- `TokenText` - literal text between expressions
- `TokenVariableStart/End` - `{{` and `}}` delimiters
- `TokenBlockStart/End` - `{{#` and `{{/` block delimiters
- `TokenIdent` - variable/filter names
- `TokenFilter` - pipe `|` character
- `TokenElse` - `{{ else }}`
- `TokenPartialStart` - `{{>`
- `TokenRawStart/End` - `{{{` and `}}}`
- `TokenCommentStart/End` - `{{!` and `}}`

## Adding New Filters

1. Implement filter function in `filters.go`:
   ```go
   func filterMyFilter(value interface{}, args []interface{}) (interface{}, error) {
       // implementation
   }
   ```
2. Register in `NewFilterRegistry()`:
   ```go
   r.Register("my_filter", filterMyFilter)
   ```
3. Add tests in `evaluator_test.go`

## Common Issues

- Filter names use `snake_case` (not `camelCase`)
- `truncate(n)` includes the "..." suffix in the count
- `group_by` returns `map[string][]interface{}` which may not chain with `keys`
- Empty/nil values are displayed as empty string
