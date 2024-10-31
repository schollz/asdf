package main

import (
	"flag"
	"fmt"

	"github.com/schollz/asdf/src/runner"
)

var flagFilename string
var flagVersion bool
var Version string

func init() {
	flag.BoolVar(&flagVersion, "version", false, "print version")
	flag.StringVar(&flagFilename, "filename", "", "filename to run")
}

func main() {
	flag.Parse()
	if flagVersion {
		fmt.Println("asdf version", Version)
		return
	}
	runner.Run(flagFilename)
}
