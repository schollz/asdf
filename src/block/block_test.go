package block

import (
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
	block := `.bpm120
c4 ~ b3 d4 
- (e f)*2
`
	steps, err := Parse(block)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(steps) != 7 {
		t.Errorf("expected 7 steps, got %d", len(steps))
	}
}
