package step

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/note"
	"github.com/schollz/asdf/src/noteorchord"
	"github.com/schollz/asdf/src/param"
	log "github.com/schollz/logger"
)

type Step struct {
	TextOriginal string
	TextNote     string
	Notes        []note.Note
	IsNote       bool
	IsRest       bool
	IsLegato     bool
	Params       []param.Param
	Emitters     []emitter.Emitter
	Beats        float64
	BPM          float64
	BeatStart    float64
	BeatEnd      float64
	TimeStart    float64
	TimeEnd      float64
	TimeLast     float64
}

func (s Step) GetParamNext(name string, defaultValue int) int {
	for _, p := range s.Params {
		if p.Name == name {
			return p.Next()
		}
	}
	return defaultValue
}

func (s Step) HasParam(name string) bool {
	for _, p := range s.Params {
		if p.Name == name {
			return true
		}
	}
	return false
}

func (s *Step) RemoveParam(name string) {
	for i, p := range s.Params {
		if p.Name == name {
			s.Params = append(s.Params[:i], s.Params[i+1:]...)
			return
		}
	}
}

func NewStep(notes []note.Note, timeStart float64, timeEnd float64, emitters []emitter.Emitter) Step {
	return Step{
		Notes:     notes,
		TimeStart: timeStart,
		TimeEnd:   timeEnd,
		Emitters:  emitters,
	}
}

func (s Step) Info() string {
	notes := []string{}
	for _, n := range s.Notes {
		notes = append(notes, fmt.Sprint(n.Midi))
	}
	return fmt.Sprintf("%s (%.2f-%.2f beats) (%.2f-%.2f sec) %0.1f at %.0fbpm (%s)", s.String(), s.BeatStart, s.BeatEnd, s.TimeStart, s.TimeEnd, s.Beats, s.BPM, strings.Join(notes, ","))
}

func (s Step) String() string {
	v := strings.Builder{}
	if s.TextNote == "" {
		for _, n := range s.Notes {
			v.WriteString(n.Name)
		}
	} else {
		v.WriteString(s.TextNote)
	}
	for _, p := range s.Params {
		v.WriteString(".")
		v.WriteString(p.String())
	}
	return v.String()
}

func (s *Step) Play(timeCurrent float64) {
	if s.TimeLast < s.TimeStart && timeCurrent >= s.TimeStart {
		for _, note := range s.Notes {

			// skip if probability is not met
			probability := s.GetParamNext("probability", 100)
			if probability < rand.Intn(100) {
				continue
			}

			// check if transpose parameter exists
			transpose := s.GetParamNext("transpose", 0)
			noteMidi := note.Midi
			if transpose != 0 {
				noteMidi += transpose
			}

			// check if velocity parameter exists
			velocity := s.GetParamNext("velocity", 64)

			for _, emitter := range s.Emitters {
				emitter.NoteOn(noteMidi, velocity)
			}
		}
	}
	for _, note := range s.Notes {
		// check if note has a gate parameter
		timeEnd := s.TimeEnd
		gate := s.GetParamNext("gate", 100)
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

func Parse(s string, midiNears ...int) (step Step, err error) {
	midiNear := 60
	if len(midiNears) > 0 {
		midiNear = midiNears[0]
	}

	fields := strings.Split(s, ".")
	log.Tracef("[%s] fields: %v", s, fields)

	step.TextOriginal = s
	if len(fields[0]) > 0 {
		step.TextNote = fields[0]
		step.Notes, err = noteorchord.Parse(fields[0], midiNear)
		if err != nil {
			log.Error(err)
			return
		}

		step.IsLegato = fields[0] == "-"
		step.IsRest = fields[0] == "~"
		step.IsNote = !step.IsLegato && !step.IsRest && len(step.Notes) > 0
	}

	log.Tracef("[%s] notes: %v", s, step.Notes)

	params := []param.Param{}
	if len(fields) > 1 {
		for _, field := range fields[1:] {
			p, err := param.Parse(field)
			if err != nil {
				log.Error(err)
				continue
			}
			params = append(params, p)
		}
	}
	step.Params = params
	return
}
