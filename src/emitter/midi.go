package emitter

import (
	"fmt"
	"strings"

	log "github.com/schollz/logger"
	midi "gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type Midi struct {
	Name string
	Conn drivers.Out
}

func NewMidi(name string) (m Midi, err error) {
	m.Name = name

	outs := midi.GetOutPorts()
	if len(outs) == 0 {
		err = fmt.Errorf("no MIDI output ports available")
		return
	}

	// find the matching output
	for _, out := range outs {
		log.Tracef("found MIDI output: %s", out.String())
		if strings.Contains(strings.ToLower(out.String()), strings.ToLower(name)) {
			m.Conn = out
		}
	}

	if m.Conn == nil {
		err = fmt.Errorf("no MIDI output port found with name %s", name)
	}

	// OPEN PORT
	err = m.Conn.Open()

	return
}

func (m *Midi) NoteOn(note int, velocity int) {
	if m.Conn == nil {
		return
	}
	// Send Note On message with the specified note and velocity
	err := m.Conn.Send(midi.NoteOn(0, uint8(note), uint8(velocity)))
	if err != nil {
		log.Errorf("Failed to send Note On: %v", err)
	}
}

func (m *Midi) NoteOff(note int) {
	if m.Conn == nil {
		return
	}
	// Send Note Off message with the specified note
	err := m.Conn.Send(midi.NoteOff(0, uint8(note)))
	if err != nil {
		log.Errorf("Failed to send Note Off: %v", err)
	}
}
