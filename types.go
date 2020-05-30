package jpl

import (
	"fmt"
	"sort"

	parsec "github.com/prataprc/goparsec"
)

// Func -
type Func func([]Expr) Expr

// Expr -
type Expr interface {
	fmt.Stringer
	Eval() Expr
	Equals(Expr) bool
	Clone() Expr
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

// Equals -
func (num *Int) Equals(e Expr) (res bool) {
	v, ok := e.(*Int)
	if !ok {
		v, ok := e.(*Float)
		return ok && num.Value == int(v.Value)
	}
	res = ok && num.Value == v.Value
	return
}

// Clone -
func (num *Int) Clone() (res Expr) {
	return &Int{Value: num.Value, Name: num.Name, Node: num.Node}
}

// Float -
type Float struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value float64
}

func (num *Float) String() (res string) {
	return fmt.Sprintf("%.4f", num.Value)
}

// Debug -
func (num *Float) Debug() (res string) {
	return fmt.Sprintf("%s:%f", num.Name, num.Value)
}

// Eval -
func (num *Float) Eval() (res Expr) {
	return num
}

// Equals -
func (num *Float) Equals(e Expr) (res bool) {
	v, ok := e.(*Float)
	if !ok {
		v, ok := e.(*Int)
		return ok && num.Value == float64(v.Value)
	}
	res = ok && num.Value == v.Value
	return
}

