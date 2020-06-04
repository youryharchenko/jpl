package jpl

import (
	"encoding/json"
	"reflect"
	"strings"
)

func jsonFuncs() map[string]Func {
	return map[string]Func{
		"toJSON":   toJSON,
		"fromJSON": fromJSON,
	}
}

func toJSON(args []Expr, ctxName string) Expr {
	if len(args) != 1 {
		return errID
	}
	e := args[0].Eval()
	b, err := json.Marshal(e)
	if err != nil {

	}
	return &Text{Name: "Text", Value: string(b), CtxName: ctxName}
}

func fromJSON(args []Expr, ctxName string) Expr {
	if len(args) != 1 {
		return errID
	}
	text, ok := args[0].Eval().(*Text)
	if !ok {
		return errID
	}

	dec := json.NewDecoder(strings.NewReader(text.Value))
	t, err := dec.Token()
	if err != nil {
		return &Text{Name: "Text", Value: "token error:" + err.Error(), CtxName: ctxName}
	}
	var v Expr
	var j []byte

	switch t.(type) {
	case float64:
		//engine.debug(reflect.TypeOf(t), t)
		v = &Float{Name: "Num", CtxName: ctxName}
		j = []byte(text.Value)
	case string:
		//engine.debug(reflect.TypeOf(t), t)
		v = &Text{Name: "Text", CtxName: ctxName}
		j = []byte(text.Value)
	case bool:
		//engine.debug(reflect.TypeOf(t), t)
		v = &ID{Name: "ID", CtxName: ctxName}
		j = []byte(text.Value)
	case json.Delim:
		engine.debug(reflect.TypeOf(t), t)
		switch t.(json.Delim).String() {
		case "[":
			v = &Alist{Name: "Alist", CtxName: ctxName}
			j = []byte(text.Value)
		case "{":
			v = &Dict{Name: "Dict", CtxName: ctxName}
			j = []byte(text.Value)
		}
	default:
		engine.debug(reflect.TypeOf(t), t)
		if t == nil {
			return &ID{Name: "ID", Value: "null", CtxName: ctxName}
		}
	}
	err = json.Unmarshal(j, v)
	if err != nil {
		return &Text{Name: "Text", Value: "error:" + err.Error(), CtxName: ctxName}
	}
	return v
}

func itemToExpr(item interface{}, ctxName string) Expr {
	//engine.debug(reflect.TypeOf(item), item)
	if item == nil {
		return nullID
	}
	switch val := item.(type) {
	case float64:
		return &Float{Name: "Float", Value: val, CtxName: ctxName}
	case string:
		return &Text{Name: "Text", Value: val, CtxName: ctxName}
	case bool:
		if val {
			return &ID{Name: "ID", Value: TRUE, CtxName: ctxName}
		}
		return &ID{Name: "ID", Value: FALSE, CtxName: ctxName}
	case []interface{}:
		list := make([]Expr, len(val))
		for i, item := range val {
			list[i] = itemToExpr(item, ctxName)
		}
		return &Alist{Name: "Alist", Value: list, CtxName: ctxName}
	case map[string]interface{}:
		dict := map[string]Expr{}
		for key, item := range val {
			dict[key] = itemToExpr(item, ctxName)
		}
		return &Dict{Name: "Dict", Value: dict, CtxName: ctxName}
	default:
		return undefID
	}
	//return undefID
}
