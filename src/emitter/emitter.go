package emitter

import (
	"github.com/hypebeast/go-osc/osc"
	log "github.com/schollz/logger"
)

type Emitter interface {
	NoteOn(note int, velocity int)
	NoteOff(note int)
	Set(param string, value int)
}

type Debugger struct{}

func (d *Debugger) NoteOn(note int, velocity int) {
	log.Debugf("NoteOn: %d %d", note, velocity)
}

func (d *Debugger) NoteOff(note int) {
	log.Debugf("NoteOff: %d", note)
}

func (d *Debugger) Set(param string, value int) {
	log.Debugf("Set: %s %d", param, value)
}

type SuperCollider struct {
	SynthDef string
	Client   *osc.Client
}

func NewSuperCollider(synthDef string) *SuperCollider {
	return &SuperCollider{
		SynthDef: synthDef,
		Client:   osc.NewClient("localhost", 57120),
	}
}

func (s *SuperCollider) NoteOn(note int, velocity int) {
	msg := osc.NewMessage("/asdf/note_on")
	msg.Append(s.SynthDef)
	msg.Append(note)
	msg.Append(velocity)
	s.Client.Send(msg)
}

func (s *SuperCollider) NoteOff(note int) {
	msg := osc.NewMessage("/asdf/note_off")
	msg.Append(s.SynthDef)
	msg.Append(note)
	s.Client.Send(msg)
}
