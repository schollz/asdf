package step

import (
	"github.com/schollz/asdf/src/param"

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

type Note struct {
	Midi   int
	Name   string
	Params []*param.Param
}

type Step struct {
	Notes     []Note
	Emitters  []*Emitter
	TimeStart float64
	TimeEnd   float64
}

func NewNote(midi int, name string, params []*param.Param) Note {
	return Note{
		Midi:   midi,
		Name:   name,
		Params: params,
	}
}

func NewStep(notes []Note, timeStart float64, timeEnd float64, emitters []*Emitter) *Step {
	return &Step{
		Notes:     notes,
		TimeStart: timeStart,
		TimeEnd:   timeEnd,
		Emitters:  emitters,
	}
}
