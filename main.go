package main

import (
	"flag"
	"fmt"
	"strings"

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
			fmt.Println("midi outs available:")
			for i, v := range outs {
				fmt.Printf("%d) '%s'\n", i+1, strings.Split(v, ":")[0])
			}
		}
		return
	}
	runner.Run(flagFilename)
}
