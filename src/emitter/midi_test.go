package emitter

import (
	"testing"
	"time"

	log "github.com/schollz/logger"
)

func TestMidi(t *testing.T) {
	log.SetLevel("trace")
	outs, _ := ListMidiOuts()
	if len(outs) < 0 {
		t.Errorf("weird error")
	}
	log.Debugf("midi outs available:")
	for i, v := range outs {
		log.Debugf("%d) '%s'", i+1, v)
	}
	m, err := NewMidi("through", 0)
	if err != nil {
		return
	}
	log.Tracef("%+v", m)
	m.NoteOn(60, 100)
	time.Sleep(1 * time.Second)
	m.NoteOff(60)
}
