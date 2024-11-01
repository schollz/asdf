package step

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/schollz/asdf/src/note"
	"github.com/schollz/asdf/src/noteorchord"
	"github.com/schollz/asdf/src/param"
	"github.com/schollz/asdf/src/player"
	log "github.com/schollz/logger"
)

type Step struct {
	TextOriginal string
	TextNote     string
	Notes        []note.Note
	MidiLast     []int
	IsNote       bool
	IsRest       bool
	IsLegato     bool
	Params       []param.Param
	Beats        float64
	BPM          float64
	BeatStart    float64
	BeatEnd      float64
	TimeStart    float64
	TimeEnd      float64
	TimeLast     float64
}

func (s *Step) GetParamNext(name string, defaultValue int) int {
	for i, p := range s.Params {
		if p.Name == name {
			v := s.Params[i].Next()
			if name == "transpose" {
				log.Debugf("transpose: %d", v)
			}
			return v
		}
	}
	return defaultValue
}

func (s *Step) GetParamCurrent(name string, defaultValue int) int {
	for i, p := range s.Params {
		if p.Name == name {
			return s.Params[i].Current()
		}
	}
	return defaultValue
}

func (s *Step) SetParm(name string, values []int) {
	for i, p := range s.Params {
		if p.Name == name {
			s.Params[i].Values = values
			return
		}
	}
	s.Params = append(s.Params, param.New(name, values))
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

func (s *Step) Play(timeLast float64, timeCurrent float64, play *player.Player) (notesOn []int, notesOff []int) {
	if timeLast <= s.TimeStart && timeCurrent >= s.TimeStart {
		for _, nn := range s.Notes {
			// adsr settings
			settings := []string{"attack", "decay", "release"}
			for _, setting := range settings {
				if s.HasParam(setting) {
					v := float64(s.GetParamNext(setting, 0))
					// scale to the duration of the note
					v = v * (s.TimeEnd - s.TimeStart) / 100.0
					play.Set(setting, v)
				}
			}
			if s.HasParam("sustain") {
				v := float64(s.GetParamNext("sustain", 0))
				play.Set("sustain", v/100.0)
			}

			// skip if probability is not met
			probability := s.GetParamNext("probability", 100)
			if probability < rand.Intn(100) {
				continue
			}

			// check if transpose parameter exists
			noteMidi := nn.Midi + s.GetParamNext("transpose", 0)

			// check if velocity parameter exists
			velocity := s.GetParamNext("velocity", 64)

			play.NoteOn(noteMidi, velocity)
			notesOn = append(notesOn, noteMidi)
		}
	}
	// check if note has a gate parameter
	timeEnd := s.TimeEnd
	gate := s.GetParamNext("gate", 100)
	if gate < 100 {
		timeEnd = s.TimeStart + (s.TimeEnd-s.TimeStart)*float64(gate)/100
	}
	if timeLast < timeEnd && timeCurrent >= timeEnd {
		transpose := s.GetParamCurrent("transpose", 0)
		for _, nn := range s.Notes {
			play.NoteOff(nn.Midi + transpose)
			notesOff = append(notesOff, nn.Midi+transpose)
		}
	}

	return
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
