// Package tmplx implements a template engine with Mustache/Jinja2-like syntax.
// It features a full lexer, recursive descent parser, AST, and tree-walking evaluator.
package tmplx

import "fmt"

// NodeType represents the type of an AST node.
type NodeType int

const (
	NodeText NodeType = iota
	NodeVariable
	NodeRaw
	NodeIf
	NodeEach
	NodePartial
	NodeComment
	NodeTemplate
	NodeFilter
	NodeUnless
	NodeWith
)

// Position tracks source location for error reporting.
type Position struct {
	Line   int
	Column int
	Offset int
}

func (p Position) String() string {
	return fmt.Sprintf("line %d, column %d", p.Line, p.Column)
}

// Node is the interface implemented by all AST nodes.
type Node interface {
	Type() NodeType
	Position() Position
	String() string
}

// TextNode represents literal text content.
type TextNode struct {
	Text  string
	Pos   Position
	Value string // resolved value after evaluation
}

func (n *TextNode) Type() NodeType     { return NodeText }
func (n *TextNode) Position() Position { return n.Pos }
func (n *TextNode) String() string     { return fmt.Sprintf("Text(%q)", n.Text) }

// VariableNode represents a variable reference with optional filters.
// {{ variable | filter1 | filter2(arg) }}
type VariableNode struct {
	Name    string
	Filters []*FilterExpr
	Pos     Position
}

func (n *VariableNode) Type() NodeType     { return NodeVariable }
func (n *VariableNode) Position() Position { return n.Pos }
func (n *VariableNode) String() string     { return fmt.Sprintf("Variable(%s)", n.Name) }

// FilterExpr represents a filter applied to a variable.
type FilterExpr struct {
	Name string
	Args []Node // can be VariableNode or TextNode (string literal)
	Pos  Position
}

// RawNode represents raw/unescaped output: {{{ content }}} or {% raw %}...{% endraw %}
type RawNode struct {
	Content Node
	Pos     Position
}

func (n *RawNode) Type() NodeType     { return NodeRaw }
func (n *RawNode) Position() Position { return n.Pos }
func (n *RawNode) String() string     { return fmt.Sprintf("Raw(%s)", n.Content) }

// IfNode represents a conditional: {{# if condition }}...{{ else }}...{{/ if }}
type IfNode struct {
	Condition Node
	Body      []Node
	ElseBody  []Node
	Pos       Position
}

func (n *IfNode) Type() NodeType     { return NodeIf }
func (n *IfNode) Position() Position { return n.Pos }
func (n *IfNode) String() string     { return fmt.Sprintf("If(%s)", n.Condition) }

// UnlessNode represents a negative conditional: {{# unless condition }}...{{/ unless }}
type UnlessNode struct {
	Condition Node
	Body      []Node
	Pos       Position
}

func (n *UnlessNode) Type() NodeType     { return NodeUnless }
func (n *UnlessNode) Position() Position { return n.Pos }
func (n *UnlessNode) String() string     { return fmt.Sprintf("Unless(%s)", n.Condition) }

// EachNode represents a loop: {{# each items }}...{{/ each }}
type EachNode struct {
	Iterable Node
	KeyVar   string // optional key variable name
	Body     []Node
	Pos      Position
}

func (n *EachNode) Type() NodeType     { return NodeEach }
func (n *EachNode) Position() Position { return n.Pos }
func (n *EachNode) String() string     { return fmt.Sprintf("Each(%s)", n.Iterable) }

// WithNode represents a scope change: {{# with user }}...{{/ with }}
type WithNode struct {
	Variable Node
	Body     []Node
	Pos      Position
}

func (n *WithNode) Type() NodeType     { return NodeWith }
func (n *WithNode) Position() Position { return n.Pos }
func (n *WithNode) String() string     { return fmt.Sprintf("With(%s)", n.Variable) }

// PartialNode represents a partial include: {{> partial_name }}
type PartialNode struct {
	Name string
	Pos  Position
}

func (n *PartialNode) Type() NodeType     { return NodePartial }
func (n *PartialNode) Position() Position { return n.Pos }
func (n *PartialNode) String() string     { return fmt.Sprintf("Partial(%s)", n.Name) }

// CommentNode represents a template comment: {{! comment }} or {{/* comment */}}
type CommentNode struct {
	Text string
	Pos  Position
}

func (n *CommentNode) Type() NodeType     { return NodeComment }
func (n *CommentNode) Position() Position { return n.Pos }
func (n *CommentNode) String() string     { return fmt.Sprintf("Comment(%q)", n.Text) }

// TemplateNode represents the root template containing child nodes.
type TemplateNode struct {
	Children []Node
	Pos      Position
}

func (n *TemplateNode) Type() NodeType     { return NodeTemplate }
func (n *TemplateNode) Position() Position { return n.Pos }
func (n *TemplateNode) String() string     { return "Template" }

// StringLiteralNode represents a string literal in filter arguments.
type StringLiteralNode struct {
	Value string
	Pos   Position
}

func (n *StringLiteralNode) Type() NodeType     { return NodeText }
func (n *StringLiteralNode) Position() Position { return n.Pos }
func (n *StringLiteralNode) String() string     { return fmt.Sprintf("String(%q)", n.Value) }

// IntLiteralNode represents an integer literal in filter arguments.
type IntLiteralNode struct {
	Value int
	Pos   Position
}

func (n *IntLiteralNode) Type() NodeType     { return NodeText }
func (n *IntLiteralNode) Position() Position { return n.Pos }
func (n *IntLiteralNode) String() string     { return fmt.Sprintf("Int(%d)", n.Value) }
