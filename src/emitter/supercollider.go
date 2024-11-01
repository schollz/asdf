package emitter

import "github.com/hypebeast/go-osc/osc"

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

func (s *SuperCollider) Set(param string, value float64) {
	msg := osc.NewMessage("/asdf/set")
	msg.Append(s.SynthDef)
	msg.Append(param)
	msg.Append(value)
	s.Client.Send(msg)
}
