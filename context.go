package jpl

import (
	parsec "github.com/prataprc/goparsec"
)

// EvalNodes -
func EvalNodes(nodes []parsec.ParsecNode) {
	for _, node := range nodes {
		//debug("evalNodes", node)
		switch node.(type) {
		case []parsec.ParsecNode:
			v := node.([]parsec.ParsecNode)
			EvalNodes(v)
		default:
			expr := nodeToExpr(node)
			res := expr.Eval()
			debug("expr:", expr, "=>", res)
		}
	}
}
