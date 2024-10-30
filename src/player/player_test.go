package player

import (
	"testing"

	"github.com/schollz/asdf/src/emitter"
	log "github.com/schollz/logger"
)

func TestPlayer(t *testing.T) {
	log.SetLevel("trace")
	p := New([]emitter.Emitter{&emitter.Debugger{}})
	p.NoteOn(60, 100)
	p.NoteOn(72, 100)
	p.NoteOff(60)
	p.Reset()
}
