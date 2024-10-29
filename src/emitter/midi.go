package emitter

import (
	"fmt"
	"strings"

	log "github.com/schollz/logger"
	"gitlab.com/gomidi/midi/v2"
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

	return
}

func (m *Midi) NoteOn(note int, velocity int) {
}
