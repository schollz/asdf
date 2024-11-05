package emitter

import (
	"github.com/hypebeast/go-osc/osc"
	log "github.com/schollz/logger"
)

type SuperCollider struct {
	ID       int
	SynthDef string
	Client   *osc.Client
}

func NewSuperCollider(synthDef string, instance int) SuperCollider {
	return SuperCollider{
		ID:       instance,
		SynthDef: synthDef,
		Client:   osc.NewClient("127.0.0.1", 7771),
	}
}

func (s *SuperCollider) NoteOn(note int, velocity int) {
	log.Tracef("[supercollider%d](%s) note_on %d %d", s.ID, s.SynthDef, note, velocity)
	msg := osc.NewMessage("/asdf")
	msg.Append("note_on")
	msg.Append(int32(s.ID))
	msg.Append(s.SynthDef)
	msg.Append(int32(note))
	msg.Append(int32(velocity))
	err := s.Client.Send(msg)
	if err != nil {
		log.Warnf("error sending message: %s", err)
	}
}

func (s *SuperCollider) NoteOff(note int) {
	msg := osc.NewMessage("/asdf")
	msg.Append("note_off")
	msg.Append(int32(s.ID))
	msg.Append(s.SynthDef)
	msg.Append(int32(note))
	err := s.Client.Send(msg)
	if err != nil {
		log.Warnf("error sending message: %s", err)
	}
}

func (s *SuperCollider) Set(param string, value float64) {
	msg := osc.NewMessage("/asdf")
	msg.Append("set")
	msg.Append(int32(s.ID))
	msg.Append(s.SynthDef)
	msg.Append(param)
	msg.Append(float32(value))
	err := s.Client.Send(msg)
	if err != nil {
		log.Warnf("error sending message: %s", err)
	}
}
