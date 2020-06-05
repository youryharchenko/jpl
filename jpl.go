package jpl

import (
	"bytes"
	"log"
	"sync"

	parsec "github.com/prataprc/goparsec"
)

var engine *JPL

// JPL -
type JPL struct {
	Y     parsec.Parser
	Debug bool

	funcs   map[string]Func
	matches map[string]Match

	global   *Context
	current  sync.Map //map[string]*Context
	treeLock sync.RWMutex

	actors     map[string]*Actor // = map[string]*Actor{}
	actorsLock sync.RWMutex      //= sync.RWMutex{}
	waitGroup  sync.WaitGroup
	stopCh     chan struct{} //= make(chan struct{})

	anyClasses map[string]AnyClass
}

// New -
func New() (jpl *JPL) {

	jpl = &JPL{
		global:  &Context{parent: nil, vars: map[string]Expr{}},
		current: sync.Map{}, // map[string]*Context{},
	}
	jpl.initParser()
	jpl.initFuncs()
	jpl.initMatches()
	jpl.initActors()
	jpl.anyClasses = anyClasses()
	jpl.current.Store("main", jpl.global) //jpl.current["main"] = jpl.global

	engine = jpl
	return
}

func (jpl *JPL) initParser() {
	var value parsec.Parser
	var values = parsec.Kleene(nil, &value)

	var point = parsec.Atom(".", "POINT")
	var refer = parsec.And(referNode, point, parsec.Ident())
	var oper = parsec.OrdChoice(nil, parsec.Atom("+", "OP"), parsec.Atom("-", "OP"), parsec.Atom("*", "OP"), parsec.Atom("/", "OP"), parsec.Atom("%", "OP"))
	var atom = parsec.OrdChoice(atomNode, parsec.Ident(), parsec.Float(), parsec.Hex(), parsec.Oct(), parsec.Int(), parsec.String(), oper, refer)

	var openSqrt = parsec.Atom("[", "OPENSQRT")
	var closeSqrt = parsec.Atom("]", "CLOSESQRT")
	var alist = parsec.And(alistNode, openSqrt, values, closeSqrt)

	var openPar = parsec.Atom("(", "OPENPAR")
	var closePar = parsec.Atom(")", "CLOSEPAR")
	var llist = parsec.And(llistNode, openPar, values, closePar)

	var openAng = parsec.Atom("<", "OPENANG")
	var closeAng = parsec.Atom(">", "CLOSEANG")
	var mlist = parsec.And(mlistNode, openAng, values, closeAng)

	var colon = parsec.Atom(":", "COLON")
	var property = parsec.And(propNode, parsec.Ident(), colon, &value)
	var properties = parsec.Kleene(nil, property)

	var openBra = parsec.Atom("{", "OPENBRA")
	var closeBra = parsec.Atom("}", "CLOSEBRA")
	var dict = parsec.And(dictNode, openBra, properties, closeBra)

	value = parsec.OrdChoice(nil, atom, alist, llist, mlist, dict)

	jpl.Y = parsec.OrdChoice(nil, values)
}

func (jpl *JPL) debug(args ...interface{}) {
	if jpl.Debug {
		log.Println(args...)
	}
}
func (jpl *JPL) initFuncs(args ...interface{}) {
	jpl.funcs = mergeFuncs(jpl.funcs,
		coreFuncs(), osFuncs(), mathFuncs(), backtrFuncs(), actorFuncs(), jsonFuncs())
}

func (jpl *JPL) initActors(args ...interface{}) {
	jpl.actors = map[string]*Actor{}
	jpl.stopCh = make(chan struct{})
}

// Parse -
func (jpl *JPL) Parse(src []byte) []parsec.ParsecNode {
	s := parsec.NewScanner(jpl.skipComments(src))
	v, _ := jpl.Y(s)
	return v.([]parsec.ParsecNode)
}

// EvalNodes -
func (jpl *JPL) EvalNodes(nodes []parsec.ParsecNode) {
	for _, node := range nodes {
		//engine.debug("evalNodes", node)
		switch node.(type) {
		case []parsec.ParsecNode:
			v := node.([]parsec.ParsecNode)
			jpl.EvalNodes(v)
		default:
			expr := nodeToExpr(node)
			res := expr.Eval()
			engine.debug("expr:", expr, "=>", res)
		}
	}
}

func (jpl *JPL) skipComments(src []byte) []byte {
	buf := bytes.Buffer{}
	skip := false
	for _, i := range src {
		if skip {
			if i == 0xA {
				skip = false
			}
			continue
		} else {
			if i == '#' {
				skip = true
				continue
			}
		}
		buf.WriteByte(i)

	}
	return buf.Bytes()
}
