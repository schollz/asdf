package step

import (
	"math/rand"

	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/param"
)

type Note struct {
	Midi   int
	Name   string
	Params []param.Param
}

type Step struct {
	Notes     []Note
	Emitters  []emitter.Emitter
	TimeStart float64
	TimeEnd   float64
	TimeLast  float64
}

func (n Note) GetParamNext(name string, defaultValue int) int {
	for _, p := range n.Params {
		if p.Name == name {
			return p.Next()
		}
	}
	return defaultValue
}

func NewNote(midi int, name string, params []param.Param) Note {
	return Note{
		Midi:   midi,
		Name:   name,
		Params: params,
	}
}

func NewStep(notes []Note, timeStart float64, timeEnd float64, emitters []emitter.Emitter) Step {
	return Step{
		Notes:     notes,
		TimeStart: timeStart,
		TimeEnd:   timeEnd,
		Emitters:  emitters,
	}
}

func (s *Step) Play(timeCurrent float64) {
	if s.TimeLast < s.TimeStart && timeCurrent >= s.TimeStart {
		for _, note := range s.Notes {

			// skip if probability is not met
			probability := note.GetParamNext("probability", 100)
			if probability < rand.Intn(100) {
				continue
			}

			// check if transpose parameter exists
			transpose := note.GetParamNext("transpose", 0)
			noteMidi := note.Midi
			if transpose != 0 {
				noteMidi += transpose
			}

			// check if velocity parameter exists
			velocity := note.GetParamNext("velocity", 64)

			for _, emitter := range s.Emitters {
				emitter.NoteOn(noteMidi, velocity)
			}
		}
	}
	for _, note := range s.Notes {
		// check if note has a gate parameter
		timeEnd := s.TimeEnd
		gate := note.GetParamNext("gate", 100)
		if gate < 100 {
			timeEnd = s.TimeStart + (s.TimeEnd-s.TimeStart)*float64(gate)/100
		}

		if s.TimeLast < timeEnd && timeCurrent >= timeEnd {
			for _, emitter := range s.Emitters {
				emitter.NoteOff(note.Midi)
			}
		}
	}
	s.TimeLast = timeCurrent
}
