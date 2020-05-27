package jpl

import (
	"flag"
	"os"
	"os/exec"
	"strings"
)

var osFuncs = map[string]Func{
	"host":     host,
	"env":      env,
	"setenv":   setenv,
	"unsetenv": unsetenv,
	"args":     flagArgs,
	"pid":      execPid,
	"pwd":      pwd,
	"cd":       chdir,
	"cmd":      runCmd,
}

func host(args []Expr) Expr {
	if len(args) != 0 {
		return errID
	}

	name, err := os.Hostname()
	if err != nil {
		return errID
	}
	return &Text{Name: "Text", Value: name}
}

func env(args []Expr) Expr {
	if len(args) > 1 {
		return errID
	}
	if len(args) == 0 {
		envs := os.Environ()
		list := make([]Expr, len(envs))
		for i, item := range envs {
			var key, val string
			//debug("env", item)
			e := strings.Split(item, "=")
			switch len(e) {
			case 0:
				continue
			case 1:
				key = e[0]
				val = ""
			case 2:
				key = e[0]
				val = e[1]
			default:
				key = e[0]
				val = strings.Join(e[1:], "=")
			}
			list[i] = &Alist{Name: "Alist", Value: []Expr{&Text{Name: "Text", Value: key}, &Text{Name: "Text", Value: val}}}
		}
		return &Alist{Name: "Alist", Value: list}
	}
	name, ok := args[0].(*Text)
	if !ok {
		return errID
	}
	return &Text{Name: "Text", Value: os.Getenv(name.Value)}
}

func flagArgs(args []Expr) Expr {
	if len(args) > 0 {
		return errID
	}
	a := flag.Args()
	list := make([]Expr, len(a)-1)
	for i, arg := range a[1:] {
		list[i] = parse([]Expr{&Text{Name: "Text", Value: arg}})
	}
	return &Alist{Name: "Alist", Value: list}
}

func execArgs(args []Expr) Expr {
	if len(args) > 0 {
		return errID
	}
	list := make([]Expr, len(os.Args)-2)
	for i, arg := range os.Args[2:] {
		list[i] = &Text{Name: "Text", Value: arg}
	}
	return &Alist{Name: "Alist", Value: list}
}

func execPid(args []Expr) Expr {
	if len(args) > 0 {
		return errID
	}
	return &Int{Name: "Num", Value: os.Getpid()}
}

func pwd(args []Expr) Expr {
	if len(args) > 0 {
		return errID
	}
	dir, err := os.Getwd()
	if err != nil {
		return errID
	}
	return &Text{Name: "Text", Value: dir}
}

func chdir(args []Expr) Expr {
	if len(args) != 1 {
		return errID
	}
	dir, ok := args[0].Eval().(*Text)
	if !ok {
		return errID
	}
	wd, err := os.Getwd()
	if err != nil {
		return errID
	}
	err = os.Chdir(dir.Value)
	if err != nil {
		return errID
	}
	return &Text{Name: "Text", Value: wd}
}

func setenv(args []Expr) Expr {
	if len(args) != 2 {
		return errID
	}
	name, ok := args[0].Eval().(*Text)
	if !ok {
		return errID
	}
	val, ok := args[1].Eval().(*Text)
	if !ok {
		return errID
	}
	old := os.Getenv(name.Value)
	err := os.Setenv(name.Value, val.Value)
	if err != nil {
		return errID
	}
	return &Text{Name: "Text", Value: old}
}

func unsetenv(args []Expr) Expr {
	if len(args) != 1 {
		return errID
	}
	name, ok := args[0].Eval().(*Text)
	if !ok {
		return errID
	}
	old := os.Getenv(name.Value)
	err := os.Unsetenv(name.Value)
	if err != nil {
		return errID
	}
	return &Text{Name: "Text", Value: old}
}

func runCmd(args []Expr) Expr {
	name, ok := args[0].Eval().(*Text)
	if !ok {
		return errID
	}
	params := make([]string, len(args)-1)
	for i := 1; i < len(args); i++ {
		a, ok := args[i].Eval().(*Text)
		if !ok {
			return errID
		}
		params[i-1] = a.Value
	}
	cmd := exec.Command(name.Value, params...)
	//debug(cmd)
	res, err := cmd.CombinedOutput()
	var errMess Expr
	if err != nil {
		errMess = &Text{Name: "Text", Value: err.Error()}
	} else {
		errMess = nullID
	}
	var lines []Expr
	if err == nil {
		outLines := strings.Split(string(res), "\n")
		lines = make([]Expr, len(outLines))
		for i, line := range outLines {
			lines[i] = &Text{Name: "Text", Value: line}
		}
	}
	return &Dict{Name: "Dict", Value: map[string]Expr{
		"error": errMess,
		"out":   &Alist{Name: "Alist", Value: lines},
	}}
}
