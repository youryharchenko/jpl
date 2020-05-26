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
		//debug("atomNode", *n)
		switch n.Name {
		case "INT":
			i, _ := strconv.Atoi(n.Value)
			return &Int{Node: n, Value: i, Name: "Num"}
		case "OCT":
			i, _ := strconv.ParseInt(n.Value, 0, 0)
			return &Int{Node: n, Value: int(i), Name: "Num"}
		case "HEX":
			i, err := strconv.ParseInt(n.Value, 0, 0)
			if err != nil {
				debug("atomNode", "HEX", err)
			}
			return &Int{Node: n, Value: int(i), Name: "Num"}
		case "FLOAT":
			f, _ := strconv.ParseFloat(n.Value, 64)
			return &Float{Node: n, Value: f, Name: "Num"}
		case "IDENT":
			return &ID{Node: n, Value: n.Value, Name: "ID"}
		case "OP":
			return &ID{Node: n, Value: n.Value, Name: "ID"}
		}
		return n
	case []parsec.ParsecNode:
		//debug("atomNode", n[0], n[1])
		return n
	case map[string]interface{}:
		return n
	case string:
		return &Text{Node: n, Value: strings.ReplaceAll(n, `"`, ""), Name: "Text"}
	case *Refer:
		return n
	default:
		debug("atomNode: unknown type", reflect.TypeOf(n).String())
	}
	return nil
}

func referNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//debug("referNode", ns)
	id := nodeToExpr(ns[1]).String()
	return &Refer{Node: ns[1], Value: id, Name: "Refer"}
}

func alistNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//debug("alistNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	list := []Expr{}
	for _, item := range ilist {
		list = append(list, nodeToExpr(item))
	}
	return &Alist{Node: ns[1], Value: list, Name: "Alist"}
}

func mlistNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//debug("mlistNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	list := []Expr{}
	for _, item := range ilist {
		list = append(list, nodeToExpr(item))
	}
	return &Mlist{Node: ns[1], Value: list, Name: "Mlist"}
}

func llistNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//debug("llistNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	list := []Expr{}
	for _, item := range ilist {
		list = append(list, nodeToExpr(item))
	}
	return &Llist{Node: ns[1], Value: list, Name: "Llist"}
}

func propNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//debug("propNode", ns)
	key := nodeToExpr(ns[0]).String()
	item := nodeToExpr(ns[2])
	return &Prop{Node: ns[1], Key: key, Value: item, Name: "Prop"}
}

func dictNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	//debug("dictNode", ns)
	ilist, ok := ns[1].([]parsec.ParsecNode)
	if !ok {
		return nil
	}
	m := map[string]Expr{}
	for _, item := range ilist {
		prop := nodeToExpr(item).(*Prop)
		m[prop.Key] = nodeToExpr(prop.Value)
	}
	return &Dict{Node: ns[1], Value: m, Name: "Dict"}
}

func commentNode(ns []parsec.ParsecNode) parsec.ParsecNode {
	if ns == nil || len(ns) < 1 {
		return nil
	}
	comment := ns[1].(*parsec.Terminal)
	//debug("commentNode", comment.Value)
	return &Comment{Node: ns[1], Value: comment.Value, Name: "Comment"}
}

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
	case *Comment:
		res = node.(*Comment)
	case []parsec.ParsecNode:
		//debug("nodeToExpr: []parsec.ParsecNode:", node)
		nodes := node.([]parsec.ParsecNode)
		if len(nodes) == 1 {
			res = nodeToExpr(nodes[0])
		} else {
			debug("nodeToExpr: []parsec.ParsecNode: len > 1", node)
		}
	case *parsec.Terminal:
		//debug("nodeToExpr: *parsec.Terminal", node)
		n := node.(*parsec.Terminal)
		switch n.Name {
		case "INT":
			i, _ := strconv.Atoi(n.Value)
			res = &Int{Node: n, Value: i, Name: "Num"}
		case "FLOAT":
			f, _ := strconv.ParseFloat(n.Value, 64)
			res = &Float{Node: n, Value: f, Name: "Num"}
		case "IDENT":
			res = &ID{Node: n, Value: n.Value, Name: "ID"}
		case "OP":
			res = &ID{Node: n, Value: n.Value, Name: "ID"}
		}
	default:
		debug("nodeToExpr: unknown type", reflect.TypeOf(node))
	}
	return
}

func debugNodes(nodes []parsec.ParsecNode, deep int) {

	for _, node := range nodes {
		switch node.(type) {
		case *parsec.Terminal:
			v := node.(*parsec.Terminal)
			debug(strings.Repeat(".", deep*2), v.Name, v.Value, v.Position, v.Attributes)
		case *Int:
			v := node.(*Int)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *Float:
			v := node.(*Float)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *ID:
			v := node.(*ID)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *Refer:
			v := node.(*Refer)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *Alist:
			v := node.(*Alist)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *Mlist:
			v := node.(*Mlist)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *Llist:
			v := node.(*Llist)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *Dict:
			v := node.(*Dict)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case *Text:
			v := node.(*Text)
			debug(strings.Repeat(".", deep*2), v.Debug())
		case []parsec.ParsecNode:
			v := node.([]parsec.ParsecNode)
			debugNodes(v, deep+1)
		default:
			debug(reflect.TypeOf(node).String())
		}
	}
}
