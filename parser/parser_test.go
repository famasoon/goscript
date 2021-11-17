package parser

import (
	"fmt"
	"goscript/ast"
	"goscript/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil!\n")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d\n", len(program.Statements))
	}

	tests := []struct {
		testName           string
		expectedIdentifier string
	}{
		{"x case", "x"},
		{"y case", "y"},
		{"foobar case", "foobar"},
	}

	for i, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			stmt := program.Statements[i]
			err := testLetStatement(t, stmt, tt.expectedIdentifier)
			if err != nil {
				t.Errorf("[Error] %v", err)
			}
		})
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) error {
	if s.TokenLiteral() != "let" {
		return fmt.Errorf("s.TokenLitral not 'let'. got=%q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		return fmt.Errorf("s not *ast.LetStatement. got=%T", s)
	}

	if letStmt.Name.Value != name {
		return fmt.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != name {
		return fmt.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
	}

	return nil
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors,", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
