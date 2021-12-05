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
		{"5 case", "5", 5},
		{"10 case", "10", 10},
		{"-5 case", "-5", -5},
		{"-10", "-10", -10},
		{"5 + 5 + 5 + 5 - 10 case", "5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2 case", "2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50 case", "-50 + 100 + -50", 0},
		{"5 * 2 + 10 case", "5 * 2 + 10", 20},
		{"5 + 2 * 10 case", "5 + 2 * 10", 25},
		{"20 + 2 * -10 case", "20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10 case", "50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10) case", "2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10 case", "3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10 case", "3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10 case", "(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
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

	return Eval(program)
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
