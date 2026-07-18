package tmplx

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of a token.
type TokenType int

const (
	TokenText          TokenType = iota
	TokenVariableStart           // {{
	TokenVariableEnd             // }}
	TokenRawStart                // {{{
	TokenRawEnd                  // }}}
	TokenBlockStart              // {{#
	TokenBlockEnd                // {{/}}
	TokenFilterStart             // |
	TokenPartialStart            // {{>
	TokenCommentStart            // {{! or {{/*
	TokenCommentEnd              // */}}
	TokenElse                    // {{ else }}
	TokenStringLit               // "string"
	TokenIntLit                  // 123
	TokenIdent                   // identifier
	TokenDot                     // .
	TokenComma                   // ,
	TokenEOF
)

var tokenNames = map[TokenType]string{
	TokenText:          "TEXT",
	TokenVariableStart: "VAR_START",
	TokenVariableEnd:   "VAR_END",
	TokenRawStart:      "RAW_START",
	TokenRawEnd:        "RAW_END",
	TokenBlockStart:    "BLOCK_START",
	TokenBlockEnd:      "BLOCK_END",
	TokenFilterStart:   "FILTER",
	TokenPartialStart:  "PARTIAL",
	TokenCommentStart:  "COMMENT_START",
	TokenCommentEnd:    "COMMENT_END",
	TokenElse:          "ELSE",
	TokenStringLit:     "STRING",
	TokenIntLit:        "INT",
	TokenIdent:         "IDENT",
	TokenDot:           "DOT",
	TokenComma:         "COMMA",
	TokenEOF:           "EOF",
}

func (t TokenType) String() string {
	if name, ok := tokenNames[t]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", int(t))
}

// Token represents a lexical token.
type Token struct {
	Type  TokenType
	Value string
	Pos   Position
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q)@%s", t.Type, t.Value, t.Pos)
}

// Lexer tokenizes template source code.
type Lexer struct {
	input  string
	pos    int
	line   int
	col    int
	tokens []Token
}

// NewLexer creates a new lexer for the given input.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   0,
		line:  1,
		col:   1,
	}
}

// Tokenize returns all tokens from the input.
func (l *Lexer) Tokenize() ([]Token, error) {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]

		// Check for template tags: {{ }}
		if ch == '{' && l.peek(1) == '{' {
			if err := l.readTemplateTag(); err != nil {
				return nil, err
			}
			continue
		}

		// Regular text
		l.readText()
	}

	l.tokens = append(l.tokens, Token{
		Type:  TokenEOF,
		Value: "",
		Pos:   l.currentPos(),
	})

	return l.tokens, nil
}

func (l *Lexer) currentPos() Position {
	return Position{Line: l.line, Column: l.col, Offset: l.pos}
}

func (l *Lexer) peek(offset int) byte {
	idx := l.pos + offset
	if idx >= len(l.input) {
		return 0
	}
	return l.input[idx]
}

func (l *Lexer) advance() byte {
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}
	return ch
}

func (l *Lexer) readText() {
	start := l.pos
	startPos := l.currentPos()
	for l.pos < len(l.input) {
		if l.input[l.pos] == '{' && l.peek(1) == '{' {
			break
		}
		l.advance()
	}
	if l.pos > start {
		l.tokens = append(l.tokens, Token{
			Type:  TokenText,
			Value: l.input[start:l.pos],
			Pos:   startPos,
		})
	}
}

func (l *Lexer) readTemplateTag() error {
	startPos := l.currentPos()

	// Skip opening {{
	l.advance() // {
	l.advance() // {

	l.skipWhitespace()

	// Check for raw: {{{ content }}} (three braces)
	if l.peek(0) == '{' {
		// We've consumed {{, now check if there's a third {
		// If so, this is a raw block: {{{ content }}}
		l.advance() // third {
		l.skipWhitespace()
		return l.readRawTag()
	}

	// Check for comment: {{! comment }} or {{/* comment */}}
	if l.peek(0) == '!' {
		return l.readComment()
	}
	if l.peek(0) == '/' && l.peek(1) == '*' {
		return l.readBlockComment()
	}

	// Check for block start: {{# if/each/with/unless }}
	if l.peek(0) == '#' {
		return l.readBlockStartTag()
	}

	// Check for block end: {{/ if/each/with/unless }}
	if l.peek(0) == '/' {
		return l.readBlockEndTag()
	}

	// Check for else: {{ else }}
	if l.peek(0) == 'e' && l.peek(1) == 'l' && l.peek(2) == 's' && l.peek(3) == 'e' {
		l.tokens = append(l.tokens, Token{
			Type:  TokenElse,
			Value: "else",
			Pos:   startPos,
		})
		// Skip "else" and closing }}
		l.advance() // e
		l.advance() // l
		l.advance() // s
		l.advance() // e
		l.skipWhitespace()
		return l.readClosingBraces(startPos)
	}

	// Check for partial: {{> name }}
	if l.peek(0) == '>' {
		return l.readPartialTag()
	}

	// Otherwise it's a variable expression: {{ variable | filter }}
	return l.readVariableExpr()
}

