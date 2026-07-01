package tmplx

import (
	"testing"
)

func TestParserVariable(t *testing.T) {
	lexer := NewLexer("Hello {{ name }}!")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(ast.Children) != 3 { // Text, Variable, Text
		t.Fatalf("Expected 3 children, got %d", len(ast.Children))
	}

	if _, ok := ast.Children[1].(*VariableNode); !ok {
		t.Errorf("Expected VariableNode, got %T", ast.Children[1])
	}
}

func TestParserIfBlock(t *testing.T) {
	lexer := NewLexer("{{# if active }}Active{{/ if }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(ast.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(ast.Children))
	}

	ifNode, ok := ast.Children[0].(*IfNode)
	if !ok {
		t.Fatalf("Expected IfNode, got %T", ast.Children[0])
	}

	if len(ifNode.Body) != 1 {
		t.Errorf("Expected 1 body node, got %d", len(ifNode.Body))
	}
}

func TestParserIfElseBlock(t *testing.T) {
	lexer := NewLexer("{{# if active }}Yes{{ else }}No{{/ if }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ifNode, ok := ast.Children[0].(*IfNode)
	if !ok {
		t.Fatalf("Expected IfNode, got %T", ast.Children[0])
	}

	if len(ifNode.ElseBody) != 1 {
		t.Errorf("Expected 1 else body node, got %d", len(ifNode.ElseBody))
	}
}

func TestParserEachBlock(t *testing.T) {
	lexer := NewLexer("{{# each items }}{{ item }}{{/ each }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	eachNode, ok := ast.Children[0].(*EachNode)
	if !ok {
		t.Fatalf("Expected EachNode, got %T", ast.Children[0])
	}

	if len(eachNode.Body) != 1 {
		t.Errorf("Expected 1 body node, got %d", len(eachNode.Body))
	}
}

func TestParserPartial(t *testing.T) {
	lexer := NewLexer("{{> header }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(ast.Children) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(ast.Children))
	}

	partialNode, ok := ast.Children[0].(*PartialNode)
	if !ok {
		t.Fatalf("Expected PartialNode, got %T", ast.Children[0])
	}

	if partialNode.Name != "header" {
		t.Errorf("Expected 'header', got %q", partialNode.Name)
	}
}

func TestParserVariableWithFilter(t *testing.T) {
	lexer := NewLexer("{{ name | upper }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	varNode, ok := ast.Children[0].(*VariableNode)
	if !ok {
		t.Fatalf("Expected VariableNode, got %T", ast.Children[0])
	}

	if len(varNode.Filters) != 1 {
		t.Fatalf("Expected 1 filter, got %d", len(varNode.Filters))
	}

	if varNode.Filters[0].Name != "upper" {
		t.Errorf("Expected 'upper', got %q", varNode.Filters[0].Name)
	}
}

func TestParserMultipleFilters(t *testing.T) {
	lexer := NewLexer("{{ name | upper | truncate(10) }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	varNode, ok := ast.Children[0].(*VariableNode)
	if !ok {
		t.Fatalf("Expected VariableNode, got %T", ast.Children[0])
	}

	if len(varNode.Filters) != 2 {
		t.Fatalf("Expected 2 filters, got %d", len(varNode.Filters))
	}
}

func TestParserNestedBlocks(t *testing.T) {
	lexer := NewLexer("{{# if active }}{{# each items }}{{ name }}{{/ each }}{{/ if }}")
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}

	parser := NewParser(tokens)
	ast, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	ifNode, ok := ast.Children[0].(*IfNode)
	if !ok {
		t.Fatalf("Expected IfNode, got %T", ast.Children[0])
	}

	if len(ifNode.Body) != 1 {
		t.Fatalf("Expected 1 body node, got %d", len(ifNode.Body))
	}

	_, ok = ifNode.Body[0].(*EachNode)
	if !ok {
		t.Fatalf("Expected EachNode in if body, got %T", ifNode.Body[0])
	}
}
