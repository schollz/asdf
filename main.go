package main

import (
	"flag"
	"fmt"

	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/runner"
)

var flagFilename string
var flagVersion bool
var flagMidiOuts bool
var Version string

func init() {
	flag.BoolVar(&flagVersion, "version", false, "print version")
	flag.BoolVar(&flagMidiOuts, "midi", false, "lists available midi outputs")
	flag.StringVar(&flagFilename, "filename", "", "filename to run")
}

func main() {
	flag.Parse()
	if flagVersion {
		fmt.Println("asdf version", Version)
		return
	}
	if flagMidiOuts {
		outs, err := emitter.ListMidiOuts()
		if err != nil {
			fmt.Println(err)
		} else {
			for i, v := range outs {
				fmt.Printf("%d) '%s'", i, v)
			}
		}
		return
	}
	runner.Run(flagFilename)
}
