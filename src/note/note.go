package note

import (
	"fmt"
	"math"
	"strings"
)

const LEFT_GROUP = "["
const RIGHT_GROUP = "]"
const HOLD = "_"
const REST = "~"

type NoteInfo struct {
	MidiValue  int
	NameSharp  string
	Frequency  float64
	NamesOther []string
}

type Note struct {
	Midi         int    `json:"midi,omitempty"`
	Name         string `json:"name,omitempty"`
	NameOriginal string `json:"name_original,omitempty"`
	IsRest       bool   `json:"is_rest,omitempty"`
	IsLegato     bool   `json:"is_legato,omitempty"`
	Length       int    `json:"length,omitempty"`
}

func findMaxPrefix(a string, b string) string {
	i := 0
	for i < len(a) && i < len(b) {
		if a[i] != b[i] {
			break
		}
		i++
	}
	return a[:i]
}

func Parse(n string, midiNears ...int) (note Note, err error) {
	midiNear := 60
	if len(midiNears) > 0 {
		midiNear = midiNears[0]
	}
	nOriginal := n

	n = strings.Replace(n, "♭", "b", -1)
	n = strings.Replace(n, "#", "s", -1)

	for _, m := range NoteDB {
		for _, noteFullName := range append(m.NamesOther, strings.ToLower(m.NameSharp)) {
			if n == noteFullName {
				note = Note{Midi: m.MidiValue, NameOriginal: nOriginal, Name: strings.ToLower(m.NameSharp)}
				return
			}
		}
	}

	// find closes to midiNear
	newnote := Note{Midi: 300, Name: ""}
	closestDistance := math.Inf(1)
	for _, m := range NoteDB {
		for octave := -1; octave <= 10; octave++ {
			for _, noteFullName := range append(m.NamesOther, strings.ToLower(m.NameSharp)) {
				noteFullNameWithoutOctave := noteFullName
				noteFullNameWithoutOctave = strings.Split(noteFullNameWithoutOctave, "-")[0]
				// remove numbers
				noteFullNameWithoutOctave = strings.Map(func(r rune) rune {
					if r >= '0' && r <= '9' {
						return -1
					}
					return r
				}, noteFullNameWithoutOctave)
				if noteFullNameWithoutOctave == n {
					if math.Abs(float64(m.MidiValue-midiNear)) < closestDistance {
						closestDistance = math.Abs(float64(m.MidiValue - midiNear))
						newnote = Note{Midi: m.MidiValue, Name: strings.ToLower(m.NameSharp)}
						// fmt.Println(newnote, closestDistance, noteFullName)
					}

				}
			}
		}
	}
	if newnote.Midi != 300 {
		note = newnote
		note.NameOriginal = nOriginal
	} else {
		err = fmt.Errorf("parsemidi could not parse %s", n)
	}
	return
}