func (l *Lexer) readRawTag() error {
	startPos := l.currentPos()
	l.tokens = append(l.tokens, Token{
		Type:  TokenRawStart,
		Value: "{{{",
		Pos:   startPos,
	})

	// Read until }}}
	for l.pos < len(l.input) {
		if l.input[l.pos] == '}' && l.peek(1) == '}' && l.peek(2) == '}' {
			l.advance() // }
			l.advance() // }
			l.advance() // }
			l.tokens = append(l.tokens, Token{
				Type:  TokenRawEnd,
				Value: "}}}",
				Pos:   l.currentPos(),
			})
			return nil
		}
		l.advance()
	}

	return fmt.Errorf("unclosed raw block at %s", startPos)
}

func (l *Lexer) readBlockComment() error {
	startPos := l.currentPos()
	l.advance() // /
	l.advance() // *

	for l.pos < len(l.input) {
		if l.input[l.pos] == '*' && l.peek(1) == '/' {
			l.advance() // *
			l.advance() // /
			// Skip closing }}
			l.skipWhitespace()
			return l.readClosingBraces(startPos)
		}
		l.advance()
	}

	return fmt.Errorf("unclosed block comment at %s", startPos)
}

func (l *Lexer) readComment() error {
	startPos := l.currentPos()
	l.advance() // !

	for l.pos < len(l.input) {
		if l.input[l.pos] == '}' && l.peek(1) == '}' {
			l.advance() // }
			l.advance() // }
			l.tokens = append(l.tokens, Token{
				Type:  TokenCommentEnd,
				Value: "}}",
				Pos:   l.currentPos(),
			})
			return nil
		}
		l.advance()
	}

	return fmt.Errorf("unclosed comment at %s", startPos)
}

func (l *Lexer) readBlockStartTag() error {
	startPos := l.currentPos()
	l.advance() // #

	l.skipWhitespace()
	ident := l.readIdentifier()
	l.skipWhitespace()

	l.tokens = append(l.tokens, Token{
		Type:  TokenBlockStart,
		Value: ident,
		Pos:   startPos,
	})

	// Read arguments (identifiers) until closing braces
	for l.peek(0) != '}' && l.peek(0) != 0 {
		if l.peek(0) == ' ' || l.peek(0) == '	' || l.peek(0) == '\n' {
			l.skipWhitespace()
			continue
		}
		// Read identifier or dotted path
		argStart := l.currentPos()
		l.readDottedIdent()
		// The identifier is already appended by readDottedIdent
		_ = argStart
	}

	return l.readClosingBraces(startPos)
}

func (l *Lexer) readBlockEndTag() error {
	startPos := l.currentPos()
	l.advance() // /

	l.skipWhitespace()
	ident := l.readIdentifier()
	l.skipWhitespace()

	l.tokens = append(l.tokens, Token{
		Type:  TokenBlockEnd,
		Value: ident,
		Pos:   startPos,
	})

	return l.readClosingBraces(startPos)
}

func (l *Lexer) readPartialTag() error {
	startPos := l.currentPos()
	l.advance() // >

	l.skipWhitespace()
	ident := l.readIdentifier()
	l.skipWhitespace()

	l.tokens = append(l.tokens, Token{
		Type:  TokenPartialStart,
		Value: ident,
		Pos:   startPos,
	})

	return l.readClosingBraces(startPos)
}

func (l *Lexer) readVariableExpr() error {
	startPos := l.currentPos()

	// Append VarStart token
	l.tokens = append(l.tokens, Token{
		Type:  TokenVariableStart,
		Value: "{{",
		Pos:   startPos,
	})

	// Read identifier or dotted path
	l.readDottedIdent()

	l.skipWhitespace()

	// Read filters: | filter | filter(arg)
	for l.peek(0) == '|' {
		l.advance() // |
		l.skipWhitespace()

		filterPos := l.currentPos()
		filterName := l.readIdentifier()
		l.tokens = append(l.tokens, Token{
			Type:  TokenFilterStart,
			Value: filterName,
			Pos:   filterPos,
		})

		l.skipWhitespace()

		// Check for filter arguments: (arg1, arg2)
		if l.peek(0) == '(' {
			l.advance() // (
			l.skipWhitespace()

			for l.peek(0) != ')' && l.peek(0) != 0 {
				l.skipWhitespace()
				l.readFilterArg()
				l.skipWhitespace()
				if l.peek(0) == ',' {
					l.advance() // ,
					l.skipWhitespace()
				}
			}
			if l.peek(0) == ')' {
				l.advance() // )
			}
		}

		l.skipWhitespace()
	}

	return l.readClosingBraces(startPos)
}

func (l *Lexer) readFilterArg() {
	if l.peek(0) == '"' || l.peek(0) == '\'' {
		l.readStringLit()
	} else if unicode.IsDigit(rune(l.peek(0))) {
		l.readIntLit()
	} else if l.peek(0) == '-' || unicode.IsDigit(rune(l.peek(0))) {
		l.readIntLit()
	} else {
		l.readDottedIdent()
	}
}

