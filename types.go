package jpl

import (
	"fmt"

	parsec "github.com/prataprc/goparsec"
)

// Func -
type Func func([]Expr) Expr

// Expr -
type Expr interface {
	fmt.Stringer
	Eval() Expr
	Debug() string
}

// Int -
type Int struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value int
}

func (num *Int) String() (res string) {
	return fmt.Sprintf("%d", num.Value)
}

// Debug -
func (num *Int) Debug() (res string) {
	return fmt.Sprintf("%s:%d", num.Name, num.Value)
}

// Eval -
func (num *Int) Eval() (res Expr) {
	return num
}

// Float -
type Float struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value float64
}

func (num *Float) String() (res string) {
	return fmt.Sprintf("%f", num.Value)
}

// Debug -
func (num *Float) Debug() (res string) {
	return fmt.Sprintf("%s:%f", num.Name, num.Value)
}

// Eval -
func (num *Float) Eval() (res Expr) {
	return num
}

// ID -
type ID struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value string
}

func (id *ID) String() (res string) {
	return fmt.Sprintf("%s", id.Value)
}

// Debug -
func (id *ID) Debug() (res string) {
	return fmt.Sprintf("%s:%s", id.Name, id.Value)
}

// Eval -
func (id *ID) Eval() (res Expr) {
	return id
}

// Refer -
type Refer struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value string
}

func (ref *Refer) String() (res string) {
	return fmt.Sprintf(".%s", ref.Value)
}

// Debug -
func (ref *Refer) Debug() (res string) {
	return fmt.Sprintf("%s:.%s", ref.Name, ref.Value)
}

// Eval -
func (ref *Refer) Eval() (res Expr) {
	return current.get(ref.Value)
}

// Alist -
type Alist struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value []Expr
}

func (alist *Alist) String() (res string) {
	res = "["
	sep := ""
	for _, item := range alist.Value {
		res += fmt.Sprintf("%s%v", sep, item)
		sep = " "
	}
	res += "]"
	return
}

// Debug -
func (alist *Alist) Debug() (res string) {
	return fmt.Sprintf("%s:%s", alist.Name, alist.String())
}

// Eval -
func (alist *Alist) Eval() (res Expr) {
	a := []Expr{}
	for _, item := range alist.Value {
		a = append(a, item.Eval())
	}
	return &Alist{Node: alist.Node, Name: alist.Name, Value: a}
}

// Llist -
type Llist struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value []Expr
}

func (llist *Llist) String() (res string) {
	res = "("
	sep := ""
	for _, item := range llist.Value {
		res += fmt.Sprintf("%s%v", sep, item)
		sep = " "
	}
	res += ")"
	return
}

// Debug -
func (llist *Llist) Debug() (res string) {
	return fmt.Sprintf("%s:%s", llist.Name, llist.String())
}

// Eval -
func (llist *Llist) Eval() (res Expr) {
	if len(llist.Value) == 0 {
		return &ID{Value: "null", Name: "ID"}
	}
	name := llist.Value[0].Eval().(*ID).Value
	fn, ok := funcs[name]
	if !ok {
		return &ID{Value: "undefined", Name: "ID"}
	}
	return fn(llist.Value[1:])
}

// Prop -
type Prop struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Key   string
	Value Expr
}

func (prop *Prop) String() (res string) {
	return fmt.Sprintf("%s:%v", prop.Key, prop.Value)
}

// Debug -
func (prop *Prop) Debug() (res string) {
	return fmt.Sprintf("%s:%v", prop.Name, prop.String())
}

// Dict -
type Dict struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value map[string]Expr
}

func (dict *Dict) String() (res string) {
	res = "{"
	sep := ""
	for key, item := range dict.Value {
		res += fmt.Sprintf("%s%s:%v", sep, key, item)
		sep = " "
	}
	res += "}"
	return
}

// Debug -
func (dict *Dict) Debug() (res string) {
	return fmt.Sprintf("%s:%s", dict.Name, dict.String())
}

// Eval -
func (dict *Dict) Eval() (res Expr) {
	d := map[string]Expr{}
	for key, item := range dict.Value {
		d[key] = item.Eval()
	}
	return &Dict{Node: dict.Node, Name: dict.Name, Value: d}
}

// Text -
type Text struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value string
}

func (text *Text) String() (res string) {
	return fmt.Sprintf("%s", text.Value)
}

// Debug -
func (text *Text) Debug() (res string) {
	return fmt.Sprintf("%s:%s", text.Name, text.Value)
}

// Eval -
func (text *Text) Eval() (res Expr) {
	return text
}