var NoteDB = []NoteInfo{
	NoteInfo{MidiValue: -12, NameSharp: "C-2", Frequency: 4.0879, NamesOther: []string{"c-2"}},
	NoteInfo{MidiValue: -11, NameSharp: "C#-2", Frequency: 4.331, NamesOther: []string{"cs-2", "db-2"}},
	NoteInfo{MidiValue: -10, NameSharp: "D-2", Frequency: 4.5885, NamesOther: []string{"d-2"}},
	NoteInfo{MidiValue: -9, NameSharp: "D#-2", Frequency: 4.8614, NamesOther: []string{"ds-2", "eb-2"}},
	NoteInfo{MidiValue: -8, NameSharp: "E-2", Frequency: 5.1504, NamesOther: []string{"e-2", "fb-2"}},
	NoteInfo{MidiValue: -7, NameSharp: "F-2", Frequency: 5.4567, NamesOther: []string{"f-2"}},
	NoteInfo{MidiValue: -6, NameSharp: "F#-2", Frequency: 5.7812, NamesOther: []string{"fs-2", "gb-2"}},
	NoteInfo{MidiValue: -5, NameSharp: "G-2", Frequency: 6.1249, NamesOther: []string{"g-2"}},
	NoteInfo{MidiValue: -4, NameSharp: "G#-2", Frequency: 6.4891, NamesOther: []string{"gs-2", "ab-2"}},
	NoteInfo{MidiValue: -3, NameSharp: "A-2", Frequency: 6.875, NamesOther: []string{"a-2"}},
	NoteInfo{MidiValue: -2, NameSharp: "A#-2", Frequency: 7.2838, NamesOther: []string{"as-2", "bb-2"}},
	NoteInfo{MidiValue: -1, NameSharp: "B-2", Frequency: 7.7169, NamesOther: []string{"b-2", "cb-2"}},
	NoteInfo{MidiValue: 0, NameSharp: "C-1", Frequency: 8.1758, NamesOther: []string{"c-1"}},
	NoteInfo{MidiValue: 1, NameSharp: "C#-1", Frequency: 8.662, NamesOther: []string{"cs-1", "db-1"}},
	NoteInfo{MidiValue: 2, NameSharp: "D-1", Frequency: 9.177, NamesOther: []string{"d-1"}},
	NoteInfo{MidiValue: 3, NameSharp: "D#-1", Frequency: 9.7227, NamesOther: []string{"ds-1", "eb-1"}},
	NoteInfo{MidiValue: 4, NameSharp: "E-1", Frequency: 10.3009, NamesOther: []string{"e-1", "fb-1"}},
	NoteInfo{MidiValue: 5, NameSharp: "F-1", Frequency: 10.9134, NamesOther: []string{"f-1"}},
	NoteInfo{MidiValue: 6, NameSharp: "F#-1", Frequency: 11.5623, NamesOther: []string{"fs-1", "gb-1"}},
	NoteInfo{MidiValue: 7, NameSharp: "G-1", Frequency: 12.2499, NamesOther: []string{"g-1"}},
	NoteInfo{MidiValue: 8, NameSharp: "G#-1", Frequency: 12.9783, NamesOther: []string{"gs-1", "ab-1"}},
	NoteInfo{MidiValue: 9, NameSharp: "A-1", Frequency: 13.75, NamesOther: []string{"a-1"}},
	NoteInfo{MidiValue: 10, NameSharp: "A#-1", Frequency: 14.5676, NamesOther: []string{"as-1", "bb-1"}},
	NoteInfo{MidiValue: 11, NameSharp: "B-1", Frequency: 15.4339, NamesOther: []string{"b-1", "cb-1"}},
	NoteInfo{MidiValue: 12, NameSharp: "C0", Frequency: 16.351, NamesOther: []string{"c0"}},
	NoteInfo{MidiValue: 13, NameSharp: "C#0", Frequency: 17.324, NamesOther: []string{"cs0", "db0"}},
	NoteInfo{MidiValue: 14, NameSharp: "D0", Frequency: 18.354, NamesOther: []string{"d0"}},
	NoteInfo{MidiValue: 15, NameSharp: "D#0", Frequency: 19.445, NamesOther: []string{"ds0", "eb0"}},
	NoteInfo{MidiValue: 16, NameSharp: "E0", Frequency: 20.601, NamesOther: []string{"e0", "fb0"}},
	NoteInfo{MidiValue: 17, NameSharp: "F0", Frequency: 21.827, NamesOther: []string{"f0"}},
	NoteInfo{MidiValue: 18, NameSharp: "F#0", Frequency: 23.124, NamesOther: []string{"fs0", "gb0"}},
	NoteInfo{MidiValue: 19, NameSharp: "G0", Frequency: 24.499, NamesOther: []string{"g0"}},
	NoteInfo{MidiValue: 20, NameSharp: "G#0", Frequency: 25.956, NamesOther: []string{"gs0", "ab0"}},
	NoteInfo{MidiValue: 21, NameSharp: "A0", Frequency: 27.5, NamesOther: []string{"a0"}},
	NoteInfo{MidiValue: 22, NameSharp: "A#0", Frequency: 29.135, NamesOther: []string{"as0", "bb0"}},
	NoteInfo{MidiValue: 23, NameSharp: "B0", Frequency: 30.868, NamesOther: []string{"b0", "cb0"}},
	NoteInfo{MidiValue: 24, NameSharp: "C1", Frequency: 32.703, NamesOther: []string{"c1"}},
	NoteInfo{MidiValue: 25, NameSharp: "C#1", Frequency: 34.648, NamesOther: []string{"cs1", "db1"}},
	NoteInfo{MidiValue: 26, NameSharp: "D1", Frequency: 36.708, NamesOther: []string{"d1"}},
	NoteInfo{MidiValue: 27, NameSharp: "D#1", Frequency: 38.891, NamesOther: []string{"ds1", "eb1"}},
	NoteInfo{MidiValue: 28, NameSharp: "E1", Frequency: 41.203, NamesOther: []string{"e1", "fb1"}},
	NoteInfo{MidiValue: 29, NameSharp: "F1", Frequency: 43.654, NamesOther: []string{"f1"}},
	NoteInfo{MidiValue: 30, NameSharp: "F#1", Frequency: 46.249, NamesOther: []string{"fs1", "gb1"}},
	NoteInfo{MidiValue: 31, NameSharp: "G1", Frequency: 48.999, NamesOther: []string{"g1"}},
	NoteInfo{MidiValue: 32, NameSharp: "G#1", Frequency: 51.913, NamesOther: []string{"gs1", "ab1"}},
	NoteInfo{MidiValue: 33, NameSharp: "A1", Frequency: 55, NamesOther: []string{"a1"}},
	NoteInfo{MidiValue: 34, NameSharp: "A#1", Frequency: 58.27, NamesOther: []string{"as1", "bb1"}},
	NoteInfo{MidiValue: 35, NameSharp: "B1", Frequency: 61.735, NamesOther: []string{"b1", "cb1"}},
	NoteInfo{MidiValue: 36, NameSharp: "C2", Frequency: 65.406, NamesOther: []string{"c2"}},
	NoteInfo{MidiValue: 37, NameSharp: "C#2", Frequency: 69.296, NamesOther: []string{"cs2", "db2"}},
	NoteInfo{MidiValue: 38, NameSharp: "D2", Frequency: 73.416, NamesOther: []string{"d2"}},
	NoteInfo{MidiValue: 39, NameSharp: "D#2", Frequency: 77.782, NamesOther: []string{"ds2", "eb2"}},
	NoteInfo{MidiValue: 40, NameSharp: "E2", Frequency: 82.407, NamesOther: []string{"e2", "fb2"}},
	NoteInfo{MidiValue: 41, NameSharp: "F2", Frequency: 87.307, NamesOther: []string{"f2"}},
	NoteInfo{MidiValue: 42, NameSharp: "F#2", Frequency: 92.499, NamesOther: []string{"fs2", "gb2"}},
	NoteInfo{MidiValue: 43, NameSharp: "G2", Frequency: 97.999, NamesOther: []string{"g2"}},
	NoteInfo{MidiValue: 44, NameSharp: "G#2", Frequency: 103.826, NamesOther: []string{"gs2", "ab2"}},
	NoteInfo{MidiValue: 45, NameSharp: "A2", Frequency: 110, NamesOther: []string{"a2"}},
	NoteInfo{MidiValue: 46, NameSharp: "A#2", Frequency: 116.541, NamesOther: []string{"as2", "bb2"}},
	NoteInfo{MidiValue: 47, NameSharp: "B2", Frequency: 123.471, NamesOther: []string{"b2", "cb2"}},
	NoteInfo{MidiValue: 48, NameSharp: "C3", Frequency: 130.813, NamesOther: []string{"c3"}},
	NoteInfo{MidiValue: 49, NameSharp: "C#3", Frequency: 138.591, NamesOther: []string{"cs3", "db3"}},
	NoteInfo{MidiValue: 50, NameSharp: "D3", Frequency: 146.832, NamesOther: []string{"d3"}},
	NoteInfo{MidiValue: 51, NameSharp: "D#3", Frequency: 155.563, NamesOther: []string{"ds3", "eb3"}},
	NoteInfo{MidiValue: 52, NameSharp: "E3", Frequency: 164.814, NamesOther: []string{"e3", "fb3"}},
	NoteInfo{MidiValue: 53, NameSharp: "F3", Frequency: 174.614, NamesOther: []string{"f3"}},
	NoteInfo{MidiValue: 54, NameSharp: "F#3", Frequency: 184.997, NamesOther: []string{"fs3", "gb3"}},
	NoteInfo{MidiValue: 55, NameSharp: "G3", Frequency: 195.998, NamesOther: []string{"g3"}},
	NoteInfo{MidiValue: 56, NameSharp: "G#3", Frequency: 207.652, NamesOther: []string{"gs3", "ab3"}},
	NoteInfo{MidiValue: 57, NameSharp: "A3", Frequency: 220, NamesOther: []string{"a3"}},
	NoteInfo{MidiValue: 58, NameSharp: "A#3", Frequency: 233.082, NamesOther: []string{"as3", "bb3"}},
	NoteInfo{MidiValue: 59, NameSharp: "B3", Frequency: 246.942, NamesOther: []string{"b3", "cb3"}},
	NoteInfo{MidiValue: 60, NameSharp: "C4", Frequency: 261.626, NamesOther: []string{"c4"}},
	NoteInfo{MidiValue: 61, NameSharp: "C#4", Frequency: 277.183, NamesOther: []string{"cs4", "db4"}},
	NoteInfo{MidiValue: 62, NameSharp: "D4", Frequency: 293.665, NamesOther: []string{"d4"}},
	NoteInfo{MidiValue: 63, NameSharp: "D#4", Frequency: 311.127, NamesOther: []string{"ds4", "eb4"}},
	NoteInfo{MidiValue: 64, NameSharp: "E4", Frequency: 329.628, NamesOther: []string{"e4", "fb4"}},
	NoteInfo{MidiValue: 65, NameSharp: "F4", Frequency: 349.228, NamesOther: []string{"f4"}},
	NoteInfo{MidiValue: 66, NameSharp: "F#4", Frequency: 369.994, NamesOther: []string{"fs4", "gb4"}},
	NoteInfo{MidiValue: 67, NameSharp: "G4", Frequency: 391.995, NamesOther: []string{"g4"}},
	NoteInfo{MidiValue: 68, NameSharp: "G#4", Frequency: 415.305, NamesOther: []string{"gs4", "ab4"}},
	NoteInfo{MidiValue: 69, NameSharp: "A4", Frequency: 440, NamesOther: []string{"a4"}},
	NoteInfo{MidiValue: 70, NameSharp: "A#4", Frequency: 466.164, NamesOther: []string{"as4", "bb4"}},
	NoteInfo{MidiValue: 71, NameSharp: "B4", Frequency: 493.883, NamesOther: []string{"b4", "cb4"}},
	NoteInfo{MidiValue: 72, NameSharp: "C5", Frequency: 523.251, NamesOther: []string{"c5"}},
	NoteInfo{MidiValue: 73, NameSharp: "C#5", Frequency: 554.365, NamesOther: []string{"cs5", "db5"}},
	NoteInfo{MidiValue: 74, NameSharp: "D5", Frequency: 587.33, NamesOther: []string{"d5"}},
	NoteInfo{MidiValue: 75, NameSharp: "D#5", Frequency: 622.254, NamesOther: []string{"ds5", "eb5"}},
	NoteInfo{MidiValue: 76, NameSharp: "E5", Frequency: 659.255, NamesOther: []string{"e5", "fb5"}},
	NoteInfo{MidiValue: 77, NameSharp: "F5", Frequency: 698.456, NamesOther: []string{"f5"}},
	NoteInfo{MidiValue: 78, NameSharp: "F#5", Frequency: 739.989, NamesOther: []string{"fs5", "gb5"}},
	NoteInfo{MidiValue: 79, NameSharp: "G5", Frequency: 783.991, NamesOther: []string{"g5"}},
	NoteInfo{MidiValue: 80, NameSharp: "G#5", Frequency: 830.609, NamesOther: []string{"gs5", "ab5"}},
	NoteInfo{MidiValue: 81, NameSharp: "A5", Frequency: 880, NamesOther: []string{"a5"}},
	NoteInfo{MidiValue: 82, NameSharp: "A#5", Frequency: 932.328, NamesOther: []string{"as5", "bb5"}},
	NoteInfo{MidiValue: 83, NameSharp: "B5", Frequency: 987.767, NamesOther: []string{"b5", "cb5"}},
	NoteInfo{MidiValue: 84, NameSharp: "C6", Frequency: 1046.502, NamesOther: []string{"c6"}},
	NoteInfo{MidiValue: 85, NameSharp: "C#6", Frequency: 1108.731, NamesOther: []string{"cs6", "db6"}},
	NoteInfo{MidiValue: 86, NameSharp: "D6", Frequency: 1174.659, NamesOther: []string{"d6"}},
	NoteInfo{MidiValue: 87, NameSharp: "D#6", Frequency: 1244.508, NamesOther: []string{"ds6", "eb6"}},
	NoteInfo{MidiValue: 88, NameSharp: "E6", Frequency: 1318.51, NamesOther: []string{"e6", "fb6"}},
	NoteInfo{MidiValue: 89, NameSharp: "F6", Frequency: 1396.913, NamesOther: []string{"f6"}},
	NoteInfo{MidiValue: 90, NameSharp: "F#6", Frequency: 1479.978, NamesOther: []string{"fs6", "gb6"}},
	NoteInfo{MidiValue: 91, NameSharp: "G6", Frequency: 1567.982, NamesOther: []string{"g6"}},
	NoteInfo{MidiValue: 92, NameSharp: "G#6", Frequency: 1661.219, NamesOther: []string{"gs6", "ab6"}},
	NoteInfo{MidiValue: 93, NameSharp: "A6", Frequency: 1760, NamesOther: []string{"a6"}},
	NoteInfo{MidiValue: 94, NameSharp: "A#6", Frequency: 1864.655, NamesOther: []string{"as6", "bb6"}},
	NoteInfo{MidiValue: 95, NameSharp: "B6", Frequency: 1975.533, NamesOther: []string{"b6", "cb6"}},
	NoteInfo{MidiValue: 96, NameSharp: "C7", Frequency: 2093.005, NamesOther: []string{"c7"}},
	NoteInfo{MidiValue: 97, NameSharp: "C#7", Frequency: 2217.461, NamesOther: []string{"cs7", "db7"}},
	NoteInfo{MidiValue: 98, NameSharp: "D7", Frequency: 2349.318, NamesOther: []string{"d7"}},
	NoteInfo{MidiValue: 99, NameSharp: "D#7", Frequency: 2489.016, NamesOther: []string{"ds7", "eb7"}},
	NoteInfo{MidiValue: 100, NameSharp: "E7", Frequency: 2637.021, NamesOther: []string{"e7", "fb7"}},
	NoteInfo{MidiValue: 101, NameSharp: "F7", Frequency: 2793.826, NamesOther: []string{"f7"}},
	NoteInfo{MidiValue: 102, NameSharp: "F#7", Frequency: 2959.955, NamesOther: []string{"fs7", "gb7"}},
	NoteInfo{MidiValue: 103, NameSharp: "G7", Frequency: 3135.964, NamesOther: []string{"g7"}},
	NoteInfo{MidiValue: 104, NameSharp: "G#7", Frequency: 3322.438, NamesOther: []string{"gs7", "ab7"}},
	NoteInfo{MidiValue: 105, NameSharp: "A7", Frequency: 3520, NamesOther: []string{"a7"}},
	NoteInfo{MidiValue: 106, NameSharp: "A#7", Frequency: 3729.31, NamesOther: []string{"as7", "bb7"}},
	NoteInfo{MidiValue: 107, NameSharp: "B7", Frequency: 3951.066, NamesOther: []string{"b7", "cb7"}},
	NoteInfo{MidiValue: 108, NameSharp: "C8", Frequency: 4186.009, NamesOther: []string{"c8"}},
	NoteInfo{MidiValue: 109, NameSharp: "C#8", Frequency: 4434.922, NamesOther: []string{"cs8", "db8"}},
	NoteInfo{MidiValue: 110, NameSharp: "D8", Frequency: 4698.636, NamesOther: []string{"d8"}},
	NoteInfo{MidiValue: 111, NameSharp: "D#8", Frequency: 4978.032, NamesOther: []string{"ds8", "eb8"}},
	NoteInfo{MidiValue: 112, NameSharp: "E8", Frequency: 5274.042, NamesOther: []string{"e8", "fb8"}},
	NoteInfo{MidiValue: 113, NameSharp: "F8", Frequency: 5587.652, NamesOther: []string{"f8"}},
	NoteInfo{MidiValue: 114, NameSharp: "F#8", Frequency: 5919.91, NamesOther: []string{"fs8", "gb8"}},
	NoteInfo{MidiValue: 115, NameSharp: "G8", Frequency: 6271.928, NamesOther: []string{"g8"}},
	NoteInfo{MidiValue: 116, NameSharp: "G#8", Frequency: 6644.876, NamesOther: []string{"gs8", "ab8"}},
	NoteInfo{MidiValue: 117, NameSharp: "A8", Frequency: 7040, NamesOther: []string{"a8"}},
	NoteInfo{MidiValue: 118, NameSharp: "A#8", Frequency: 7458.62, NamesOther: []string{"as8", "bb8"}},
	NoteInfo{MidiValue: 119, NameSharp: "B8", Frequency: 7902.132, NamesOther: []string{"b8", "cb8"}},
	NoteInfo{MidiValue: 120, NameSharp: "C9", Frequency: 8372.018, NamesOther: []string{"c9"}},
	NoteInfo{MidiValue: 121, NameSharp: "C#9", Frequency: 8869.844, NamesOther: []string{"cs9", "db9"}},
	NoteInfo{MidiValue: 122, NameSharp: "D9", Frequency: 9397.272, NamesOther: []string{"d9"}},
	NoteInfo{MidiValue: 123, NameSharp: "D#9", Frequency: 9956.064, NamesOther: []string{"ds9", "eb9"}},
	NoteInfo{MidiValue: 124, NameSharp: "E9", Frequency: 10548.084, NamesOther: []string{"e9", "fb9"}},
	NoteInfo{MidiValue: 125, NameSharp: "F9", Frequency: 11175.304, NamesOther: []string{"f9"}},
	NoteInfo{MidiValue: 126, NameSharp: "F#9", Frequency: 11839.82, NamesOther: []string{"fs9", "gb9"}},
	NoteInfo{MidiValue: 127, NameSharp: "G9", Frequency: 12543.856, NamesOther: []string{"g9"}},
	NoteInfo{MidiValue: 128, NameSharp: "G#9", Frequency: 13289.752, NamesOther: []string{"gs9", "ab9"}},
	NoteInfo{MidiValue: 129, NameSharp: "A9", Frequency: 14080, NamesOther: []string{"a9"}},
	NoteInfo{MidiValue: 130, NameSharp: "A#9", Frequency: 14917.24, NamesOther: []string{"as9", "bb9"}},
	NoteInfo{MidiValue: 131, NameSharp: "B9", Frequency: 15804.264, NamesOther: []string{"b9", "cb9"}},
}

var NotesWhite = []string{"C", "D", "E", "F", "G", "A", "B"}
var NotesScaleSharp = []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
var NotesScaleAcc1 = []string{"B#", "Db", "D", "Eb", "Fb", "E#", "Gb", "G", "Ab", "A", "Bb", "Cb"}
var NotesScaleAcc2 = []string{"C", "Cs", "D", "Ds", "E", "F", "Fs", "G", "Gs", "A", "As", "B"}
var NotesScaleAcc3 = []string{"Bs", "Db", "D", "Eb", "Fb", "Es", "Gb", "G", "Ab", "A", "Bb", "Cb"}
var NotesAdds = []string{"", "#", "b", "s"}
var AllNotes = []string{}

func init() {
	for _, n := range NotesWhite {
		for _, a := range NotesAdds {
			AllNotes = append(AllNotes, n+a)
		}
	}
}
