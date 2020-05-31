package jpl

import "fmt"

var lastFork *Fork

//
const (
	OK   = "ok"
	FAIL = "fail"
	NEXT = "next"
)

var (
	okID   = &ID{Value: OK, Name: "ID"}
	failID = &ID{Value: FAIL, Name: "ID"}
	nextID = &ID{Value: NEXT, Name: "ID"}
)

func backtrFuncs() map[string]Func {
	return map[string]Func{
		"among": among,
		"setu":  setu,
	}
}

func setu(args []Expr, ctxName string) Expr {
	if len(args) < 2 {
		return errID
	}
	id, ok := args[0].Eval().(*ID)
	if !ok {
		return errID
	}
	e := args[1].Eval()
	var old = engine.current[ctxName].set(id.Value, e)
	lastFork.addUndo(&Llist{Name: "Llist", Value: []Expr{&ID{Name: "ID", Value: "set", CtxName: ctxName}, id, old}, CtxName: ctxName})
	return old
}

func among(args []Expr, ctxName string) Expr {
	if len(args) < 4 {
		return errID
	}
	var res Expr = nullID
	dict, ok := args[0].Eval().(*Dict)
	if !ok {
		return errID
	}
	engine.current[ctxName].push(dict.Value, ctxName)
	v, ok := args[1].Eval().(*ID)
	if !ok {
		return errID
	}
	alist, ok := args[2].Eval().(*Alist)
	if !ok {
		return errID
	}
	e := args[3]
	res = runAmong(v.Value, alist.Value, e, ctxName)
	engine.current[ctxName].pop(ctxName)
	return res
}

func runAmong(v string, list []Expr, e Expr, ctxName string) Expr {
	var res Expr = failID
	deep := 0
	if lastFork != nil {
		deep = lastFork.deep + 1
	}
	forkID := &ID{Name: "ID", Value: fmt.Sprintf("among-%d", deep), CtxName: ctxName}
	lastFork = makeFork(forkID, deep, list)

	for _, item := range list {
		val := item.Eval()
		engine.current[ctxName].set(v, val)
		res = e.Eval()
		id, ok := res.(*ID)
		if ok && id.Equals(failID) {
			lastFork.runUndo()
			continue
		}
		if ok && id.Equals(okID) {
			lastFork.up(&ID{Name: "ID", Value: fmt.Sprintf("among-%d", deep), CtxName: ctxName})
			return res
		}
		if ok && id.Equals(nextID) {
			lastFork.runUndo()
			continue
		}
		res = runAmong(v, list, e, ctxName)
		resid, ok := res.(*ID)
		if ok && resid == okID {
			lastFork.up(&ID{Name: "ID", Value: fmt.Sprintf("among-%d", deep), CtxName: ctxName})
			return res
		}
		lastFork.runUndo()

	}
	lastFork.up(&ID{Name: "ID", Value: fmt.Sprintf("among-%d", deep), CtxName: ctxName})
	return res
}

func makeFork(id *ID, deep int, alt []Expr) *Fork {
	return &Fork{
		parent: lastFork,
		id:     id,
		deep:   deep,
		alt:    alt,
	}
}

// Fork -
type Fork struct {
	parent *Fork
	alt    []Expr
	undo   []*Llist
	id     *ID
	deep   int
}

func (fork *Fork) runUndo() {
	for _, u := range fork.undo {
		u.Eval()
	}
	fork.undo = []*Llist{}
}

func (fork *Fork) addUndo(l *Llist) {
	//engine.debug("fork.addUndo", l.Debug())
	fork.undo = append(fork.undo, l)
}

func (fork *Fork) up(id *ID) {
	f := lastFork
	for !f.id.Equals(id) {
		f = f.parent
	}
	lastFork = f.parent
}
