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
	Name    string
	Channel int
	Conn    drivers.Out
}

func ListMidiOuts() (outs []string, err error) {
	outPorts := midi.GetOutPorts()
	if len(outs) == 0 {
		err = fmt.Errorf("no MIDI output ports available")
		return
	}

	// find the matching output
	for _, out := range outPorts {
		outName := strings.ToLower(out.String())
		outName = strings.ReplaceAll(outName, "-", "")
		outs = append(outs, outName)
	}

	return
}

func NewMidi(name string, channel int) (m Midi, err error) {
	m.Name = name
	m.Channel = channel

	outs := midi.GetOutPorts()
	if len(outs) == 0 {
		err = fmt.Errorf("no MIDI output ports available")
		return
	}

	// find the matching output
	found := false
	for _, out := range outs {
		outName := strings.ToLower(out.String())
		outName = strings.ReplaceAll(outName, "-", "")
		log.Tracef("found MIDI output: '%s'", outName)
		if strings.Contains(outName, strings.ToLower(name)) {
			m.Conn = out
			found = true
		}
	}
	if !found {
		err = fmt.Errorf("no MIDI output port found with name %s", name)
		return
	}

	if m.Conn == nil {
		err = fmt.Errorf("no MIDI output port found with name %s", name)
	}

	// OPEN PORT
	err = m.Conn.Open()

	return
}

func (m Midi) NoteOn(note int, velocity int) {
	if m.Conn == nil {
		return
	}
	// Send Note On message with the specified note and velocity
	log.Tracef("[%s] note_on n=%d,v=%d", m.Name, note, velocity)

	err := m.Conn.Send(midi.NoteOn(uint8(m.Channel), uint8(note), uint8(velocity)))
	if err != nil {
		log.Errorf("Failed to send Note On: %v", err)
	}
}

func (m Midi) NoteOff(note int) {
	if m.Conn == nil {
		return
	}
	// Send Note Off message with the specified note
	log.Tracef("[%s] note_ff n=%d", m.Name, note)
	err := m.Conn.Send(midi.NoteOff(uint8(m.Channel), uint8(note)))
	if err != nil {
		log.Errorf("Failed to send Note Off: %v", err)
	}
}

func (m Midi) Set(param string, value int) {
	if m.Conn == nil {
		return
	}
	// Send Control Change message with the specified controller and value
	// err := m.Conn.Send(midi.ControlChange(0, uint8(value), uint8(value)))
	// if err != nil {
	// 	log.Errorf("Failed to send Control Change: %v", err)
	// }
}
