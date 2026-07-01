package tmplx

import (
	"fmt"
	"io"
	"os"
)

// Engine is the main template engine interface.
type Engine struct {
 evaluator *Evaluator
}

// New creates a new template engine.
func New() *Engine {
 return &Engine{
  evaluator: NewEvaluator(),
 }
}

// SetVariable sets a template variable.
func (e *Engine) SetVariable(name string, value interface{}) {
 e.evaluator.SetVariable(name, value)
}

// SetVariables sets multiple template variables.
func (e *Engine) SetVariables(vars map[string]interface{}) {
 e.evaluator.SetVariables(vars)
}

// SetPartial registers a partial template.
func (e *Engine) SetPartial(name, content string) {
 e.evaluator.SetPartial(name, content)
}

// RegisterFilter registers a custom filter.
func (e *Engine) RegisterFilter(name string, fn FilterFunc) {
 e.evaluator.RegisterFilter(name, fn)
}

// Render parses and renders a template string.
func (e *Engine) Render(template string) (string, error) {
 tokens, err := e.lex(template)
 if err != nil {
  return "", err
 }

 ast, err := e.parse(tokens)
 if err != nil {
  return "", err
 }

 return e.evaluator.Render(ast)
}

// RenderToFile parses, renders, and writes to a file.
func (e *Engine) RenderToFile(template, outputPath string) error {
 result, err := e.Render(template)
 if err != nil {
  return err
 }
 return os.WriteFile(outputPath, []byte(result), 0644)
}

// RenderFile reads a template file, parses, and renders it.
func (e *Engine) RenderFile(path string) (string, error) {
 data, err := os.ReadFile(path)
 if err != nil {
  return "", fmt.Errorf("read template: %w", err)
 }
 return e.Render(string(data))
}

// RenderWriter parses, renders, and writes to a writer.
func (e *Engine) RenderWriter(template string, w io.Writer) error {
 result, err := e.Render(template)
 if err != nil {
  return err
 }
 _, err = w.Write([]byte(result))
 return err
}

func (e *Engine) lex(template string) ([]Token, error) {
 lexer := NewLexer(template)
 return lexer.Tokenize()
}

func (e *Engine) parse(tokens []Token) (*TemplateNode, error) {
 parser := NewParser(tokens)
 return parser.Parse()
}

// MustRender is like Render but panics on error.
func MustRender(template string, vars map[string]interface{}) string {
 e := New()
 e.SetVariables(vars)
 result, err := e.Render(template)
 if err != nil {
  panic(err)
 }
 return result
}
