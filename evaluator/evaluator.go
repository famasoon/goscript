package evaluator

import (
	"goscript/ast"
	"goscript/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
	}

	return nil
}

func evalStatements(stmt []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmt {
		result = Eval(statement)
	}

	return result
}
