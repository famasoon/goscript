package lexer

import (
	"goscript/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = fn(x, y) {
		x + y;
	};
	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;
	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	10 == 10;
	10 != 9;
	"foobar"
	"foo bar"
	[1, 2];
	`

	tests := []struct {
		testName        string
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{"Let test case", token.LET, "let"},
		{"IDENT test case", token.IDENT, "five"},
		{"ASSIGN test case", token.ASSIGN, "="},
		{"INT test case", token.INT, "5"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"Let test case", token.LET, "let"},
		{"IDENT test case", token.IDENT, "ten"},
		{"ASSIGN test case", token.ASSIGN, "="},
		{"INT test case", token.INT, "10"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"LET test case", token.LET, "let"},
		{"IDENT test case", token.IDENT, "add"},
		{"ASSIGN test case", token.ASSIGN, "="},
		{"FUNCTION test case", token.FUNCTION, "fn"},
		{"LPAREN test case", token.LPAREN, "("},
		{"IDENT test case", token.IDENT, "x"},
		{"COMMA test case", token.COMMA, ","},
		{"IDENT test case", token.IDENT, "y"},
		{"RPAREN test case", token.RPAREN, ")"},
		{"LBRACE test case", token.LBRACE, "{"},
		{"IDENT test case", token.IDENT, "x"},
		{"PLUS test case", token.PLUS, "+"},
		{"IDENT test case", token.IDENT, "y"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"RBRACE test case", token.RBRACE, "}"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"Let test case", token.LET, "let"},
		{"IDENT test case", token.IDENT, "result"},
		{"ASSIGN test case", token.ASSIGN, "="},
		{"IDENT test case", token.IDENT, "add"},
		{"LPAREN test case", token.LPAREN, "("},
		{"IDENT test case", token.IDENT, "five"},
		{"COMMA test case", token.COMMA, ","},
		{"IDENT test case", token.IDENT, "ten"},
		{"RPAREN test case", token.RPAREN, ")"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"BANG test case", token.BANG, "!"},
		{"MINUS test case", token.MINUS, "-"},
		{"SLASH test case", token.SLASH, "/"},
		{"ASTERISK test case", token.ASTERISK, "*"},
		{"INT test case", token.INT, "5"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"INT test case", token.INT, "5"},
		{"LT test case", token.LT, "<"},
		{"INT test case", token.INT, "10"},
		{"GT test case", token.GT, ">"},
		{"INT test case", token.INT, "5"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"IF test case", token.IF, "if"},
		{"LPAREN test case", token.LPAREN, "("},
		{"INT test case", token.INT, "5"},
		{"LT test case", token.LT, "<"},
		{"INT test case", token.INT, "10"},
		{"RPAREN test case", token.RPAREN, ")"},
		{"LBRACE test case", token.LBRACE, "{"},
		{"RETURN test case", token.RETURN, "return"},
		{"TRUE test case", token.TRUE, "true"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"RBRACE test case", token.RBRACE, "}"},
		{"ELSE test case", token.ELSE, "else"},
		{"LBRACE test case", token.LBRACE, "{"},
		{"RETURN test case", token.RETURN, "return"},
		{"FALSE test case", token.FALSE, "false"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"RBRACE test case", token.RBRACE, "}"},
		{"INT test case", token.INT, "10"},
		{"EQ test case", token.EQ, "=="},
		{"INT test case", token.INT, "10"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"INT test case", token.INT, "10"},
		{"NOT_EQ test case", token.NOT_EQ, "!="},
		{"INT test case", token.INT, "9"},
		{"SEMICOLON test case", token.SEMICOLON, ";"},
		{"foobar test case", token.STRING, "foobar"},
		{"foo bar test case", token.STRING, "foo bar"},
		{"[", token.LBRACKET, "["},
		{"1", token.INT, "1"},
		{",", token.COMMA, ","},
		{"2", token.INT, "2"},
		{"]", token.RBRACKET, "]"},
		{";", token.SEMICOLON, ";"},
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
