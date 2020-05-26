package jpl

import (
	"log"

	parsec "github.com/prataprc/goparsec"
)

// Debug -
var Debug = true

// Y - root Parser
var Y parsec.Parser

// circular rats
var value parsec.Parser

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

var comment = parsec.And(commentNode, parsec.Atom("#", "HASH"), parsec.TokenExact("[^\n]*", "ALFA"), parsec.TokenExact("\n", "LF"))

var values = parsec.Kleene(nil, &value)

var funcs map[string]Func

func init() {
	funcs = initFuncs()

	value = parsec.OrdChoice(nil, atom, alist, llist, mlist, dict, comment)
	Y = parsec.OrdChoice(nil, values)
}

func debug(args ...interface{}) {
	if Debug {
		log.Println(args...)
	}
}

// Parse -
func Parse(src []byte) []parsec.ParsecNode {
	s := parsec.NewScanner(src)
	v, _ := Y(s)
	return v.([]parsec.ParsecNode)
}
