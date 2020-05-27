package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/youryharchenko/jpl"
)

func main() {
	verb := flag.Bool("v", false, "Verbose")
	flag.Parse()

	jpl.Debug = *verb
	if len(flag.Args()) < 1 {
		log.Println("try with parameters")
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
