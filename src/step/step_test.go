package step

import (
	"testing"

	"github.com/schollz/asdf/src/note"
	"github.com/schollz/asdf/src/param"
	log "github.com/schollz/logger"
)

func TestString(t *testing.T) {
	log.SetLevel("trace")
	tests := []struct {
		line     string
		expected string
	}{
		{"Cmaj.gate50,30.prob30", "Cmaj.gate50,30.probability30"},
		{".gate50,30.prob30", ".gate50,30.probability30"},
	}
	for _, test := range tests {
		step, err := Parse(test.line)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if step.String() != test.expected {
			t.Errorf("expected %v but got %v", test.expected, step.String())
		}
		log.Debugf("%+v", step)
	}
	// remove gate
	tests = []struct {
		line     string
		expected string
	}{
		{"Cmaj.gate40.prob30", "Cmaj.probability30"},
		{".gate56.prob30", ".probability30"},
	}
	for _, test := range tests {
		step, err := Parse(test.line)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		step.RemoveParam(param.PARAM_GATE)
		if step.String() != test.expected {
			t.Errorf("expected %v but got %v", test.expected, step.String())
		}
		log.Debugf("%+v", step)
	}

}
func TestParse(t *testing.T) {
	log.SetLevel("trace")
	tests := []struct {
		line     string
		expected Step
	}{
		{"Cmaj.gate50,30", Step{
			Notes: []note.Note{
				note.Note{Midi: 60},
				note.Note{Midi: 64},
				note.Note{Midi: 67},
			},
			Params: []param.Param{
				param.Param{Name: "gate", Values: []int{50, 30}},
			},
		}},
		{"c3d3e3.probability50", Step{
			Notes: []note.Note{
				note.Note{Midi: 48},
				note.Note{Midi: 50},
				note.Note{Midi: 52},
			},
			Params: []param.Param{
				param.Param{Name: "probability", Values: []int{50}},
			},
		}},
		{"-", Step{}},
		{".probability50", Step{
			Params: []param.Param{
				param.Param{Name: "probability", Values: []int{50}},
			},
		}},
	}
	for _, test := range tests {
		step, err := Parse(test.line)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for i, v := range step.Notes {
			if v.Midi != test.expected.Notes[i].Midi {
				t.Errorf("expected %v but got %v", test.expected.Notes[i].Midi, v.Midi)
			}
		}
		for i, v := range step.Params {
			if v.Name != test.expected.Params[i].Name {
				t.Errorf("expected %v but got %v", test.expected.Params[i].Name, v.Name)
			}
			for j, w := range v.Values {
				if w != test.expected.Params[i].Values[j] {
					t.Errorf("expected %v but got %v", test.expected.Params[i].Values[j], w)
				}
			}
		}
		log.Debugf("%+v", step)
	}
}
