package note

import "testing"

func TestMatch(t *testing.T) {
	tests := []struct {
		note     string
		midiNear int
		expected Note
	}{
		{"c", 71, Note{Midi: 72, Name: "c5"}},
		{"c4", 0, Note{Midi: 60, Name: "c4"}},
		{"b-1", 0, Note{Midi: 11, Name: "b-1"}},
		{"g♭", 32, Note{Midi: 30, Name: "f#1"}},
		{"f#", 80, Note{Midi: 78, Name: "f#5"}},
		{"f#3", 80, Note{Midi: 54, Name: "f#3"}},
	}
	for _, test := range tests {
		note, err := Parse(test.note, test.midiNear)
		if test.midiNear == 0 {
			note, err = Parse(test.note)
		}
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if note.Midi != test.expected.Midi || note.Name != test.expected.Name {
			t.Errorf("test '%s':\nexpected\n\t%v\nbut got\n\t%v", test.note, test.expected, note)
		}
	}

}

func TestOctave(t *testing.T) {
	tests := []struct {
		expectedoctave int
		n              Note
	}{
		{4, Note{Midi: 60, Name: "c4"}},
		{5, Note{Midi: 72, Name: "c5"}},
		{1, Note{Midi: 35, Name: "b1"}},
		{-1, Note{Midi: 35 - 24, Name: "b-1"}},
		{2, Note{Midi: 47, Name: "b2"}},
	}
	for _, test := range tests {
		octave := test.n.Octave()
		if octave != test.expectedoctave {
			t.Errorf("expected %d but got %d", test.expectedoctave, octave)
		}
	}
}
