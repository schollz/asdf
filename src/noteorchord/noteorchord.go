package noteorchord

import (
	"strings"

	"github.com/schollz/asdf/src/chord"
	"github.com/schollz/asdf/src/note"
)

func Parse(midiString string, midiNears ...int) (notes []note.Note, err error) {
	if len(midiString) == 0 {
		return
	}
	midiNear := 60
	if len(midiNears) > 0 {
		midiNear = midiNears[0]
	}
	_ = midiNear

	midiString = strings.TrimSpace(midiString)
	if midiString == "-" || midiString == "~" {
		return
	}

	// if the first character is capital, then it is a chord
	if midiString[0] >= 'A' && midiString[0] <= 'Z' {
		// parse chord
		notes, err = chord.Parse(midiString, midiNear)
		return
	}

	// can be a single midi note like "c" in which case we need to find the closest note to midiNear
	// or can be a single note like "c4" in which case we want an exact match
	// or can be a sequence of notes like "c4eg" in which case we want to need to split them
	midiString = strings.ToLower(midiString)

	// check if split if it has multiple of any letter [a-g]
	noteStrings := []string{}
	lastAdded := 0
	for i := 1; i < len(midiString); i++ {
		if midiString[i] >= 'a' && midiString[i] <= 'g' {
			noteStrings = append(noteStrings, midiString[lastAdded:i])
			lastAdded = i
		}
	}
	if lastAdded != len(midiString) {
		noteStrings = append(noteStrings, midiString[lastAdded:])
	}

	for _, noteString := range noteStrings {
		noteString = strings.TrimSpace(noteString)
		n, noteErr := note.Parse(noteString, midiNear)
		if noteErr != nil {
			err = noteErr
			return
		}
		notes = append(notes, n)
		midiNear = n.Midi
	}
	return
}
