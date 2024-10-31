package main

import (
	"flag"

	"github.com/schollz/asdf/src/runner"
)

var flagFilename string

func init() {
	flag.StringVar(&flagFilename, "filename", "", "filename to run")
}

func main() {
	flag.Parse()
	runner.Run(flagFilename)
}
