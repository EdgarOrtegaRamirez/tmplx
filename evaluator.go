package tmplx

import (
	"fmt"
	"strings"
)

// Evaluator walks the AST and produces output.
type Evaluator struct {
	filters   *FilterRegistry
	partials  map[string]string
	variables map[string]interface{}
}

// NewEvaluator creates a new evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		filters:   NewFilterRegistry(),
		partials:  make(map[string]string),
		variables: make(map[string]interface{}),
	}
}

// SetVariable sets a template variable.
func (e *Evaluator) SetVariable(name string, value interface{}) {
	e.variables[name] = value
}

// SetVariables sets multiple template variables.
func (e *Evaluator) SetVariables(vars map[string]interface{}) {
	for k, v := range vars {
		e.variables[k] = v
	}
}

// SetPartial registers a partial template.
func (e *Evaluator) SetPartial(name, content string) {
	e.partials[name] = content
}

// RegisterFilter registers a custom filter.
func (e *Evaluator) RegisterFilter(name string, fn FilterFunc) {
	e.filters.Register(name, fn)
}

// Render renders a template AST to a string.
func (e *Evaluator) Render(template *TemplateNode) (string, error) {
	var sb strings.Builder
	for _, child := range template.Children {
		output, err := e.renderNode(child)
		if err != nil {
			return "", err
		}
		sb.WriteString(output)
	}
	return sb.String(), nil
}

func (e *Evaluator) renderNode(node Node) (string, error) {
	if node == nil {
		return "", nil
	}

	switch n := node.(type) {
	case *TextNode:
		return n.Text, nil

	case *VariableNode:
		return e.renderVariable(n)

	case *RawNode:
		output, err := e.renderNode(n.Content)
		if err != nil {
			return "", err
		}
		return output, nil

	case *IfNode:
		return e.renderIf(n)

	case *UnlessNode:
		return e.renderUnless(n)

	case *EachNode:
		return e.renderEach(n)

	case *WithNode:
		return e.renderWith(n)

	case *PartialNode:
		return e.renderPartial(n)

	case *CommentNode:
		return "", nil

	case *TemplateNode:
		var sb strings.Builder
		for _, child := range n.Children {
			output, err := e.renderNode(child)
			if err != nil {
				return "", err
			}
			sb.WriteString(output)
		}
		return sb.String(), nil

	default:
		return "", fmt.Errorf("unknown node type: %T", node)
	}
}

func (e *Evaluator) renderVariable(node *VariableNode) (string, error) {
	value, err := e.resolveVariable(node.Name)
	if err != nil {
		return "", fmt.Errorf("variable %q: %w", node.Name, err)
	}

	// Apply filters
	for _, filter := range node.Filters {
		fn, ok := e.filters.Get(filter.Name)
		if !ok {
			return "", fmt.Errorf("unknown filter: %s", filter.Name)
		}

		// Resolve filter arguments
		args := make([]interface{}, len(filter.Args))
		for i, arg := range filter.Args {
			switch a := arg.(type) {
			case *VariableNode:
				resolved, err := e.resolveVariable(a.Name)
				if err != nil {
					return "", fmt.Errorf("filter arg %q: %w", a.Name, err)
				}
				args[i] = resolved
			case *StringLiteralNode:
				args[i] = a.Value
			case *IntLiteralNode:
				args[i] = a.Value
			default:
				args[i] = arg
			}
		}

		value, err = fn(value, args)
		if err != nil {
			return "", fmt.Errorf("filter %s: %w", filter.Name, err)
		}
	}

	return toStringVal(value), nil
}

func (e *Evaluator) resolveVariable(name string) (interface{}, error) {
	parts := strings.Split(name, ".")
	var current interface{} = e.variables

	for _, part := range parts {
		if current == nil {
			return nil, nil
		}

		switch v := current.(type) {
		case map[string]interface{}:
			val, ok := v[part]
			if !ok {
				return nil, nil
			}
			current = val
		case map[interface{}]interface{}:
			val, ok := v[part]
			if !ok {
				return nil, nil
			}
			current = val
		default:
			return nil, fmt.Errorf("cannot access %q on %T", part, current)
		}
	}

	return current, nil
}

