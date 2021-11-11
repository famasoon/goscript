package lexer

import (
	"goscript/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		testName        string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{"= test case", token.ASSIGN, "="},
		{"+ test case", token.PLUS, "+"},
		{"( test case", token.LPAREN, "("},
		{") test case", token.RPAREN, ")"},
		{"{ test case", token.LBRACE, "{"},
		{"} test case", token.RBRACE, "}"},
		{", test case", token.COMMA, ","},
		{"; test case", token.SEMMICOLON, ";"},
		{"EOF test case", token.EOF, ""},
	}

	l := New(input)

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			token := l.NextToken()

			if token.Type != tt.expectedType {
				t.Errorf("%s: tokenType wrong. expected=%q, got=%q\n", tt.testName, tt.expectedType, token.Type)
			}

			if token.Literal != tt.expectedLiteral {
				t.Errorf("%s: tokenLiteral wrong. expected=%q, got=%q\n", tt.testName, tt.expectedLiteral, token.Literal)
			}
		})
	}
}
