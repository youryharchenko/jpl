package jpl

import (
	"fmt"
	"math"
)

func mathFuncs() map[string]Func {
	return map[string]Func{
		"+":     sum,
		"-":     sub,
		"*":     prod,
		"/":     div,
		"%":     mod,
		"pow":   pow,
		"abs":   abs,
		"gt":    gt,
		"lt":    lt,
		"ge":    ge,
		"le":    le,
		"range": rangeInt,
		"int":   toInt,
		"float": toFloat,
	}
}

func toInt(args []Expr, ctxName string) Expr {
	if len(args) != 1 {
		return errID
	}
	switch a := args[0].Eval().(type) {
	case *Int:
		return a
	case *Float:
		return &Int{Name: "Int", Value: int(a.Value), CtxName: ctxName}
	default:
		return errID
	}
}

func toFloat(args []Expr, ctxName string) Expr {
	if len(args) != 1 {
		return errID
	}
	switch a := args[0].Eval().(type) {
	case *Float:
		return a
	case *Int:
		return &Float{Name: "Float", Value: float64(a.Value), CtxName: ctxName}
	default:
		return errID
	}
}

func sum(args []Expr, ctxName string) Expr {
	if len(args) < 2 {
		return errID
	}
	bInt := true
	si := 0
	sf := 0.0
	for _, arg := range args {
		a := arg.Eval()
		switch a.(type) {
		case *Int:
			si += a.(*Int).Value
		case *Float:
			sf += a.(*Float).Value
			bInt = false
		}
	}
	if bInt {
		return &Int{Value: si, Name: "Num"}
	}
	sf += float64(si)
	return &Float{Value: sf, Name: "Num", CtxName: ctxName}
}

func sub(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	bInt := true
	var si int
	var sf float64
	for i, arg := range args[:2] {
		a := arg.Eval()
		switch a.(type) {
		case *Int:
			if i == 0 {
				si = a.(*Int).Value
				sf = float64(a.(*Int).Value)
			} else {
				si -= a.(*Int).Value
				sf -= float64(a.(*Int).Value)
			}
		case *Float:
			if i == 0 {
				sf = a.(*Float).Value
			} else {
				sf -= a.(*Float).Value
			}
			bInt = false
		}
	}
	if bInt {
		return &Int{Value: si, Name: "Num", CtxName: ctxName}
	}
	return &Float{Value: sf, Name: "Num", CtxName: ctxName}
}

func prod(args []Expr, ctxName string) Expr {
	if len(args) < 2 {
		return errID
	}
	bInt := true
	si := 1
	sf := 1.0
	for _, arg := range args {
		a := arg.Eval()
		switch a.(type) {
		case *Int:
			si *= a.(*Int).Value
		case *Float:
			sf *= a.(*Float).Value
			bInt = false
		}
	}
	if bInt {
		return &Int{Value: si, Name: "Num", CtxName: ctxName}
	}
	sf *= float64(si)
	return &Float{Value: sf, Name: "Num", CtxName: ctxName}
}

func div(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	var sf float64
	for i, arg := range args[:2] {
		a := arg.Eval()
		switch a.(type) {
		case *Int:
			if i == 0 {
				sf = float64(a.(*Int).Value)
			} else {
				sf /= float64(a.(*Int).Value)
			}
		case *Float:
			if i == 0 {
				sf = a.(*Float).Value
			} else {
				sf /= a.(*Float).Value
			}
		}
	}
	return &Float{Value: sf, Name: "Num", CtxName: ctxName}
}

func mod(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	var si int
	for i, arg := range args[:2] {
		a := arg.Eval()
		switch a.(type) {
		case *Int:
			if i == 0 {
				si = a.(*Int).Value
			} else {
				si %= a.(*Int).Value
			}
		case *Float:
			if i == 0 {
				si = int(a.(*Float).Value)
			} else {
				si %= int(a.(*Float).Value)
			}
		}
	}
	return &Int{Value: si, Name: "Num", CtxName: ctxName}
}

func pow(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	bInt := true
	var x float64
	var y float64
	for i, arg := range args[:2] {
		a := arg.Eval()
		switch a.(type) {
		case *Int:
			if i == 0 {
				x = float64(a.(*Int).Value)
			} else {
				y = float64(a.(*Int).Value)
			}
		case *Float:
			if i == 0 {
				x = a.(*Float).Value
			} else {
				y = a.(*Float).Value
			}
			bInt = false
		}
	}
	if bInt {
		return &Int{Value: int(math.Pow(x, y)), Name: "Num", CtxName: ctxName}
	}
	return &Float{Value: math.Pow(x, y), Name: "Num"}
}

