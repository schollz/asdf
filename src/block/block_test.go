package block

import (
	"fmt"
	"testing"

	log "github.com/schollz/logger"
)

func TestExpand(t *testing.T) {
	log.SetLevel("trace")
	block := `.bpm120
c4 d4
- Em.arp.u4.gate50
f g (a b c)*2`
	block_expected := `.bpm120
c4 d4
- - - - e4.gate50 g4.gate50 b4.gate50 e5.gate50
f - - - - - g - - - - - a b c a b c
`
	result, err := expand(block)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != block_expected {
		t.Errorf("\n%s -->\n%v != %v", block, result, block_expected)
	}
}

func TestParse(t *testing.T) {
	log.SetLevel("trace")
	blockString := `.bpm120
c4 ~ b3 c4
- (Em f)*2
Em7.arp.u4.gate50 ~
`
	block, err := Parse(blockString)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(block.Steps) != 11 {
		t.Errorf("expected 11 steps, got %d", len(block.Steps))
	}
	if block.TotalTime != 6.0 {
		t.Errorf("expected 6.0 total time, got %f", block.TotalTime)
	}
}

func TestParse2(t *testing.T) {
	log.SetLevel("trace")
	blockString := `
c4.bpm60 ~ b3 c4
c4.bpm30 ~ b3 c4
`
	block, err := Parse(blockString)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(block.Steps) != 6 {
		t.Errorf("expected 6 steps, got %d", len(block.Steps))
	}
	if block.TotalTime != 12.0 {
		t.Errorf("expected 12.0 total time, got %f", block.TotalTime)
	}
}

func TestParse3(t *testing.T) {
	log.SetLevel("trace")
	blockString := `
c4.bpm60.beats3 b3 c4
c4.bpm30 b4 c4
`
	block, err := Parse(blockString)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(block.Steps) != 6 {
		t.Errorf("expected 6 steps, got %d", len(block.Steps))
	}
	if block.TotalTime != 9.0 {
		t.Errorf("expected 12.0 total time, got %f", block.TotalTime)
	}
}

func TestMerge(t *testing.T) {
	log.SetLevel("trace")
	block1, err := Parse("a3 b c d")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	block2, err := Parse("e f g a")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	block1.Add(block2)
	fmt.Println("merged block")
	for _, s := range block1.Steps {
		fmt.Printf("step: %+v\n", s.Info())
	}
}
