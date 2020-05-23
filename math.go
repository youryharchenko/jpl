package jpl

import "math"

var mathFuncs = map[string]Func{
	"+":   sum,
	"-":   sub,
	"*":   prod,
	"/":   div,
	"%":   mod,
	"pow": pow,
}

func sum(args []Expr) Expr {
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
	return &Float{Value: sf, Name: "Num"}
}

func sub(args []Expr) Expr {
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
		return &Int{Value: si, Name: "Num"}
	}
	return &Float{Value: sf, Name: "Num"}
}

func prod(args []Expr) Expr {
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
		return &Int{Value: si, Name: "Num"}
	}
	sf *= float64(si)
	return &Float{Value: sf, Name: "Num"}
}

func div(args []Expr) Expr {
	var sf float64
	for i, arg := range args[:2] {
		a := arg.Eval()
		switch a.(type) {
		case *Int:
			if i == 0 {
				sf = float64(a.(*Int).Value)
			} else {
				sf -= float64(a.(*Int).Value)
			}
		case *Float:
			if i == 0 {
				sf = a.(*Float).Value
			} else {
				sf -= a.(*Float).Value
			}
		}
	}
	return &Float{Value: sf, Name: "Num"}
}

func mod(args []Expr) Expr {
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
	return &Int{Value: si, Name: "Num"}
}

func pow(args []Expr) Expr {
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
		return &Int{Value: int(math.Pow(x, y)), Name: "Num"}
	}
	return &Float{Value: math.Pow(x, y), Name: "Num"}
}
