package emitter

import (
	"testing"

	log "github.com/schollz/logger"
)

func TestMidi(t *testing.T) {
	log.SetLevel("trace")
	m, err := NewMidi("through")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	log.Tracef("%+v", m)
}