func abs(args []Expr, ctxName string) Expr {
	if len(args) != 1 {
		return errID
	}
	switch a := args[0].Eval().(type) {
	case *Int:
		if a.Value < 0 {
			return &Int{Name: "Int", Value: a.Value * -1, CtxName: ctxName}
		}
		return a
	case *Float:
		if a.Value < 0 {
			return &Float{Name: "Float", Value: a.Value * -1.0, CtxName: ctxName}
		}
		return a
	default:
		return errID
	}
}

func compareFloat(e1 Expr, e2 Expr, f func(x float64, y float64) bool) (bool, error) {
	f1, ok := e1.(*Float)
	if !ok {
		return false, fmt.Errorf("first arg not Float")
	}
	f2, ok := e2.(*Float)
	if !ok {
		return false, fmt.Errorf("second arg not Float")
	}
	return f(f1.Value, f2.Value), nil
}

func compareToFloat(e1 Expr, e2 Expr, f func(x float64, y float64) bool) (bool, error) {
	ef1, ok := e1.(*Float)
	var f1 float64
	if !ok {
		i1, ok := e1.(*Int)
		if !ok {
			return false, fmt.Errorf("first arg not Float abd not Int")
		}
		f1 = float64(i1.Value)
	} else {
		f1 = ef1.Value
	}
	ef2, ok := e2.(*Float)
	var f2 float64
	if !ok {
		i2, ok := e2.(*Int)
		if !ok {
			return false, fmt.Errorf("second arg not Float abd not Int")
		}
		f2 = float64(i2.Value)
	} else {
		f2 = ef2.Value
	}
	return f(f1, f2), nil
}

func compareInt(e1 Expr, e2 Expr, f func(x int, y int) bool) (bool, error) {
	i1, ok := e1.(*Int)
	if !ok {
		return false, fmt.Errorf("first arg not Int")
	}
	i2, ok := e2.(*Int)
	if !ok {
		return false, fmt.Errorf("second arg not Int")
	}
	return f(i1.Value, i2.Value), nil
}

func compareNum(e1 Expr, e2 Expr, fInt func(x int, y int) bool, fFloat func(x float64, y float64) bool) (bool, error) {
	b, err := compareInt(e1, e2, fInt)
	if err == nil {
		if b {
			return true, nil
		}
		return false, nil
	}
	b, err = compareFloat(e1, e2, fFloat)
	if err == nil {
		if b {
			return true, nil
		}
		return false, nil
	}
	b, err = compareToFloat(e1, e2, fFloat)
	if err == nil {
		if b {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("error compareNum")
}

func gt(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	fInt := func(x int, y int) bool { return x > y }
	fFloat := func(x float64, y float64) bool { return x > y }

	b, err := compareNum(args[0].Eval(), args[1].Eval(), fInt, fFloat)
	if err == nil {
		if b {
			return trueID
		}
		return falseID
	}
	return errID
}

func lt(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	fInt := func(x int, y int) bool { return x < y }
	fFloat := func(x float64, y float64) bool { return x < y }

	b, err := compareNum(args[0].Eval(), args[1].Eval(), fInt, fFloat)
	if err == nil {
		if b {
			return trueID
		}
		return falseID
	}
	return errID
}

func ge(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	fInt := func(x int, y int) bool { return x >= y }
	fFloat := func(x float64, y float64) bool { return x >= y }

	b, err := compareNum(args[0].Eval(), args[1].Eval(), fInt, fFloat)
	if err == nil {
		if b {
			return trueID
		}
		return falseID
	}
	return errID
}

func le(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	fInt := func(x int, y int) bool { return x <= y }
	fFloat := func(x float64, y float64) bool { return x <= y }

	b, err := compareNum(args[0].Eval(), args[1].Eval(), fInt, fFloat)
	if err == nil {
		if b {
			return trueID
		}
		return falseID
	}
	return errID
}

func rangeInt(args []Expr, ctxName string) Expr {
	if len(args) != 2 {
		return errID
	}
	i1, ok1 := args[0].(*Int)
	i2, ok2 := args[1].(*Int)
	if !(ok1 && ok2) {
		return errID
	}
	list := []Expr{}
	if i2.Value >= i1.Value {
		list = make([]Expr, i2.Value-i1.Value)
		for i := i1.Value; i < i2.Value; i++ {
			//list = append(list, &Int{Name: "Num", Node: nil, Value: i})
			list[i-i1.Value] = &Int{Name: "Num", Node: nil, Value: i, CtxName: ctxName}
		}
	} else {
		list = make([]Expr, i1.Value-i2.Value)
		cnt := 0
		for i := i1.Value; i > i2.Value; i-- {
			//list = append(list, &Int{Name: "Num", Node: nil, Value: i})
			list[cnt] = &Int{Name: "Num", Node: nil, Value: i, CtxName: ctxName}
			cnt++
		}
	}
	return &Alist{Name: "Alist", Node: nil, Value: list, CtxName: ctxName}
}
