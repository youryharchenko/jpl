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
