package jpl

import (
	"sync"

	parsec "github.com/prataprc/goparsec"
)

var global = &Context{parent: nil, vars: map[string]Expr{}}
var current = global
var contextLock sync.RWMutex

// Context -
type Context struct {
	parent *Context
	vars   map[string]Expr
}

func (ctx *Context) push(vars map[string]Expr) {
	contextLock.Lock()
	current = &Context{parent: current, vars: vars}
	contextLock.Unlock()
}

func (ctx *Context) pop() {
	contextLock.Lock()
	current = current.parent
	contextLock.Unlock()
}

func (ctx *Context) set(id string, val Expr) Expr {
	contextLock.Lock()
	defer contextLock.Unlock()

	c := ctx
	for c.parent != nil {
		old, ok := c.vars[id]
		if ok {
			c.vars[id] = val
			//debug("set", val)
			if old == nil {
				old = nullID
			}
			return old
		}
		c = c.parent
	}
	old, ok := c.vars[id]

	if !ok {
		old = undefID
	}
	c.vars[id] = val
	return old
}

func (ctx *Context) get(id string) Expr {
	contextLock.RLock()
	defer contextLock.RUnlock()
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
	return undefID
}

func (ctx *Context) bound(id string) bool {
	contextLock.RLock()
	defer contextLock.RUnlock()
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
	contextLock.RLock()
	defer contextLock.RUnlock()
	return &Dict{Value: ctx.vars, Name: "Dict"}
}

func (ctx *Context) clone() *Context {
	contextLock.RLock()
	defer contextLock.RUnlock()
	vars := map[string]Expr{}
	for key, item := range ctx.vars {
		vars[key] = item.Clone()
	}
	return &Context{parent: ctx.parent, vars: vars}
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
