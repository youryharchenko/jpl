package jpl

import "sync"

//var global = &Context{parent: nil, vars: map[string]Expr{}}
//var current = global
//var contextLock sync.RWMutex

// Context -
type Context struct {
	parent      *Context
	vars        map[string]Expr
	contextLock sync.RWMutex
}

func (ctx *Context) push(vars map[string]Expr, ctxName string) {
	engine.treeLock.Lock()
	c, _ := engine.current.Load(ctxName)
	nc := &Context{parent: c.(*Context), vars: vars}
	engine.current.Store(ctxName, nc)
	//engine.current[ctxName] = &Context{parent: engine.current[ctxName], vars: vars}
	engine.treeLock.Unlock()
}

func (ctx *Context) pop(ctxName string) {
	engine.treeLock.Lock()
	c, _ := engine.current.Load(ctxName)
	engine.current.Store(ctxName, c.(*Context).parent)
	//engine.current[ctxName] = engine.current[ctxName].parent
	engine.treeLock.Unlock()
}

func (ctx *Context) set(id string, val Expr) Expr {
	//engine.contextLock.Lock()
	//defer engine.contextLock.Unlock()

	c := ctx
	for c.parent != nil {
		c.contextLock.RLock()
		old, ok := c.vars[id]
		c.contextLock.RUnlock()
		if ok {
			c.contextLock.Lock()
			c.vars[id] = val
			c.contextLock.Unlock()
			//engine.debug("set", val)
			if old == nil {
				old = nullID
			}
			return old
		}
		c = c.parent
	}
	c.contextLock.RLock()
	old, ok := c.vars[id]
	c.contextLock.RUnlock()
	if !ok {
		old = undefID
	}
	c.contextLock.Lock()
	c.vars[id] = val
	c.contextLock.Unlock()
	return old
}

func (ctx *Context) get(id string) Expr {
	//engine.contextLock.RLock()
	//defer engine.contextLock.RUnlock()
	c := ctx
	for c.parent != nil {
		c.contextLock.RLock()
		val, ok := c.vars[id]
		c.contextLock.RUnlock()
		if ok {
			return val
		}
		c = c.parent
	}
	c.contextLock.RLock()
	val, ok := c.vars[id]
	c.contextLock.RUnlock()
	if ok {
		return val
	}
	return undefID
}

func (ctx *Context) bound(id string) bool {
	//engine.contextLock.RLock()
	//defer engine.contextLock.RUnlock()
	c := ctx
	for c != nil {
		c.contextLock.RLock()
		_, ok := c.vars[id]
		c.contextLock.RUnlock()
		if ok {
			return true
		}
		c = c.parent
	}
	return false
}

func (ctx *Context) dict() Expr {
	ctx.contextLock.RLock()
	defer ctx.contextLock.RUnlock()
	return &Dict{Value: ctx.vars, Name: "Dict"}
}

func (ctx *Context) clone() *Context {
	ctx.contextLock.RLock()
	defer ctx.contextLock.RUnlock()
	vars := map[string]Expr{}
	for key, item := range ctx.vars {
		vars[key] = item.Clone()
	}
	return &Context{parent: ctx.parent, vars: vars}
}
