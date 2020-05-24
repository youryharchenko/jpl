package jpl

// Match -
type Match func([]Expr, Expr) bool

var matches = map[string]Match{
	"atom":  patAtom,
	"id":    patID,
	"num":   patNum,
	"int":   patInt,
	"float": patFloat,
	"text":  patText,
	"refer": patRefer,
}

func patAtom(args []Expr, e Expr) bool {
	if patNum(args, e) || patID(args, e) || patText(args, e) || patRefer(args, e) {
		return true
	}
	return false
}

func patID(args []Expr, e Expr) bool {
	if _, ok := e.(*ID); ok {
		return true
	}
	return false
}

func patInt(args []Expr, e Expr) bool {
	if _, ok := e.(*Int); ok {
		return true
	}
	return false
}

func patFloat(args []Expr, e Expr) bool {
	if _, ok := e.(*Float); ok {
		return true
	}
	return false
}

func patNum(args []Expr, e Expr) bool {
	if patInt(args, e) || patFloat(args, e) {
		return true
	}
	return false
}

func patText(args []Expr, e Expr) bool {
	if _, ok := e.(*Text); ok {
		return true
	}
	return false
}

func patRefer(args []Expr, e Expr) bool {
	if _, ok := e.(*Refer); ok {
		return true
	}
	return false
}

func match(pat Expr, e Expr) Expr {
	patCtx := &Pattern{}
	patCtx.begin()
	if patCtx.matchExpr(pat, e) {
		patCtx.commit()
		return trueID
	}
	patCtx.rollback()
	return falseID
}

// Pattern -
type Pattern struct {
	clon *Context
}

func (pat *Pattern) matchExpr(p Expr, e Expr) (res bool) {
	switch pt := p.(type) {
	case *ID:
		res = pat.matchID(pt, e)
	case *Int:
		res = pat.matchInt(pt, e)
	case *Float:
		res = pat.matchFloat(pt, e)
	case *Llist:
		res = pat.matchLlist(pt, e)
	case *Refer:
		res = pat.matchRefer(pt, e)
	case *Alist:
		res = pat.matchAlist(pt, e)
	case *Text:
		res = pat.matchText(pt, e)
	case *Dict:
		res = pat.matchDict(pt, e)
	case *Mlist:
		res = pat.matchMlist(pt, e)
	}
	return
}

func (pat *Pattern) matchID(p *ID, e Expr) (res bool) {
	return p.Equals(e)
}

func (pat *Pattern) matchInt(p *Int, e Expr) (res bool) {
	return p.Equals(e)
}

func (pat *Pattern) matchFloat(p *Float, e Expr) (res bool) {
	return p.Equals(e)
}

func (pat *Pattern) matchRefer(p *Refer, e Expr) (res bool) {
	if current.bound(p.Value) {
		if current.get(p.Value).Equals(nullID) {
			current.set(p.Value, e)
			return true
		}
		return p.Eval().Equals(e)
	}
	return false
}

func (pat *Pattern) matchAlist(p *Alist, e Expr) (res bool) {
	ealist, ok := e.(*Alist)
	if !ok {
		return false
	}
	if len(p.Value) != len(ealist.Value) {
		return false
	}
	for i, item := range p.Value {
		if !pat.matchExpr(item, ealist.Value[i]) {
			return false
		}
	}
	return true
}

func (pat *Pattern) matchDict(p *Dict, e Expr) (res bool) {
	edict, ok := e.(*Dict)
	if !ok {
		return false
	}
	for key, item := range p.Value {
		v, ok := edict.Value[key]
		if !ok {
			return false
		}
		if !pat.matchExpr(item, v) {
			return false
		}
	}
	return true
}

func (pat *Pattern) matchLlist(p *Llist, e Expr) (res bool) {
	return p.Equals(e)
}

func (pat *Pattern) matchText(p *Text, e Expr) (res bool) {
	return p.Equals(e)
}

func (pat *Pattern) matchMlist(p *Mlist, e Expr) (res bool) {
	nameID, ok := p.Value[0].Eval().(*ID)
	if !ok {
		return false
	}
	name := nameID.Value
	f, ok := matches[name]
	if !ok {
		return false
	}
	return f(p.Value[1:], e)
}

func (pat *Pattern) begin() {
	pat.clon = current.clone()
	cl := pat.clon
	for cl.parent != nil {
		cl.parent = cl.parent.clone()
		cl = cl.parent
	}
}

func (pat *Pattern) commit() {
	pat.clon = nil
}

func (pat *Pattern) rollback() {
	current = pat.clon
}
