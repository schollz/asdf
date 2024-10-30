package emitter

import (
	"testing"
	"time"

	log "github.com/schollz/logger"
)

func TestMidi(t *testing.T) {
	log.SetLevel("trace")
	m, err := NewMidi("through")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	log.Tracef("%+v", m)
	m.NoteOn(60, 100)
	time.Sleep(1 * time.Second)
	m.NoteOff(60)
}
