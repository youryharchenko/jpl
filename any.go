package jpl

import "fmt"

// Constructor -
type Constructor func(cls AnyClass, args []Expr, ctxName string) Expr

// Method -
type Method func(any *Any, args []Expr, ctxName string) Expr

func anyClasses() map[string]AnyClass {
	return map[string]AnyClass{
		"TviewApp":     makeTviewApp(),
		"TviewBox":     makeTviewBox(),
		"TviewPages":   makeTviewPages(),
		"TviewModal":   makeTviewModal(),
		"HttpServer":   makeHTTPServer(),
		"HttpRequest":  makeHTTPRequest(),
		"HttpClient":   makeHTTPClient(),
		"HttpResponse": makeHTTPResponse(),
	}
}

// AnyClass -
type AnyClass struct {
	Name        string
	Constructor Constructor
	Methods     map[string]Method
	Properties  map[string]Property
}

func (cls AnyClass) adaptDict(e Expr) (map[string]interface{}, error) {
	dict, ok := e.(*Dict)
	if !ok {
		return nil, fmt.Errorf("adaptDict: expr %v is not *Dict", e)
	}
	props := map[string]interface{}{}
	for key, item := range cls.Properties {
		i := item.Adapter(dict, key)
		if i == nil {
			if item.Default == nil {
				return nil, fmt.Errorf("adaptDict: value %s and default is nil", key)
			}
			i = item.Default
		}
		props[key] = i
	}
	return props, nil
}

// Property -
type Property struct {
	Name    string
	Adapter func(*Dict, string) interface{}
	Default interface{}
}

func anyToInterface(dict *Dict, name string) interface{} {
	e, ok := dict.Value[name]
	if !ok {
		return nil
	}
	any, ok := e.(*Any)
	if !ok {
		return nil
	}
	return any.Value
}

func textToString(dict *Dict, name string) interface{} {
	e, ok := dict.Value[name]
	if !ok {
		return nil
	}
	text, ok := e.(*Text)
	if !ok {
		return nil
	}
	return text.Value
}

func exprToExpr(dict *Dict, name string) interface{} {
	e, ok := dict.Value[name]
	if !ok {
		return nil
	}
	return e
}

func alistToStrings(dict *Dict, name string) interface{} {
	e, ok := dict.Value[name]
	if !ok {
		return nil
	}
	alist, ok := e.(*Alist)
	if !ok {
		return nil
	}
	list := make([]string, len(alist.Value))
	for i, button := range alist.Value {
		switch t := button.(type) {
		case *Text:
			list[i] = t.Value
		case *ID:
			list[i] = t.Value
		default:
			list[i] = undefID.Value
		}
	}
	return list
}

func idToBool(dict *Dict, name string) interface{} {
	e, ok := dict.Value[name]
	if !ok {
		return nil
	}
	id, ok := e.(*ID)
	if !ok {
		return nil
	}
	return !id.Equals(falseID)
}

func dictToMap(dict *Dict, name string) interface{} {
	e, ok := dict.Value[name]
	if !ok {
		return nil
	}
	d, ok := e.(*Dict)
	if !ok {
		return nil
	}
	m := map[string]string{}
	for key, button := range d.Value {
		//engine.debug(key)
		switch t := button.(type) {
		case *Text:
			m[key] = t.Value
		case *ID:
			m[key] = t.Value
		default:
			m[key] = undefID.Value
		}
	}
	return m
}
