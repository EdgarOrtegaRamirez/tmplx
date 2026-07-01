package tmplx

import (
	"fmt"
	"strconv"
)

// Parser converts tokens into an AST.
type Parser struct {
	tokens *TokenList
}

// NewParser creates a new parser from tokens.
func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: NewTokenList(tokens)}
}

// Parse parses the tokens into a template AST.
func (p *Parser) Parse() (*TemplateNode, error) {
	root := &TemplateNode{
		Pos: p.tokens.Peek().Pos,
	}

	for p.tokens.HasMore() {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			root.Children = append(root.Children, node)
		}
	}

	return root, nil
}

func (p *Parser) parseNode() (Node, error) {
	tok := p.tokens.Peek()

	switch tok.Type {
	case TokenText:
		p.tokens.Advance()
		return &TextNode{
			Text: tok.Value,
			Pos:  tok.Pos,
		}, nil

	case TokenVariableStart:
		return p.parseVariable()

	case TokenRawStart:
		return p.parseRaw()

	case TokenBlockStart:
		return p.parseBlock()

	case TokenElse:
		// else is handled by the parent block parser
		return nil, nil

	case TokenBlockEnd:
		// block end is handled by the parent block parser
		// but we need to advance to prevent infinite loops
		p.tokens.Advance()
		return nil, nil

	case TokenCommentEnd:
		p.tokens.Advance()
		return nil, nil

	case TokenVariableEnd:
		// Closing braces - skip
		p.tokens.Advance()
		return nil, nil

	case TokenPartialStart:
		p.tokens.Advance()
		return &PartialNode{
			Name: tok.Value,
			Pos:  tok.Pos,
		}, nil

	case TokenEOF:
		return nil, nil

	default:
		p.tokens.Advance()
		return nil, fmt.Errorf("unexpected token %s at %s", tok.Type, tok.Pos)
	}
}

func (p *Parser) parseVariable() (Node, error) {
	tok := p.tokens.Expect(TokenVariableStart)

	// Read the variable name (first identifier)
	nameTok := p.tokens.Expect(TokenIdent)
	name := nameTok.Value

	// Read dotted path
	for p.tokens.Peek().Type == TokenDot {
		p.tokens.Advance() // dot
		partTok := p.tokens.Expect(TokenIdent)
		name += "." + partTok.Value
	}

	// Read filters
	var filters []*FilterExpr
	for p.tokens.Peek().Type == TokenFilterStart {
		filterTok := p.tokens.Advance()
		filter := &FilterExpr{
			Name: filterTok.Value,
			Pos:  filterTok.Pos,
		}

		// Read filter arguments if present
		// Arguments are already tokenized as STRING, INT, or IDENT tokens
		// We need to collect them until we hit a filter or closing braces
		for {
			next := p.tokens.Peek()
			if next.Type == TokenFilterStart || next.Type == TokenVariableEnd || next.Type == TokenEOF {
				break
			}
			// Consume argument token
			argTok := p.tokens.Advance()
			switch argTok.Type {
			case TokenStringLit:
				filter.Args = append(filter.Args, &StringLiteralNode{
					Value: argTok.Value,
					Pos:  argTok.Pos,
				})
			case TokenIntLit:
				val, _ := strconv.Atoi(argTok.Value)
				filter.Args = append(filter.Args, &IntLiteralNode{
					Value: val,
					Pos:  argTok.Pos,
				})
			case TokenIdent:
				filter.Args = append(filter.Args, &VariableNode{
					Name: argTok.Value,
					Pos:  argTok.Pos,
				})
			}
		}

		filters = append(filters, filter)
	}

	// Skip closing }}
	p.expectClosingBraces(tok.Pos)

	return &VariableNode{
		Name:    name,
		Filters: filters,
		Pos:     tok.Pos,
	}, nil
}

func (p *Parser) parseRaw() (Node, error) {
	tok := p.tokens.Expect(TokenRawStart)

	// Read text content until raw end
	var content string
	if p.tokens.Peek().Type == TokenText {
		content = p.tokens.Advance().Value
	}

	p.tokens.Expect(TokenRawEnd)

	return &RawNode{
		Content: &TextNode{Text: content, Pos: tok.Pos},
		Pos:     tok.Pos,
	}, nil
}

