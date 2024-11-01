package emitter

import (
	log "github.com/schollz/logger"
)

type Emitter interface {
	NoteOn(note int, velocity int)
	NoteOff(note int)
	Set(param string, value float64)
}

type Debugger struct{}

func (d *Debugger) NoteOn(note int, velocity int) {
	log.Debugf("NoteOn: %d %d", note, velocity)
}

func (d *Debugger) NoteOff(note int) {
	log.Debugf("NoteOff: %d", note)
}

func (d *Debugger) Set(param string, value float64) {
	log.Debugf("Set: %s %2.3f", param, value)
}
