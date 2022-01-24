package parser

import (
	"fmt"
	"goscript/ast"
	"goscript/lexer"
	"testing"
)

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

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

func TestLetStatements2(t *testing.T) {
	tests := []struct {
		testName           string
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5; case", "let x = 5;", "x", 5},
		{"let y = true; case", "let y = true;", "y", true},
		{"let foobar = y; case", "let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)

			program := p.ParseProgram()
			checkParserErrors(t, p)
			if program == nil {
				t.Fatalf("ParseProgram() returned nil!\n")
			}

			stmt := program.Statements[0]
			err := testLetStatement(t, stmt, tt.expectedIdentifier)
			if err != nil {
				t.Errorf("[Error] %v", err)
			}

			val := stmt.(*ast.LetStatement).Value
			err = testLiteralExpression(t, val, tt.expectedValue)
			if err != nil {
				t.Errorf("[Error] %v", err)
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 99932;
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

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", input, ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", input, ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Errorf("ident.Value not %s. got=%d", input, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", input, literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		testName     string
		input        string
		operator     string
		integerValue int64
	}{
		{"! case", "!5;", "!", 5},
		{"- case", "-15;", "-", 15},
	}

	for _, tt := range prefixTests {
		t.Run(tt.testName, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("stmt is not ast.PrefixExpression. go=%T", stmt.Expression)
			}
			if exp.Operator != tt.operator {
				t.Errorf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
			}

			err := testIntegerLiteral(t, exp.Right, tt.integerValue)
			if err != nil {
				t.Errorf("[Error] %v", err)
			}
		})
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) error {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		return fmt.Errorf("il not *ast.IntegerLiteral. got=%T", il)
	}

	if integ.Value != value {
		return fmt.Errorf("integ.Value not %d. got=%d", value, integ.Value)
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		return fmt.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
	}

	return nil
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		testName   string
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5 case", "5 + 5", 5, "+", 5},
		{"5 - 5 case", "5 - 5", 5, "-", 5},
		{"5 * 5 case", "5 * 5", 5, "*", 5},
		{"5 / 5 case", "5 / 5", 5, "/", 5},
		{"5 > 5 case", "5 > 5", 5, ">", 5},
		{"5 < 5 case", "5 < 5", 5, "<", 5},
		{"5 == 5 case", "5 == 5", 5, "==", 5},
		{"5 != 5 case", "5 != 5", 5, "!=", 5},
		{"true == true case", "true == true", true, "==", true},
		{"true != false case", "true != false", true, "!=", false},
		{"false == false case", "false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		t.Run(tt.testName, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			err := testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
			if err != nil {
				t.Errorf("[Error] %v", err)
			}
		})
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) error {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("exp not *ast.Identifier. got=%T", exp)
	}

	if ident.Value != value {
		return fmt.Errorf("ident.Value not %s. got=%s", value, ident.Value)
	}

	if ident.TokenLiteral() != value {
		return fmt.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
	}

	return nil
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) error {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	return fmt.Errorf("type of exp not handled. got=%T", exp)
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) error {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		return fmt.Errorf("exp is not ast.Expression. got=%T(%s)", exp, exp)
	}

	err := testLiteralExpression(t, opExp.Left, left)
	if err != nil {
		return err
	}

	if opExp.Operator != operator {
		return fmt.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
	}

	err = testLiteralExpression(t, opExp.Right, right)
	if err != nil {
		return err
	}

	return nil
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected bool
	}{
		{"true case", "true", true},
		{"false case", "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			p := New(lexer.New(tt.input))
			program := p.ParseProgram()
			checkParserErrors(t, p)

			l := len(program.Statements)
			if l != 1 {
				t.Fatalf("program has not enough statements. got=%d", l)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T",
					program.Statements[0])
			}

			err := testBooleanLiteral(t, stmt.Expression, tt.expected)
			if err != nil {
				t.Errorf("%T", err)
			}
		})
	}
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) error {
	b, ok := expr.(*ast.Boolean)
	if !ok {
		return fmt.Errorf("b not *ast.Boolean. got=%T", expr)
	}
	if b.Value != value {
		return fmt.Errorf("b.Value not %t. got=%t", value, b.Value)
	}
	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		return fmt.Errorf("b.TokenLiteral() not %t. got=%s", value, b.TokenLiteral())
	}

	return nil
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	err := testInfixExpression(t, exp.Condition, "x", "<", "y")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	err = testIdentifier(t, consequence.Expression, "x")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"

	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Body does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("stmt.Expression is not ast.Expression. got=%T", stmt.Expression)
	}

	err := testInfixExpression(t, exp.Condition, "x", "<", "y")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not %d statements. got=%d\n", 1, len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	err = testIdentifier(t, consequence.Expression, "x")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	alt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Alternative.Statements[0])
	}

	err = testIdentifier(t, alt.Expression, "y")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x,y) { x + y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("stmt.Expression is not ast.Expression. got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}

	err := testLiteralExpression(t, function.Parameters[0], "x")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	err = testLiteralExpression(t, function.Parameters[1], "y")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}

	err = testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected []string
	}{
		{"fn() {} case", "fn() {}", []string{}},
		{"fn(x) {}; case", "fn(x) {};", []string{"x"}},
		{"fn(x, y, z) {}; case", "fn(x, y, z) {};", []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			function := stmt.Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.expected) {
				t.Errorf("length parameters wrong. want=%d, got=%d", len(tt.expected), len(function.Parameters))
			}

			for i, ident := range tt.expected {
				err := testLiteralExpression(t, function.Parameters[i], ident)
				if err != nil {
					t.Errorf("[Error] %v", err)
				}
			}
		})
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}

	err := testIdentifier(t, exp.Function, "add")
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d,", len(exp.Arguments))
	}

	err = testLiteralExpression(t, exp.Arguments[0], 1)
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	err = testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	if err != nil {
		t.Errorf("[Error] %v", err)
	}

	err = testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
	if err != nil {
		t.Errorf("[Error] %v", err)
	}
}

/* func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected string
	}{
		{
			"-a * b case",
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a case",
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c case",
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c case",
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c case",
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c case",
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c case",
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f case",
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5 case",
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4 case",
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4 case",
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5 case",
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true case",
			"true",
			"true",
		},
		{
			"false case",
			"false",
			"false",
		},
		{
			"3 > 5 == false case",
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true case",
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4 case",
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2 case",
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5) case",
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5) case",
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true) case",
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d case",
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8)) case",
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g) case",
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(t, p)

			actual := program.String()
			if actual != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, actual)
			}
		})
	}
}
*/
