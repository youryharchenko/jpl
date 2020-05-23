package jpl

import (
	parsec "github.com/prataprc/goparsec"
)

var global = &Context{parent: nil, vars: map[string]Expr{}}
var current = global

// Context -
type Context struct {
	parent *Context
	vars   map[string]Expr
}

func (ctx *Context) push(vars map[string]Expr) {
	current = &Context{parent: current, vars: vars}
}

func (ctx *Context) pop() {
	current = current.parent
}

func (ctx *Context) set(id string, val Expr) Expr {
	c := ctx
	for c.parent != nil {
		old, ok := c.vars[id]
		if ok {
			c.vars[id] = val
			//debug("set", val)
			if old == nil {
				old = &ID{Value: "null", Name: "ID"}
			}
			return old
		}
		c = c.parent
	}
	old, ok := c.vars[id]
	if !ok {
		old = &ID{Value: "undefined", Name: "ID"}
	}
	c.vars[id] = val
	return old
}

func (ctx *Context) get(id string) Expr {
	c := ctx
	for c.parent != nil {
		val, ok := c.vars[id]
		if ok {
			return val
		}
		c = c.parent
	}
	val, ok := c.vars[id]
	if ok {
		return val
	}
	return &ID{Value: "undefined", Name: "ID"}
}

func (ctx *Context) bound(id string) bool {
	c := ctx
	for c != nil {
		_, ok := c.vars[id]
		if ok {
			return true
		}
		c = c.parent
	}
	return false
}

func (ctx *Context) dict() Expr {
	return &Dict{Value: ctx.vars, Name: "Dict"}
}

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
