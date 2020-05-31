package jpl

import (
	"flag"
	"io/ioutil"
	"testing"
)

func TestAll(t *testing.T) {
	file, err := ioutil.ReadFile("test/test.jpl")
	if err != nil {
		t.Error(err)
		return
	}
	flag.Parse()
	eng := New()
	nodes := eng.Parse(file)
	eng.EvalNodes(nodes)
}

/*
func Test1(t *testing.T) {
	s := parsec.NewScanner([]byte("  "))
	v, _ := Y(s)
	t.Error(v)
}

func Test2(t *testing.T) {
	s := parsec.NewScanner([]byte(`1 -2 3.8 [3 5] a b [] () (+ 1 3){} {a:3 b:null c:[0 1 2] d:{} e:(- 3 2) f:"yes"} "" "Hello world" ["yes" "no" "ok"]`))
	v, _ := Y(s)
	nodes := v.([]parsec.ParsecNode)
	debugNodes(nodes, 0)
	t.Error("")
}

func Test3(t *testing.T) {
	s := parsec.NewScanner([]byte(`1 -2 3.8 abc "" "Hello world" [][7 8 a]{} {a:3 b:null c:[0 1 2] d:{}} (+ 1 3) () (x) (+ 2 3 5 7) (* 2 3 5 7)`))
	v, _ := Y(s)
	nodes := v.([]parsec.ParsecNode)
	EvalNodes(nodes)
	t.Error("")
}

func Test4(t *testing.T) {
	s := parsec.NewScanner([]byte(`(+ 2 3.1 5 7) (* 2 3 5.4 7) (+ 1 (+ 1 (+ 1 1))) (* 1 (* 2 (* 3 4)))  (* 1 (* 2 (* 3 4.0)))`))
	v, _ := Y(s)
	nodes := v.([]parsec.ParsecNode)
	EvalNodes(nodes)
	t.Error("")
}

func Test5(t *testing.T) {
	i, err := strconv.ParseInt("0xA", 0, 0)
	if err != nil {
		t.Error(err, err.(*strconv.NumError), err.(*strconv.NumError).Num, err.(*strconv.NumError).Err)
	} else {
		t.Error(i)
	}
	t.Error("")
}

func Test6(t *testing.T) {
	i, err := strconv.ParseUint("0xA", 0, 0)
	if err != nil {
		t.Error(err)
	} else {
		t.Error(i)
	}
	t.Error("")
}
*/
