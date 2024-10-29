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
		{"Cmaj.arpeggio16.prob30", "Cmaj.probability30"},
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