func (p *Parser) parseBlock() (Node, error) {
	tok := p.tokens.Expect(TokenBlockStart)
	blockType := tok.Value

	switch blockType {
	case "if":
		return p.parseIfBlock(tok)
	case "unless":
		return p.parseUnlessBlock(tok)
	case "each":
		return p.parseEachBlock(tok)
	case "with":
		return p.parseWithBlock(tok)
	default:
		return nil, fmt.Errorf("unknown block type %q at %s", blockType, tok.Pos)
	}
}

func (p *Parser) parseIfBlock(tok Token) (Node, error) {
	// Read condition (variable name)
	condTok := p.tokens.Expect(TokenIdent)
	condition := &VariableNode{
		Name: condTok.Value,
		Pos:  condTok.Pos,
	}

	// Skip closing }}
	p.expectClosingBraces(tok.Pos)

	// Parse body until {{ else }} or {{/ if }}
	var body []Node
	var elseBody []Node

	for p.tokens.HasMore() {
		next := p.tokens.Peek()

		if next.Type == TokenElse {
			p.tokens.Advance()
			p.expectClosingBraces(next.Pos)

			// Parse else body
			for p.tokens.HasMore() {
				next2 := p.tokens.Peek()
				if next2.Type == TokenBlockEnd && next2.Value == "if" {
					break
				}
				node, err := p.parseNode()
				if err != nil {
					return nil, err
				}
				if node != nil {
					elseBody = append(elseBody, node)
				}
			}
			break
		}

		if next.Type == TokenBlockEnd && next.Value == "if" {
			p.tokens.Advance()
			p.expectClosingBraces(next.Pos)
			break
		}

		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			body = append(body, node)
		}
	}

	return &IfNode{
		Condition: condition,
		Body:      body,
		ElseBody:  elseBody,
		Pos:       tok.Pos,
	}, nil
}

func (p *Parser) parseUnlessBlock(tok Token) (Node, error) {
	condTok := p.tokens.Expect(TokenIdent)
	condition := &VariableNode{
		Name: condTok.Value,
		Pos:  condTok.Pos,
	}

	p.expectClosingBraces(tok.Pos)

	var body []Node
	for p.tokens.HasMore() {
		next := p.tokens.Peek()
		if next.Type == TokenBlockEnd && next.Value == "unless" {
			p.tokens.Advance()
			p.expectClosingBraces(next.Pos)
			break
		}
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			body = append(body, node)
		}
	}

	return &UnlessNode{
		Condition: condition,
		Body:      body,
		Pos:       tok.Pos,
	}, nil
}

func (p *Parser) parseEachBlock(tok Token) (Node, error) {
	iterTok := p.tokens.Expect(TokenIdent)
	iterable := &VariableNode{
		Name: iterTok.Value,
		Pos:  iterTok.Pos,
	}

	// Optional "as" key variable: {{# each items as item }}
	keyVar := ""
	if p.tokens.Peek().Type == TokenIdent && p.tokens.Peek().Value == "as" {
		p.tokens.Advance() // "as"
		keyTok := p.tokens.Expect(TokenIdent)
		keyVar = keyTok.Value
	}

	p.expectClosingBraces(tok.Pos)

	var body []Node
	for p.tokens.HasMore() {
		next := p.tokens.Peek()
		if next.Type == TokenBlockEnd && next.Value == "each" {
			p.tokens.Advance()
			p.expectClosingBraces(next.Pos)
			break
		}
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			body = append(body, node)
		}
	}

	return &EachNode{
		Iterable: iterable,
		KeyVar:   keyVar,
		Body:     body,
		Pos:      tok.Pos,
	}, nil
}

func (p *Parser) parseWithBlock(tok Token) (Node, error) {
	varTok := p.tokens.Expect(TokenIdent)
	variable := &VariableNode{
		Name: varTok.Value,
		Pos:  varTok.Pos,
	}

	p.expectClosingBraces(tok.Pos)

	var body []Node
	for p.tokens.HasMore() {
		next := p.tokens.Peek()
		if next.Type == TokenBlockEnd && next.Value == "with" {
			p.tokens.Advance()
			p.expectClosingBraces(next.Pos)
			break
		}
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			body = append(body, node)
		}
	}

	return &WithNode{
		Variable: variable,
		Body:     body,
		Pos:      tok.Pos,
	}, nil
}

func (p *Parser) expectClosingBraces(startPos Position) {
	// The closing }} is now a VAR_END token produced by the lexer
	if p.tokens.Peek().Type == TokenVariableEnd {
		p.tokens.Advance()
	}
	// Also handle any empty Text tokens
	for p.tokens.Peek().Type == TokenText && p.tokens.Peek().Value == "" {
		p.tokens.Advance()
	}
}
