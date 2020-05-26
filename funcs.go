package jpl

import (
	"log"
	"reflect"
)

// Constants -
const (
	TRUE      = "true"
	FALSE     = "false"
	ERROR     = "error"
	NULL      = "null"
	UNDEFINED = "undefined"
	BREAK     = "break"
	CONTINUE  = "contunue"
)

var (
	trueID  = &ID{Value: TRUE, Name: "ID"}
	falseID = &ID{Value: FALSE, Name: "ID"}
	errID   = &ID{Value: ERROR, Name: "ID"}
	nullID  = &ID{Value: NULL, Name: "ID"}
	undefID = &ID{Value: UNDEFINED, Name: "ID"}
)

func initFuncs() map[string]Func {
	return mergeFuncs(funcs, coreFuncs, osFuncs, mathFuncs)
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

func applyFunc(fn Expr, args []Expr) Expr {
	switch fnExpr := fn.Eval().(type) {
	case *ID:
		name := fnExpr.Value
		f, ok := funcs[name]
		if !ok {
			return undefID
		}
		return f(args)
	case *Lamb:
		return fnExpr.Apply(args)
	}
	return undefID
}

var coreFuncs = map[string]Func{
	"print": printExprs,
	"quote": quote,
	"eval":  eval,
	"set":   set,
	"get":   get,
	"put":   put,
	"let":   let,
	"do":    do,
	"while": while,
	"eq":    eq,
	"is":    is,
	"not":   not,
	"if":    iff,
	"func":  lambda,
	"form":  form,
	"bool":  toBool,
	"map":   mapl,
	"fold":  foldl,
	//"foldl": foldl,
	//"foldr": foldr,
}

func printExprs(args []Expr) Expr {
	for _, arg := range args {
		log.Println(arg.Eval())
	}
	return &Int{Value: len(args), Name: "Num"}
}

func quote(args []Expr) Expr {
	if len(args) == 0 {
		return nullID
	}
	return args[0]
}

func eval(args []Expr) Expr {
	if len(args) == 0 {
		return nullID
	}
	return args[0].Eval().Eval()
}

func set(args []Expr) Expr {
	if len(args) != 2 {
		return errID
	}
	return current.set(args[0].Eval().String(), args[1].Eval())
}

func get(args []Expr) Expr {
	if len(args) != 2 {
		return errID
	}
	dict, ok := args[0].Eval().(*Dict)
	if !ok {
		return errID
	}
	key, ok := args[1].Eval().(*ID)
	if !ok {
		return errID
	}
	val, ok := dict.Value[key.Value]
	if !ok {
		return nullID
	}
	return val
}

func put(args []Expr) Expr {
	if len(args) != 3 {
		return errID
	}
	dict, ok := args[0].Eval().(*Dict)
	if !ok {
		return errID
	}
	key, ok := args[1].Eval().(*ID)
	if !ok {
		return errID
	}
	var res Expr = nullID
	res, ok = dict.Value[key.Value]
	if !ok {
		res = nullID
	}
	dict.Value[key.Name] = args[3].Eval()
	return res
}

func eq(args []Expr) Expr {
	if len(args) != 2 {
		return errID
	}
	if args[0].Eval().Equals(args[1].Eval()) {
		return trueID
	}
	return falseID
}

func is(args []Expr) Expr {
	if len(args) != 2 {
		return errID
	}
	return match(args[0], args[1].Eval())
}

func not(args []Expr) Expr {
	if len(args) != 1 {
		return errID
	}
	e := args[0].Eval()
	if e.Equals(trueID) {
		return falseID
	} else if e.Equals(falseID) {
		return trueID
	}
	return errID
}

func iff(args []Expr) Expr {
	if len(args) < 2 || len(args) > 3 {
		return errID
	}
	cond, ok := args[0].Eval().(*ID)
	if !ok {
		return errID
	}
	if cond.Value == TRUE {
		return args[1].Eval()
	}
	if cond.Value == FALSE && len(args) == 3 {
		return args[2].Eval()
	}
	return errID
}

func let(args []Expr) Expr {
	if len(args) < 1 {
		return errID
	}
	d, ok := args[0].Eval().(*Dict)
	if !ok {
		return errID
	}
	var res Expr
	current.push(d.Value)
	do(args[1:])
	res = current.dict()
	current.pop()
	return res
}

func while(args []Expr) (res Expr) {
	if len(args) < 2 {
		return errID
	}
	res = nullID
	for args[0].Eval().Equals(trueID) {
		res = do(args[1:])
		id, ok := res.(*ID)
		if ok && id.Value == BREAK {
			break
		}
		if ok && id.Value == CONTINUE {
			continue
		}
	}
	return res
}

func do(args []Expr) Expr {
	var res Expr = nullID
	for _, item := range args {
		res = item.Eval()
		id, ok := res.(*ID)
		if ok && id.Value == BREAK {
			break
		}
		if ok && id.Value == CONTINUE {
			break
		}
	}
	return res
}

func lambda(args []Expr) Expr {
	if len(args) != 2 {
		return errID
	}
	alist, ok := args[0].(*Alist)
	if !ok {
		return errID
	}
	params := []*ID{}
	for _, item := range alist.Value {
		param, ok := item.Eval().(*ID)
		if !ok {
			return errID
		}
		params = append(params, param)
	}
	body := args[1]
	return &Lamb{Params: params, Body: body, Name: "Lambda"}
}

func form(args []Expr) Expr {
	if len(args) != 1 {
		return errID
	}
	var res Expr
	switch e := args[0].(type) {
	case *Llist:
		list := []Expr{}
		for _, item := range e.Value {
			list = append(list, item.Eval())
		}
		res = &Llist{Name: e.Name, Node: e.Node, Value: list}
	case *Mlist:
		list := []Expr{}
		for _, item := range e.Value {
			list = append(list, item.Eval())
		}
		res = &Mlist{Name: e.Name, Node: e.Node, Value: list}
	default:
		res = e.Eval()
	}
	return res
}

func toBool(args []Expr) Expr {
	if len(args) != 1 {
		return errID
	}
	res := trueID
	switch e := args[0].Eval().(type) {
	case *ID:
		if e.Equals(falseID) || e.Equals(nullID) || e.Equals(undefID) {
			res = falseID
		}
	case *Int:
		if e.Value == 0 {
			res = falseID
		}
	case *Float:
		if e.Value == 0.0 {
			res = falseID
		}
	case *Alist:
		if len(e.Value) == 0 {
			res = falseID
		}
	case *Mlist:
		if len(e.Value) == 0 {
			res = falseID
		}
	case *Dict:
		if len(e.Value) == 0 {
			res = falseID
		}
	case *Text:
		//debug("toBool", "text", len(e.Value), e.Value)
		if len(e.Value) == 0 {
			res = falseID
		}
	default:
		debug("toBool", reflect.TypeOf(e))
	}
	return res
}

func mapl(args []Expr) Expr {
	if len(args) != 2 {
		return errID
	}
	alist, ok := args[0].Eval().(*Alist)
	if !ok {
		return errID
	}
	list := make([]Expr, len(alist.Value))
	for i, item := range alist.Eval().(*Alist).Value {
		//list = append(list, applyFunc(args[1].Eval(), []Expr{item}))
		list[i] = applyFunc(args[1].Eval(), []Expr{item})
	}
	return &Alist{Node: alistNode, Name: alist.Name, Value: list}
}

func foldl(args []Expr) Expr {
	if len(args) != 3 {
		return errID
	}
	alist, ok := args[0].Eval().(*Alist)
	if !ok {
		return errID
	}
	var res Expr = args[1].Eval()
	for _, item := range alist.Eval().(*Alist).Value {
		res = applyFunc(args[2].Eval(), []Expr{res, item})
	}
	return res
}
