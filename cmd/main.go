package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/youryharchenko/jpl"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		log.Println("try with one parameter")
		os.Exit(0)
	}
	file := flag.Args()[0]
	src, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}

	nodes := jpl.Parse(src)
	jpl.EvalNodes(nodes)
}
