package jpl

import (
	"log"
)

func initFuncs() map[string]Func {
	return mergeFuncs(funcs, coreFuncs, mathFuncs)
}

func mergeFuncs(fns map[string]Func, ext ...map[string]Func) map[string]Func {
	if fns == nil {
		fns = map[string]Func{}
	}
	for _, item := range ext {
		for key, fn := range item {
			fns[key] = fn
		}
	}
	return fns
}

var coreFuncs = map[string]Func{
	"print": printExprs,
	"quote": quote,
	"eval":  eval,
	"set":   set,
	"let":   let,
}

func printExprs(args []Expr) Expr {
	for _, arg := range args {
		log.Println(arg.Eval())
	}
	return &Int{Value: len(args), Name: "Num"}
}

func quote(args []Expr) Expr {
	if len(args) == 0 {
		return &ID{Value: "null", Name: "ID"}
	}
	return args[0]
}

func eval(args []Expr) Expr {
	if len(args) == 0 {
		return &ID{Value: "null", Name: "ID"}
	}
	return args[0].Eval().Eval()
}

func set(args []Expr) Expr {
	if len(args) != 2 {
		return &ID{Value: "error", Name: "ID"}
	}
	return current.set(args[0].Eval().String(), args[1].Eval())
}

func let(args []Expr) Expr {
	errID := &ID{Value: "error", Name: "ID"}
	if len(args) < 1 {
		return errID
	}
	d, ok := args[0].Eval().(*Dict)
	if !ok {
		return errID
	}
	var res Expr
	current.push(d.Value)
	for _, item := range args[1:] {
		res = item.Eval()
		id, ok := res.(*ID)
		if ok && id.Value == "break" {
			break
		}
	}
	res = current.dict()
	current.pop()
	return res
}
