package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/EdgarOrtegaRamirez/tmplx"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "render", "r":
		cmdRender()
	case "eval", "e":
		cmdEval()
	case "lint", "l":
		cmdLint()
	case "filters", "f":
		cmdFilters()
	case "version", "v":
		fmt.Printf("tmplx v%s\n", version)
	case "help", "h":
		printUsage()
	default:
		// Treat as template string
		cmdEvalString(cmd)
	}
}

func printUsage() {
	fmt.Println(`tmplx - A template engine with Mustache/Jinja2-like syntax

Usage:
  tmplx render <template-file> [--data <data-file>] [--json <json-string>]
  tmplx eval <template-string> [--json <json-string>]
  tmplx lint <template-file>
  tmplx filters
  tmplx version

Examples:
  tmplx render greeting.html --json '{"name": "World"}'
  tmplx eval 'Hello {{ name | upper }}!' --json '{"name": "world"}'
  tmplx lint my-template.html

Template Syntax:
  {{ variable }}           - Variable interpolation
  {{ variable | filter }}  - Apply filter
  {{# if condition }}      - Conditional
  {{ else }}               - Else branch
  {{/ if }}                - End conditional
  {{# each items }}        - Loop
  {{/ each }}                - End loop
  {{> partial_name }}      - Include partial
  {{! comment }}           - Comment

Filters:
  {{ name | upper }}       - Uppercase
  {{ name | lower }}       - Lowercase
  {{ name | truncate(20) }} - Truncate to 20 chars
  {{ items | join(", ") }} - Join list with separator
  {{ value | default("N/A") }} - Default value`)
}

func cmdRender() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: template file required")
		os.Exit(1)
	}

	templatePath := os.Args[2]
	data, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading template: %v\n", err)
		os.Exit(1)
	}

	vars := parseDataArgs()

	engine := tmplx.New()
	engine.SetVariables(vars)

	// Auto-load partials from same directory
	loadPartials(engine, templatePath)

	result, err := engine.Render(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering template: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

func cmdEval() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: template string required")
		os.Exit(1)
	}

	template := os.Args[2]
	vars := parseDataArgs()

	engine := tmplx.New()
	engine.SetVariables(vars)

	result, err := engine.Render(template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering template: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

func cmdLint() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Error: template file required")
		os.Exit(1)
	}

	templatePath := os.Args[2]
	data, err := os.ReadFile(templatePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading template: %v\n", err)
		os.Exit(1)
	}

	engine := tmplx.New()
	_, err = engine.Render(string(data))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lint error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Template is valid")
}

func cmdFilters() {
	registry := tmplx.NewFilterRegistry()
	filters := []string{
		"upper", "lower", "title", "capitalize", "trim", "ltrim", "rtrim",
		"truncate", "truncate_words", "replace", "replace_first", "replace_last",
		"repeat", "reverse", "pad_left", "pad_right", "center", "wrap",
		"snake_case", "camel_case", "kebab_case", "pascal_case",
		"html_escape", "html_unescape", "url_encode", "url_decode",
		"base64_encode", "base64_decode",
		"abs", "ceil", "floor", "round", "clamp", "min", "max", "sum", "average",
		"percentage", "format_number",
		"int", "float", "str", "bool", "string", "number",
		"length", "len", "size", "first", "last", "slice", "sort", "sort_by",
		"reverse_collection", "flatten", "unique", "uniq", "compact",
		"join", "map", "filter", "select", "reject", "group_by", "count",
		"keys", "values", "to_json", "to_yaml", "to_csv", "to_lines",
		"default", "d", "ternary", "not", "is_empty", "is_blank", "is_nil", "is_none",
		"regex_match", "regex_find", "regex_find_all", "regex_replace", "regex_split",
		"contains", "starts_with", "ends_with", "words", "lines",
		"char_count", "word_count", "line_count", "type_of", "now", "timestamp",
	}

	fmt.Printf("Available filters (%d):\n", len(filters))
	for _, name := range filters {
		_, ok := registry.Get(name)
		status := "✓"
		if !ok {
			status = "✗"
		}
		fmt.Printf("  %s %s\n", status, name)
	}
}

func cmdEvalString(template string) {
	vars := parseDataArgs()

	engine := tmplx.New()
	engine.SetVariables(vars)

	result, err := engine.Render(template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering template: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(result)
}

func parseDataArgs() map[string]interface{} {
	vars := make(map[string]interface{})

	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--json", "-j":
			if i+1 < len(os.Args) {
				i++
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(os.Args[i]), &data); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
					os.Exit(1)
				}
				for k, v := range data {
					vars[k] = v
				}
			}
		case "--data", "-d":
			if i+1 < len(os.Args) {
				i++
				data, err := os.ReadFile(os.Args[i])
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading data file: %v\n", err)
					os.Exit(1)
				}
				var vars2 map[string]interface{}
				if err := json.Unmarshal(data, &vars2); err != nil {
					fmt.Fprintf(os.Stderr, "Error parsing data file: %v\n", err)
					os.Exit(1)
				}
				for k, v := range vars2 {
					vars[k] = v
				}
			}
		case "--set", "-s":
			if i+2 < len(os.Args) {
				i++
				key := os.Args[i]
				i++
				val := os.Args[i]
				vars[key] = val
			}
		}
	}

	return vars
}

func loadPartials(engine *tmplx.Engine, templatePath string) {
	dir := templatePath[:strings.LastIndex(templatePath, "/")]
	if dir == "" {
		dir = "."
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".partial") || strings.HasPrefix(name, "_") {
			partialPath := dir + "/" + name
			data, err := os.ReadFile(partialPath)
			if err != nil {
				continue
			}
			partialName := strings.TrimSuffix(strings.TrimPrefix(name, "_"), ".partial")
			engine.SetPartial(partialName, string(data))
		}
	}
}

func init() {
	// Ensure we read from stdin if needed
	if info, err := os.Stdin.Stat(); err == nil && (info.Mode()&os.ModeCharDevice == 0) {
		// stdin has data, could be used for template input
		_, _ = io.ReadAll(os.Stdin) // Just check it's readable
	}
}
