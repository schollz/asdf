package fileparser

import (
	"fmt"
	"testing"
)

func TestFileparser(t *testing.T) {
	filename := "test.js"
	sequences, err := Parse(filename)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	for _, b := range sequences.Blocks {
		fmt.Printf("block %s\n", b.Name)
		for _, s := range b.Steps {
			fmt.Printf("step %+v\n", s.Info())
		}
	}

	fmt.Println("outputs")
	for _, o := range sequences.Sprockets {
		fmt.Printf("output %s\n", o.Name)
		for _, s := range o.Block.Steps {
			fmt.Printf("step %+v\n", s.Info())
		}
	}

}
