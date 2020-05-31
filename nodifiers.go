package jpl

import (
	"reflect"
	"strconv"
	"strings"

	parsec "github.com/prataprc/goparsec"
)

func atomNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	switch n := ns[0].(type) {
	case *parsec.Terminal:
		//engine.debug("atomNode", *n)
		switch n.Name {
		case "INT":
			i, _ := strconv.Atoi(n.Value)
			return &Int{Node: n, Value: i, Name: "Num", CtxName: "main"}
		case "OCT":
			i, _ := strconv.ParseInt(n.Value, 0, 0)
			return &Int{Node: n, Value: int(i), Name: "Num", CtxName: "main"}
		case "HEX":
			i, err := strconv.ParseInt(n.Value, 0, 0)
			if err != nil {
				engine.debug("atomNode", "HEX", err)
			}
			return &Int{Node: n, Value: int(i), Name: "Num", CtxName: "main"}
		case "FLOAT":
			f, _ := strconv.ParseFloat(n.Value, 64)
			return &Float{Node: n, Value: f, Name: "Num", CtxName: "main"}
		case "IDENT":
			return &ID{Node: n, Value: n.Value, Name: "ID", CtxName: "main"}
		case "OP":
			return &ID{Node: n, Value: n.Value, Name: "ID", CtxName: "main"}
		}
		return n
	case []parsec.ParsecNode:
		//engine.debug("atomNode", n[0], n[1])
		return n
	case map[string]interface{}:
		return n
	case string:
		return &Text{Node: n, Value: strings.ReplaceAll(n, `"`, ""), Name: "Text", CtxName: "main"}
	case *Refer:
		return n
	default:
		engine.debug("atomNode: unknown type", reflect.TypeOf(n).String())
	}
	return nil
}

func referNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//engine.debug("referNode", ns)
	id := nodeToExpr(ns[1]).String()
	return &Refer{Node: ns[1], Value: id, Name: "Refer", CtxName: "main"}
}

func alistNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//engine.debug("alistNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	list := []Expr{}
	for _, item := range ilist {
		list = append(list, nodeToExpr(item))
	}
	return &Alist{Node: ns[1], Value: list, Name: "Alist", CtxName: "main"}
}

func mlistNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//engine.debug("mlistNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	list := []Expr{}
	for _, item := range ilist {
		list = append(list, nodeToExpr(item))
	}
	return &Mlist{Node: ns[1], Value: list, Name: "Mlist", CtxName: "main"}
}

func llistNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//engine.debug("llistNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	list := []Expr{}
	for _, item := range ilist {
		list = append(list, nodeToExpr(item))
	}
	return &Llist{Node: ns[1], Value: list, Name: "Llist", CtxName: "main"}
}

func propNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//engine.debug("propNode", ns)
	key := nodeToExpr(ns[0]).String()
	item := nodeToExpr(ns[2])
	return &Prop{Node: ns[1], Key: key, Value: item, Name: "Prop", CtxName: "main"}
}

func dictNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//engine.debug("dictNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	m := map[string]Expr{}
	for _, item := range ilist {
		prop := nodeToExpr(item).(*Prop)
		m[prop.Key] = nodeToExpr(prop.Value)
	}
	return &Dict{Node: ns[1], Value: m, Name: "Dict", CtxName: "main"}
}

/*
func commentNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	comment := ns[1].(*parsec.Terminal)
	//engine.debug("commentNode", comment.Value)
	return &Comment{Node: ns[1], Value: comment.Value, Name: "Comment"}
}
*/

func nodeToExpr(node parsec.ParsecNode) (res Expr) {
	switch node.(type) {
	case *Int:
		res = node.(*Int)
	case *Float:
		res = node.(*Float)
	case *ID:
		res = node.(*ID)
	case *Refer:
		res = node.(*Refer)
	case *Alist:
		res = node.(*Alist)
	case *Mlist:
		res = node.(*Mlist)
	case *Llist:
		res = node.(*Llist)
	case *Prop:
		res = node.(*Prop)
	case *Dict:
		res = node.(*Dict)
	case *Text:
		res = node.(*Text)
	//case *Comment:
	//	res = node.(*Comment)
	case []parsec.ParsecNode:
		//engine.debug("nodeToExpr: []parsec.ParsecNode:", node)
		nodes := node.([]parsec.ParsecNode)
		if len(nodes) == 1 {
			res = nodeToExpr(nodes[0])
		} else {
			engine.debug("nodeToExpr: []parsec.ParsecNode: len > 1", node)
		}
	case *parsec.Terminal:
		//engine.debug("nodeToExpr: *parsec.Terminal", node)
		n := node.(*parsec.Terminal)
		switch n.Name {
		case "INT":
			i, _ := strconv.Atoi(n.Value)
			res = &Int{Node: n, Value: i, Name: "Num", CtxName: "main"}
		case "FLOAT":
			f, _ := strconv.ParseFloat(n.Value, 64)
			res = &Float{Node: n, Value: f, Name: "Num", CtxName: "main"}
		case "IDENT":
			res = &ID{Node: n, Value: n.Value, Name: "ID", CtxName: "main"}
		case "OP":
			res = &ID{Node: n, Value: n.Value, Name: "ID", CtxName: "main"}
		}
	default:
		engine.debug("nodeToExpr: unknown type", reflect.TypeOf(node))
	}
	return
}

func debugNodes(nodes []parsec.ParsecNode, deep int) {

	for _, node := range nodes {
		switch node.(type) {
		case *parsec.Terminal:
			v := node.(*parsec.Terminal)
			engine.debug(strings.Repeat(".", deep*2), v.Name, v.Value, v.Position, v.Attributes)
		case *Int:
			v := node.(*Int)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *Float:
			v := node.(*Float)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *ID:
			v := node.(*ID)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *Refer:
			v := node.(*Refer)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *Alist:
			v := node.(*Alist)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *Mlist:
			v := node.(*Mlist)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *Llist:
			v := node.(*Llist)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *Dict:
			v := node.(*Dict)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case *Text:
			v := node.(*Text)
			engine.debug(strings.Repeat(".", deep*2), v.Debug())
		case []parsec.ParsecNode:
			v := node.([]parsec.ParsecNode)
			debugNodes(v, deep+1)
		default:
			engine.debug(reflect.TypeOf(node).String())
		}
	}
}
