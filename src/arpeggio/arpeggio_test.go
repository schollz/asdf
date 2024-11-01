package arpeggio

import (
	"testing"

	log "github.com/schollz/logger"
)

func TestExpand(t *testing.T) {
	log.SetLevel("trace")
	tests := []struct {
		line     string
		expected string
	}{
		{"Cmaj.arpeggio.up4.down2.thumb1.up4", "(c4 e4 g4 c5 g4 e4 c4 e4 g4 c5 e5)"},
		{"d4fa.arp.do4.p50", "(d4.probability50 f3.probability50 a3.probability50 d2.probability50)"},
		{"d4fa.p50", "d4fa.p50"},
		{"d4fa.arp.up3.p10,20,30", "(d4.probability10,20,30 f4.probability20,30,10 a4.probability30,10,20)"},
	}
	for _, test := range tests {
		result, err := Expand(test.line)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != test.expected {
			t.Errorf("expected %v but got %v", test.expected, result)
		}
		log.Debugf("%+v", result)
	}
}
