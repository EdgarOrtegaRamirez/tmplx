package tmplx

import (
	"testing"
)

func TestLexerBasic(t *testing.T) {
	lexer := NewLexer("Hello {{ name }}!")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	// Expected: Text, VarStart, Ident, VarEnd, Text, EOF
	if len(tokens) != 6 {
		t.Fatalf("Expected 6 tokens, got %d: %v", len(tokens), tokens)
	}

	if tokens[0].Type != TokenText {
		t.Errorf("Expected Text, got %s", tokens[0].Type)
	}
	if tokens[0].Value != "Hello " {
		t.Errorf("Expected 'Hello ', got %q", tokens[0].Value)
	}
}

func TestLexerFilters(t *testing.T) {
	lexer := NewLexer("{{ name | upper | truncate(10) }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	// Should have: VarStart, Ident(name), Filter(upper), Filter(truncate), VarEnd, EOF
	found := false
	for _, tok := range tokens {
		if tok.Type == TokenFilterStart && tok.Value == "upper" {
			found = true
		}
	}
	if !found {
		t.Error("Expected 'upper' filter token")
	}
}

func TestLexerComment(t *testing.T) {
	lexer := NewLexer("Hello {{! this is a comment }} World")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	// Should have: Text("Hello "), CommentEnd, Text(" World")
	if len(tokens) < 2 {
		t.Fatalf("Expected at least 2 tokens, got %d", len(tokens))
	}
}

func TestLexerBlock(t *testing.T) {
	lexer := NewLexer("{{# if active }}Active{{/ if }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	foundStart := false
	foundEnd := false
	for _, tok := range tokens {
		if tok.Type == TokenBlockStart && tok.Value == "if" {
			foundStart = true
		}
		if tok.Type == TokenBlockEnd && tok.Value == "if" {
			foundEnd = true
		}
	}
	if !foundStart {
		t.Error("Expected block start token for 'if'")
	}
	if !foundEnd {
		t.Error("Expected block end token for 'if'")
	}
}

func TestLexerRaw(t *testing.T) {
	lexer := NewLexer("Before {{{ raw content }}} After")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	foundStart := false
	for _, tok := range tokens {
		if tok.Type == TokenRawStart {
			foundStart = true
		}
	}
	if !foundStart {
		t.Error("Expected raw start token")
	}
}

func TestLexerPartial(t *testing.T) {
	lexer := NewLexer("{{> header }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	found := false
	for _, tok := range tokens {
		if tok.Type == TokenPartialStart && tok.Value == "header" {
			found = true
		}
	}
	if !found {
		t.Error("Expected partial token for 'header'")
	}
}

func TestLexerEach(t *testing.T) {
	lexer := NewLexer("{{# each items }}{{ item }}{{/ each }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	foundEach := false
	for _, tok := range tokens {
		if tok.Type == TokenBlockStart && tok.Value == "each" {
			foundEach = true
		}
	}
	if !foundEach {
		t.Error("Expected each block start token")
	}
}

func TestLexerUnless(t *testing.T) {
	lexer := NewLexer("{{# unless hidden }}Visible{{/ unless }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	found := false
	for _, tok := range tokens {
		if tok.Type == TokenBlockStart && tok.Value == "unless" {
			found = true
		}
	}
	if !found {
		t.Error("Expected unless block start token")
	}
}

func TestLexerWith(t *testing.T) {
	lexer := NewLexer("{{# with user }}{{ name }}{{/ with }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	found := false
	for _, tok := range tokens {
		if tok.Type == TokenBlockStart && tok.Value == "with" {
			found = true
		}
	}
	if !found {
		t.Error("Expected with block start token")
	}
}

func TestLexerNested(t *testing.T) {
	lexer := NewLexer("{{# if active }}{{# each items }}{{ name }}{{/ each }}{{/ if }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	ifCount := 0
	eachCount := 0
	for _, tok := range tokens {
		if tok.Type == TokenBlockStart {
			if tok.Value == "if" {
				ifCount++
			}
			if tok.Value == "each" {
				eachCount++
			}
		}
	}
	if ifCount != 1 {
		t.Errorf("Expected 1 if block, got %d", ifCount)
	}
	if eachCount != 1 {
		t.Errorf("Expected 1 each block, got %d", eachCount)
	}
}
