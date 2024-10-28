package noteorchord

import (
	"testing"

	"github.com/schollz/asdf/src/note"
	log "github.com/schollz/logger"
)

func TestParse(t *testing.T) {
	log.SetLevel("debug")
	// table driven tests
	tests := []struct {
		midiString string
		midiNear   int
		expected   []note.Note
	}{
		{"c", 71, []note.Note{{Midi: 72, Name: "c5"}}},
		{"c6", 20, []note.Note{{Midi: 84, Name: "c6"}}},
		{"c", 62, []note.Note{{Midi: 60, Name: "c4"}}},
		{"d", 32, []note.Note{{Midi: 26, Name: "d1"}}},
		{"f#3", 32, []note.Note{{Midi: 54, Name: "f#3"}}},
		{"g7", 100, []note.Note{{Midi: 103, Name: "g7"}}},
		{"gb", 100, []note.Note{{Midi: 103, Name: "g7"}, {Midi: 107, Name: "b7"}}},
		{"g♭c", 100, []note.Note{{Midi: 102, Name: "f#7"}, {Midi: 96, Name: "c7"}}},
		{"c4eg", 52, []note.Note{
			{Midi: 60, Name: "c4"},
			{Midi: 64, Name: "e4"},
			{Midi: 67, Name: "g4"},
		}},
		{"ceg6", 52, []note.Note{
			{Midi: 48, Name: "c3"},
			{Midi: 52, Name: "e3"},
			{Midi: 91, Name: "g6"},
		}},
		{"Cmaj", 52, []note.Note{
			{Midi: 48, Name: "c3"},
			{Midi: 52, Name: "e3"},
			{Midi: 55, Name: "g3"},
		}},
	}
	for _, test := range tests {
		notes, err := Parse(test.midiString, test.midiNear)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			continue
		}
		if len(notes) != len(test.expected) {
			t.Errorf("expected %v, got %v", test.expected, notes)
			continue
		}
		for i, note := range notes {
			if note.Midi != test.expected[i].Midi || note.Name != test.expected[i].Name {
				t.Errorf("test: %s (%d), expected %v, got %v", test.midiString, test.midiNear, test.expected, notes)
				break
			}
		}
	}
}