func (e *Evaluator) renderIf(node *IfNode) (string, error) {
	condValue, err := e.resolveVariable(node.Condition.(*VariableNode).Name)
	if err != nil {
		return "", err
	}

	if isTruthy(condValue) {
		var sb strings.Builder
		for _, child := range node.Body {
			output, err := e.renderNode(child)
			if err != nil {
				return "", err
			}
			sb.WriteString(output)
		}
		return sb.String(), nil
	}

	// Render else body
	var sb strings.Builder
	for _, child := range node.ElseBody {
		output, err := e.renderNode(child)
		if err != nil {
			return "", err
		}
		sb.WriteString(output)
	}
	return sb.String(), nil
}

func (e *Evaluator) renderUnless(node *UnlessNode) (string, error) {
	condValue, err := e.resolveVariable(node.Condition.(*VariableNode).Name)
	if err != nil {
		return "", err
	}

	if !isTruthy(condValue) {
		var sb strings.Builder
		for _, child := range node.Body {
			output, err := e.renderNode(child)
			if err != nil {
				return "", err
			}
			sb.WriteString(output)
		}
		return sb.String(), nil
	}

	return "", nil
}

func (e *Evaluator) renderEach(node *EachNode) (string, error) {
	iterValue, err := e.resolveVariable(node.Iterable.(*VariableNode).Name)
	if err != nil {
		return "", err
	}

	slice, ok := toSlice(iterValue)
	if !ok {
		// Try to iterate over map
		m, ok := toMap(iterValue)
		if !ok {
			return "", nil
		}
		// Convert map to slice of maps for iteration
		var items []interface{}
		for k, v := range m {
			item := map[string]interface{}{
				"_key":   k,
				"_value": v,
			}
			items = append(items, item)
		}
		slice = items
	}

	// Save current scope
	origVars := make(map[string]interface{})
	for k, v := range e.variables {
		origVars[k] = v
	}

	var sb strings.Builder
	for _, item := range slice {
		// Set loop variable
		e.variables["item"] = item
		e.variables["_"] = item

		// If there's a key variable, set it
		if node.KeyVar != "" {
			if m, ok := toMap(item); ok {
				e.variables[node.KeyVar] = m["_key"]
			}
		}

		for _, child := range node.Body {
			output, err := e.renderNode(child)
			if err != nil {
				return "", err
			}
			sb.WriteString(output)
		}
	}

	// Restore scope
	e.variables = origVars

	return sb.String(), nil
}

func (e *Evaluator) renderWith(node *WithNode) (string, error) {
	varValue, err := e.resolveVariable(node.Variable.(*VariableNode).Name)
	if err != nil {
		return "", err
	}

	// Save current scope
	origVars := make(map[string]interface{})
	for k, v := range e.variables {
		origVars[k] = v
	}

	// Set the variable in scope
	if m, ok := toMap(varValue); ok {
		for k, v := range m {
			e.variables[k] = v
		}
	}

	var sb strings.Builder
	for _, child := range node.Body {
		output, err := e.renderNode(child)
		if err != nil {
			return "", err
		}
		sb.WriteString(output)
	}

	// Restore scope
	e.variables = origVars

	return sb.String(), nil
}

func (e *Evaluator) renderPartial(node *PartialNode) (string, error) {
	content, ok := e.partials[node.Name]
	if !ok {
		return "", fmt.Errorf("partial not found: %s", node.Name)
	}

	// Parse and render the partial
	lexer := NewLexer(content)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return "", fmt.Errorf("partial %s: %w", node.Name, err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		return "", fmt.Errorf("partial %s: %w", node.Name, err)
	}

	return e.Render(ast)
}
