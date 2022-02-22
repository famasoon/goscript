package evaluator

import (
	"fmt"
	"goscript/lexer"
	"goscript/object"
	"goscript/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected int64
	}{
		{"5", "5", 5},
		{"10", "10", 10},
		{"-5", "-5", -5},
		{"-10", "-10", -10},
		{"5 + 5", "5+5", 10},
		{"5 + 5 + 5 + 5 - 10", "5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", "2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", "-50 + 100 + -50", 0},
		{"5 * 2 + 10", "5 * 2 + 10", 20},
		{"5 + 2 * 10", "5 + 2 * 10", 25},
		{"20 + 2 * -10", "20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", "50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", "2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", "3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", "3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			evaluated := testEval(tt.input)
			err := testIntegerObject(evaluated, tt.expected)
			if err != nil {
				t.Errorf("[ERROR] %v", err)
			}
		})
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(obj object.Object, expected int64) error {
	result, ok := obj.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func TestEvaBooleanExpression(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected bool
	}{
		{"true case", "true", true},
		{"false case", "false", false},
		{"1 < 2", "1 < 2", true},
		{"1 > 2", "1 > 2", false},
		{"1 < 1", "1 < 1", false},
		{"1 > 1", "1 > 1", false},
		{"1 == 1", "1 == 1", true},
		{"1 != 1", "1 != 1", false},
		{"1 == 2", "1 == 2", false},
		{"1 != 2", "1 != 2", true},
		{"true == true", "true == true", true},
		{"false == false", "false == false", true},
		{"true == false", "true == false", false},
		{"true != false", "true != false", true},
		{"false != true", "false != true", true},
		{"(1 < 2) == true", "(1 < 2) == true", true},
		{"(1 < 2) == false", "(1 < 2) == false", false},
		{"(1 > 2) == true", "(1 > 2) == true", false},
		{"(1 > 2) == false", "(1 > 2) == false", true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			evaluated := testEval(tt.input)
			err := testBooleanObject(evaluated, tt.expected)
			if err != nil {
				t.Errorf("[ERROR] %v", err)
			}
		})
	}
}

func testBooleanObject(obj object.Object, expected bool) error {
	result, ok := obj.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected bool
	}{
		{"!true case", "!true", false},
		{"!false case", "!false", true},
		{"!5 case", "!5", false},
		{"!!true case", "!!true", true},
		{"!!false case", "!!false", false},
		{"!!5 case", "!!5", true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			evaluated := testEval(tt.input)
			err := testBooleanObject(evaluated, tt.expected)
			if err != nil {
				t.Errorf("[ERROR] %v", err)
			}
		})
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", "if (true) { 10 }", 10},
		{"if (false) { 10 }", "if (false) { 10 }", nil},
		{"if (1) { 10 }", "if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", "if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", "if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", "if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", "if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				err := testIntegerObject(evaluated, int64(integer))
				if err != nil {
					t.Errorf("[ERROR] %v", err)
				}
			} else {
				err := testNullObject(evaluated)
				if err != nil {
					t.Errorf("[ERROR] %v", err)
				}
			}
		})
	}
}

func testNullObject(obj object.Object) error {
	if obj != NULL {
		return fmt.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
	}
	return nil
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected int64
	}{
		{"return 10;", "return 10;", 10},
		{"return 10; 9;", "return 10; 9;", 10},
		{"return 2 * 5; 9;", "return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", "9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { return 10; }", "if (10 > 1) { return 10; }", 10},
		{"if case",
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			evaluated := testEval(tt.input)
			err := testIntegerObject(evaluated, tt.expected)
			if err != nil {
				t.Errorf("[ERROR] %v", err)
			}
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		testName        string
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{"return true + false;",
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World" case`,
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			evaluated := testEval(tt.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			}
			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected int64
	}{
		{"a case", "let a = 5; a;", 5},
		{"b case", "let a = 5 * 5; a;", 25},
		{"ab case", "let a = 5; let b = a; b;", 5},
		{"abc case", "let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := testIntegerObject(testEval(tt.input), tt.expected)
			if err != nil {
				t.Errorf("[ERROR] %v", err)
			}
		})
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		testName string
		input    string
		expected int64
	}{
		{"1st case", "let identity = fn(x) { x; }; identity(5);", 5},
		{"2nd case", "let identity = fn(x) { return x; }; identity(5);", 5},
		{"3rd case", "let double = fn(x) { x * 2; }; double(5);", 10},
		{"4th case", "let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"5th case", "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"6th case", "fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := testIntegerObject(testEval(tt.input), tt.expected)
			if err != nil {
				t.Errorf("[ERROR] %v", err)
			}
		})
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			evaluated := testEval(tt.input)

			switch expected := tt.expected.(type) {
			case int:
				err := testIntegerObject(evaluated, int64(expected))
				if err != nil {
					t.Errorf("[ERROR] %v", err)
				}
			case string:
				errObj, ok := evaluated.(*object.Error)
				if !ok {
					t.Errorf("object is not Error. got=%T (%+v)",
						evaluated, evaluated)
				}
				if errObj.Message != expected {
					t.Errorf("wrong error message. expected=%q, got=%q",
						expected, errObj.Message)
				}
			}
		})
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	err := testIntegerObject(result.Elements[0], 1)
	if err != nil {
		t.Errorf("[ERROR] %v", err)
	}
	err = testIntegerObject(result.Elements[1], 4)
	if err != nil {
		t.Errorf("[ERROR] %v", err)
	}
	err = testIntegerObject(result.Elements[2], 6)
	if err != nil {
		t.Errorf("[ERROR] %v", err)
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		name string
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T){
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				err := testIntegerObject(evaluated, int64(integer))
				if err != nil {
					t.Errorf("[ERROR] %v\n", err)
				}
			} else {
				err := testNullObject(evaluated)
				if err != nil {
					t.Errorf("[ERROR] %v\n", err)
				}
			}
		})
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		err := testIntegerObject(pair.Value, expectedValue)
		if err != nil {
			t.Errorf("[ERROR] %v\n", err)
		}
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		name string
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T)  {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				err := testIntegerObject(evaluated, int64(integer))
				if err != nil {
					t.Errorf("[ERROR] %v\n", err)
				}
			} else {
				err := testNullObject(evaluated)
				if err != nil {
					t.Errorf("[ERROR] %v\n", err)
				}
			}
		})
	}
}