// Clone -
func (num *Float) Clone() (res Expr) {
	return &Float{Value: num.Value, Name: num.Name, Node: num.Node}
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

// Equals -
func (id *ID) Equals(e Expr) (res bool) {
	v, ok := e.(*ID)
	res = ok && id.Value == v.Value
	return
}

// Clone -
func (id *ID) Clone() (res Expr) {
	return &ID{Value: id.Value, Name: id.Name, Node: id.Node}
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

// Equals -
func (ref *Refer) Equals(e Expr) (res bool) {
	v, ok := e.(*Refer)
	res = ok && ref.Value == v.Value
	return
}

// Clone -
func (ref *Refer) Clone() (res Expr) {
	return &Refer{Value: ref.Value, Name: ref.Name, Node: ref.Node}
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

// Equals -
func (alist *Alist) Equals(e Expr) (res bool) {
	v, ok := e.(*Alist)
	if !ok {
		return false
	}
	res = true
	for i, item := range alist.Value {
		if !item.Equals(v.Value[i]) {
			return false
		}
	}
	return
}

// Clone -
func (alist *Alist) Clone() (res Expr) {
	return &Alist{Value: alist.Value, Name: alist.Name, Node: alist.Node}
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
		return nullID
	}
	return applyFunc(llist.Value[0].Eval(), llist.Value[1:])
}

// Equals -
func (llist *Llist) Equals(e Expr) (res bool) {
	v, ok := e.(*Llist)
	if !ok {
		return false
	}
	res = true
	for i, item := range llist.Value {
		if !item.Equals(v.Value[i]) {
			return false
		}
	}
	return
}

// Clone -
func (llist *Llist) Clone() (res Expr) {
	return &Llist{Value: llist.Value, Name: llist.Name, Node: llist.Node}
}

// Mlist -
type Mlist struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value []Expr
}

func (mlist *Mlist) String() (res string) {
	res = "<"
	sep := ""
	for _, item := range mlist.Value {
		res += fmt.Sprintf("%s%v", sep, item)
		sep = " "
	}
	res += ">"
	return
}

// Debug -
func (mlist *Mlist) Debug() (res string) {
	return fmt.Sprintf("%s:%s", mlist.Name, mlist.String())
}

// Eval -
func (mlist *Mlist) Eval() (res Expr) {
	return mlist
}

// Equals -
func (mlist *Mlist) Equals(e Expr) (res bool) {
	v, ok := e.(*Mlist)
	if !ok {
		return false
	}
	res = true
	for i, item := range mlist.Value {
		if !item.Equals(v.Value[i]) {
			return false
		}
	}
	return
}

// Clone -
func (mlist *Mlist) Clone() (res Expr) {
	return &Llist{Value: mlist.Value, Name: mlist.Name, Node: mlist.Node}
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
	keys := []string{}
	for key := range dict.Value {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		res += fmt.Sprintf("%s%s:%v", sep, key, dict.Value[key])
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

// Equals -
func (dict *Dict) Equals(e Expr) (res bool) {
	v, ok := e.(*Dict)
	if !ok {
		return false
	}
	res = true
	for key, item := range dict.Value {
		p, ok := v.Value[key]
		if !(ok && item.Equals(p)) {
			return false
		}
	}
	return
}

// Clone -
func (dict *Dict) Clone() (res Expr) {
	return &Dict{Value: dict.Value, Name: dict.Name, Node: dict.Node}
}

// Text -
type Text struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value string
}

func (text *Text) String() (res string) {
	return fmt.Sprintf(`"%s"`, text.Value)
}

// Debug -
func (text *Text) Debug() (res string) {
	return fmt.Sprintf(`%s:"%s"`, text.Name, text.Value)
}

// Eval -
func (text *Text) Eval() (res Expr) {
	return text
}

// Equals -
func (text *Text) Equals(e Expr) (res bool) {
	v, ok := e.(*Text)
	res = ok && text.Value == v.Value
	return
}

// Clone -
func (text *Text) Clone() (res Expr) {
	return &Text{Value: text.Value, Name: text.Name, Node: text.Node}
}

// Lamb -
type Lamb struct {
	Expr
	Name   string
	Params []*ID
	Body   Expr
}

func (lamb *Lamb) String() (res string) {
	return fmt.Sprintf("%v=>%v", lamb.Params, lamb.Body)
}

// Debug -
func (lamb *Lamb) Debug() (res string) {
	return fmt.Sprintf("%s:%s", lamb.Name, lamb.String())
}

// Eval -
func (lamb *Lamb) Eval() (res Expr) {
	return lamb
}

// Equals -
func (lamb *Lamb) Equals(e Expr) (res bool) {
	v, ok := e.(*Lamb)
	res = ok && lamb == v
	return
}

// Clone -
func (lamb *Lamb) Clone() (res Expr) {
	return &Lamb{Params: lamb.Params, Body: lamb.Body, Name: lamb.Name}
}

//var lambLock sync.RWMutex

// Apply -
func (lamb *Lamb) Apply(args []Expr) (res Expr) {
	//debug(lamb.Debug(), args)
	if len(lamb.Params) != len(args) {
		return errID
	}
	vars := map[string]Expr{}
	for i, item := range lamb.Params {
		vars[item.Value] = args[i].Eval()
	}
	//debug("Lamb Apply", "locking...", lamb.Debug(), args)
	//lambLock.Lock()
	//debug("Lamb Apply", "locked", lamb.Debug(), args)
	current.push(vars)
	res = lamb.Body.Eval()
	current.pop()
	//lambLock.Unlock()
	//debug("Lamb Apply", "Unlocked", lamb.Debug(), args)
	return
}

// Comment -
type Comment struct {
	Expr
	Node  parsec.ParsecNode
	Name  string
	Value string
}

func (com *Comment) String() (res string) {
	return fmt.Sprintf("%s", com.Value)
}

// Debug -
func (com *Comment) Debug() (res string) {
	return fmt.Sprintf("%s:%s", com.Name, com.Value)
}

// Eval -
func (com *Comment) Eval() (res Expr) {
	return com
}

// Equals -
func (com *Comment) Equals(e Expr) (res bool) {
	v, ok := e.(*Comment)
	res = ok && com.Value == v.Value
	return
}

// Clone -
func (com *Comment) Clone() (res Expr) {
	return &Comment{Value: com.Value, Name: com.Name, Node: com.Node}
}
