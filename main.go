package main

import (
	"flag"
	"fmt"
	"strings"

	cmd "github.com/schollz/asdf/micro/cmd/micro"
	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/runner"
)

var flagFilename string
var flagVersion bool
var flagMidiOuts bool

func init() {
	flag.BoolVar(&flagMidiOuts, "midi", false, "lists available midi outputs")
	flag.StringVar(&flagFilename, "filename", "", "filename to run")
}

func main() {
	flag.Parse()
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
	if flagFilename != "" {
		runner.Run(flagFilename)
		return
	} else {
		cmd.Run()
	}

}
