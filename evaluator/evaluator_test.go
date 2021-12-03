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
