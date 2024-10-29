package emitter

import (
	log "github.com/schollz/logger"
)

type Emitter interface {
	NoteOn(note int, velocity int)
	NoteOff(note int)
	Set(param string, value int)
}

type DebugEmitter struct{}

func (d *DebugEmitter) NoteOn(note int, velocity int) {
	log.Debugf("NoteOn: %d %d", note, velocity)
}

func (d *DebugEmitter) NoteOff(note int) {
	log.Debugf("NoteOff: %d", note)
}

func (d *DebugEmitter) Set(param string, value int) {
	log.Debugf("Set: %s %d", param, value)
}