func (l *Lexer) readStringLit() {
	quote := l.advance() // opening quote
	start := l.pos
	startPos := l.currentPos()
	for l.pos < len(l.input) && l.input[l.pos] != quote {
		if l.input[l.pos] == '\\' {
			l.advance() // skip escape
		}
		l.advance()
	}
	value := l.input[start:l.pos]
	if l.pos < len(l.input) {
		l.advance() // closing quote
	}
	l.tokens = append(l.tokens, Token{
		Type:  TokenStringLit,
		Value: value,
		Pos:   startPos,
	})
}

func (l *Lexer) readIntLit() {
	start := l.pos
	startPos := l.currentPos()
	if l.peek(0) == '-' {
		l.advance()
	}
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.peek(0))) {
		l.advance()
	}
	l.tokens = append(l.tokens, Token{
		Type:  TokenIntLit,
		Value: l.input[start:l.pos],
		Pos:   startPos,
	})
}

func (l *Lexer) readDottedIdent() {
	start := l.currentPos()

	// Read first part of identifier
	first := l.readIdentifier()
	if first == "" {
		return
	}

	l.tokens = append(l.tokens, Token{
		Type:  TokenIdent,
		Value: first,
		Pos:   start,
	})

	// Read dotted parts: .name.subname
	for l.peek(0) == '.' && l.peek(1) != '}' && l.peek(1) != '|' && l.peek(1) != ' ' {
		l.advance() // .
		dotPos := l.currentPos()
		part := l.readIdentifier()
		if part != "" {
			l.tokens = append(l.tokens, Token{
				Type:  TokenDot,
				Value: ".",
				Pos:   dotPos,
			})
			l.tokens = append(l.tokens, Token{
				Type:  TokenIdent,
				Value: part,
				Pos:   dotPos,
			})
		}
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for l.pos < len(l.input) {
		ch := rune(l.input[l.pos])
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '-' {
			l.advance()
		} else {
			break
		}
	}
	return l.input[start:l.pos]
}

func (l *Lexer) readClosingBraces(startPos Position) error {
	// Skip to }}
	for l.pos < len(l.input) && (l.input[l.pos] != '}' || l.peek(1) != '}') {
		l.advance()
	}
	if l.pos+1 >= len(l.input) {
		return fmt.Errorf("unclosed template tag at %s", startPos)
	}
	l.advance() // }
	l.advance() // }

	// Append VarEnd token
	l.tokens = append(l.tokens, Token{
		Type:  TokenVariableEnd,
		Value: "}}",
		Pos:   l.currentPos(),
	})

	return nil
}

func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.advance()
		} else {
			break
		}
	}
}

// TokenList is a helper to work with a slice of tokens.
type TokenList struct {
	tokens []Token
	pos    int
}

// NewTokenList creates a new token list wrapper.
func NewTokenList(tokens []Token) *TokenList {
	return &TokenList{tokens: tokens, pos: 0}
}

// Peek returns the current token without advancing.
func (tl *TokenList) Peek() Token {
	if tl.pos >= len(tl.tokens) {
		return Token{Type: TokenEOF}
	}
	return tl.tokens[tl.pos]
}

// PeekAt returns the token at offset from current position.
func (tl *TokenList) PeekAt(offset int) Token {
	idx := tl.pos + offset
	if idx >= len(tl.tokens) {
		return Token{Type: TokenEOF}
	}
	return tl.tokens[idx]
}

// Advance returns the current token and moves to the next.
func (tl *TokenList) Advance() Token {
	t := tl.Peek()
	if tl.pos < len(tl.tokens) {
		tl.pos++
	}
	return t
}

// Expect advances and returns the token, panicking if type doesn't match.
func (tl *TokenList) Expect(tt TokenType) Token {
	t := tl.Advance()
	if t.Type != tt {
		panic(fmt.Sprintf("expected %s, got %s at %s", tt, t.Type, t.Pos))
	}
	return t
}

// HasMore returns true if there are more tokens.
func (tl *TokenList) HasMore() bool {
	return tl.pos < len(tl.tokens) && tl.tokens[tl.pos].Type != TokenEOF
}

// Position returns the current position in the token list.
func (tl *TokenList) Position() int {
	return tl.pos
}

// Backup moves back one token.
func (tl *TokenList) Backup() {
	if tl.pos > 0 {
		tl.pos--
	}
}

// Tokens returns the raw token slice.
func (tl *TokenList) Tokens() []Token {
	return tl.tokens
}

// String returns a debug representation of remaining tokens.
func (tl *TokenList) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i := tl.pos; i < len(tl.tokens) && i < tl.pos+5; i++ {
		if i > tl.pos {
			sb.WriteString(", ")
		}
		sb.WriteString(tl.tokens[i].String())
	}
	if len(tl.tokens)-tl.pos > 5 {
		sb.WriteString(", ...")
	}
	sb.WriteString("]")
	return sb.String()
}
