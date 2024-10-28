package chord

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/schollz/asdf/src/note"
	log "github.com/schollz/logger"
)

var dbChords = [][]string{
	[]string{"1P 3M 5P", "major", "maj", "^", ""},
	[]string{"1P 3M 5P 7M", "major seventh", "maj7", "ma7", "Maj7", "^7"},
	[]string{"1P 3M 5P 7M 9M", "major ninth", "maj9", "^9"},
	[]string{"1P 3M 5P 7M 9M 13M", "major thirteenth", "maj13", "Maj13 ^13"},
	[]string{"1P 3M 5P 6M", "sixth", "6", "add6", "add13"},
	[]string{"1P 3M 5P 6M 9M", "sixth/ninth", "6/9", "69"},
	[]string{"1P 3M 6m 7M", "major seventh flat sixth", "maj7b6", "^7b6"},
	[]string{"1P 3M 5P 7M 11A", "major seventh sharp eleventh", "majs4", "^7#11", "maj7#11"},
	// ==Minor==
	// '''Normal'''
	[]string{"1P 3m 5P", "minor", "m", "min", "-"},
	[]string{"1P 3m 5P 7m", "minor seventh", "m7", "min7", "mi7", "-7"},
	[]string{"1P 3m 5P 7M", "minor/major seventh", "maj7", "majmaj7", "mM7", "mMaj7", "m/M7", "-^7"},
	[]string{"1P 3m 5P 6M", "minor sixth", "m6", "-6"},
	[]string{"1P 3m 5P 7m 9M", "minor ninth", "m9", "-9"},
	[]string{"1P 3m 5P 7M 9M", "minor/major ninth", "minmaj9", "mMaj9", "-^9"},
	[]string{"1P 3m 5P 7m 9M 11P", "minor eleventh", "m11", "-11"},
	[]string{"1P 3m 5P 7m 9M 13M", "minor thirteenth", "m13", "-13"},
	// '''Diminished'''
	[]string{"1P 3m 5d", "diminished", "dim", "°", "o"},
	[]string{"1P 3m 5d 7d", "diminished seventh", "dim7", "°7", "o7"},
	[]string{"1P 3m 5d 7m", "half-diminished", "m7b5", "ø", "-7b5", "h7", "h"},
	// ==Dominant/Seventh==
	// '''Normal'''
	[]string{"1P 3M 5P 7m", "dominant seventh", "7", "dom"},
	[]string{"1P 3M 5P 7m 9M", "dominant ninth", "9"},
	[]string{"1P 3M 5P 7m 9M 13M", "dominant thirteenth", "13"},
	[]string{"1P 3M 5P 7m 11A", "lydian dominant seventh", "7s11", "7#4"},
	// '''Altered'''
	[]string{"1P 3M 5P 7m 9m", "dominant flat ninth", "7b9"},
	[]string{"1P 3M 5P 7m 9A", "dominant sharp ninth", "7s9"},
	[]string{"1P 3M 7m 9m", "altered", "alt7"},
	// '''Suspended'''
	[]string{"1P 4P 5P", "suspended fourth", "sus4", "sus"},
	[]string{"1P 2M 5P", "suspended second", "sus2"},
	[]string{"1P 4P 5P 7m", "suspended fourth seventh", "7sus4", "7sus"},
	[]string{"1P 5P 7m 9M 11P", "eleventh", "11"},
	[]string{"1P 4P 5P 7m 9m", "suspended fourth flat ninth", "b9sus", "phryg", "7b9sus", "7b9sus4"},
	// ==Other==
	[]string{"1P 5P", "fifth", "5"},
	[]string{"1P 3M 5A", "augmented", "aug", "+", "+5", "^#5"},
	[]string{"1P 3m 5A", "minor augmented", "ms5", "-#5", "m+"},
	[]string{"1P 3M 5A 7M", "augmented seventh", "maj75", "maj7+5", "+maj7", "^7#5"},
	[]string{"1P 3M 5P 7M 9M 11A", "major sharp eleventh (lydian)", "maj9s11", "^9#11"},
	// ==Legacy==
	[]string{"1P 2M 4P 5P", "", "sus24", "sus4add9"},
	[]string{"1P 3M 5A 7M 9M", "", "maj9s5", "Maj9s5"},
	[]string{"1P 3M 5A 7m", "", "7s5", "+7", "7+", "7aug", "aug7"},
	[]string{"1P 3M 5A 7m 9A", "", "7s5s9", "7s9s5", "7alt"},
	[]string{"1P 3M 5A 7m 9M", "", "9s5", "9+"},
	[]string{"1P 3M 5A 7m 9M 11A", "", "9s5s11"},
	[]string{"1P 3M 5A 7m 9m", "", "7s5b9", "7b9s5"},
	[]string{"1P 3M 5A 7m 9m 11A", "", "7s5b9s11"},
	[]string{"1P 3M 5A 9A", "", "padds9"},
	[]string{"1P 3M 5A 9M", "", "ms5add9", "padd9"},
	[]string{"1P 3M 5P 6M 11A", "", "M6s11", "M6b5", "6s11", "6b5"},
	[]string{"1P 3M 5P 6M 7M 9M", "", "maj7add13"},
	[]string{"1P 3M 5P 6M 9M 11A", "", "69s11"},
	[]string{"1P 3m 5P 6M 9M", "", "m69", "-69"},
	[]string{"1P 3M 5P 6m 7m", "", "7b6"},
	[]string{"1P 3M 5P 7M 9A 11A", "", "maj7s9s11"},
	[]string{"1P 3M 5P 7M 9M 11A 13M", "", "M13s11", "maj13s11", "M13+4", "M13s4"},
	[]string{"1P 3M 5P 7M 9m", "", "maj7b9"},
	[]string{"1P 3M 5P 7m 11A 13m", "", "7s11b13", "7b5b13"},
	[]string{"1P 3M 5P 7m 13M", "", "7add6", "67", "7add13"},
	[]string{"1P 3M 5P 7m 9A 11A", "", "7s9s11", "7b5s9", "7s9b5"},
	[]string{"1P 3M 5P 7m 9A 11A 13M", "", "13s9s11"},
	[]string{"1P 3M 5P 7m 9A 11A 13m", "", "7s9s11b13"},
	[]string{"1P 3M 5P 7m 9A 13M", "", "13s9"},
	[]string{"1P 3M 5P 7m 9A 13m", "", "7s9b13"},
	[]string{"1P 3M 5P 7m 9M 11A", "", "9s11", "9+4", "9s4"},
	[]string{"1P 3M 5P 7m 9M 11A 13M", "", "13s11", "13+4", "13s4"},
	[]string{"1P 3M 5P 7m 9M 11A 13m", "", "9s11b13", "9b5b13"},
	[]string{"1P 3M 5P 7m 9m 11A", "", "7b9s11", "7b5b9", "7b9b5"},
	[]string{"1P 3M 5P 7m 9m 11A 13M", "", "13b9s11"},
	[]string{"1P 3M 5P 7m 9m 11A 13m", "", "7b9b13s11", "7b9s11b13", "7b5b9b13"},
	[]string{"1P 3M 5P 7m 9m 13M", "", "13b9"},
	[]string{"1P 3M 5P 7m 9m 13m", "", "7b9b13"},
	[]string{"1P 3M 5P 7m 9m 9A", "", "7b9s9"},
	[]string{"1P 3M 5P 9M", "", "Madd9", "2", "add9", "add2"},
	[]string{"1P 3M 5P 9m", "", "majaddb9"},
	[]string{"1P 3M 5d", "", "majb5"},
	[]string{"1P 3M 5d 6M 7m 9M", "", "13b5"},
	[]string{"1P 3M 5d 7M", "", "maj7b5"},
	[]string{"1P 3M 5d 7M 9M", "", "maj9b5"},
	[]string{"1P 3M 5d 7m", "", "7b5"},
	[]string{"1P 3M 5d 7m 9M", "", "9b5"},
	[]string{"1P 3M 7m", "", "7no5"},
	[]string{"1P 3M 7m 13m", "", "7b13"},
	[]string{"1P 3M 7m 9M", "", "9no5"},
	[]string{"1P 3M 7m 9M 13M", "", "13no5"},
	[]string{"1P 3M 7m 9M 13m", "", "9b13"},
	[]string{"1P 3m 4P 5P", "", "madd4"},
	[]string{"1P 3m 5P 6m 7M", "", "mmaj7b6"},
	[]string{"1P 3m 5P 6m 7M 9M", "", "mmaj9b6"},
	[]string{"1P 3m 5P 7m 11P", "", "m7add11", "m7add4"},
	[]string{"1P 3m 5P 9M", "", "madd9"},
	[]string{"1P 3m 5d 6M 7M", "", "o7maj7"},
	[]string{"1P 3m 5d 7M", "", "omaj7"},
	[]string{"1P 3m 6m 7M", "", "mb6maj7"},
	[]string{"1P 3m 6m 7m", "", "m7s5"},
	[]string{"1P 3m 6m 7m 9M", "", "m9s5"},
	[]string{"1P 3m 5A 7m 9M 11P", "", "m11A"},
	[]string{"1P 3m 6m 9m", "", "mb6b9"},
	[]string{"1P 2M 3m 5d 7m", "", "m9b5"},
	[]string{"1P 4P 5A 7M", "", "maj7s5sus4"},
	[]string{"1P 4P 5A 7M 9M", "", "maj9s5sus4"},
	[]string{"1P 4P 5A 7m", "", "7s5sus4"},
	[]string{"1P 4P 5P 7M", "", "maj7sus4"},
	[]string{"1P 4P 5P 7M 9M", "", "maj9sus4"},
	[]string{"1P 4P 5P 7m 9M", "", "9sus4", "9sus"},
	[]string{"1P 4P 5P 7m 9M 13M", "", "13sus4", "13sus"},
	[]string{"1P 4P 5P 7m 9m 13m", "", "7sus4b9b13", "7b9b13sus4"},
	[]string{"1P 4P 7m 10m", "", "4", "quartal"},
	[]string{"1P 5P 7m 9m 11P", "", "11b9"},
}

func Parse(chordString string, midiNear int) (result []note.Note, err error) {
	chordStringOriginal := chordString
	chordMatch := ""
	_ = chordMatch

	log.Tracef("chordString: %s", chordStringOriginal)

	octave := 4
	if midiNear != 0 {
		octave = midiNear/12 - 1
	}
	if strings.Contains(chordString, ";") {
		chordSplit := strings.Split(chordString, ";")
		chordString = chordSplit[0]
		if len(chordSplit) > 1 {
			octave, err = strconv.Atoi(chordSplit[1])
			if err != nil {
				return
			}
		}
	}
	log.Tracef("octave: %d", octave)
	log.Tracef("chordString: %s", chordString)

	transposeNote := ""
	if strings.Contains(chordString, "/") {
		chordSplit := strings.Split(chordString, "/")
		chordString = chordSplit[0]
		if len(chordSplit) > 1 {
			transposeNote = chordSplit[1]
		}
	}
	log.Tracef("transposeNote: %s", transposeNote)

	// find the root note name
	noteMatch := ""
	transposeNoteMatch := ""
	chordRest := ""
	for _, n := range note.AllNotes {
		if transposeNote != "" && len(n) > len(transposeNoteMatch) {
			if strings.ToLower(n) == transposeNote || n == transposeNote {
				transposeNoteMatch = n
			}
		}
		if len(n) > len(noteMatch) {
			// check if has prefix
			if strings.HasPrefix(chordString, n) {
				noteMatch = n
				chordRest = chordString[len(n):]
			}
			if strings.HasPrefix(strings.ToLower(chordString), n) {
				noteMatch = n
				chordRest = chordString[len(n):]
			}
		}
	}
	if noteMatch == "" {
		err = fmt.Errorf("no chord found")
	}
	log.Tracef("noteMatch: %s", noteMatch)
	log.Tracef("chordRest: %s", chordRest)

	// convert to canonical sharp scale
	// e.g. Fb -> E, Gs -> G#
	for i, n := range note.NotesScaleAcc1 {
		if noteMatch == n {
			noteMatch = note.NotesScaleSharp[i]
			break
		}
	}
	for i, n := range note.NotesScaleAcc2 {
		if noteMatch == n {
			noteMatch = note.NotesScaleSharp[i]
			break
		}
	}
	for i, n := range note.NotesScaleAcc3 {
		if noteMatch == n {
			noteMatch = note.NotesScaleSharp[i]
			break
		}
	}
	if transposeNoteMatch != "" {
		for i, n := range note.NotesScaleAcc1 {
			if transposeNoteMatch == n {
				transposeNoteMatch = note.NotesScaleSharp[i]
				break
			}
		}
		for i, n := range note.NotesScaleAcc2 {
			if transposeNoteMatch == n {
				transposeNoteMatch = note.NotesScaleSharp[i]
				break
			}
		}
		for i, n := range note.NotesScaleAcc3 {
			if transposeNoteMatch == n {
				transposeNoteMatch = note.NotesScaleSharp[i]
				break
			}
		}
	}
	log.Tracef("noteMatch: %s", noteMatch)
	log.Tracef("transposeNoteMatch: %s", transposeNoteMatch)

	// find longest matching chord pattern
	chordMatch = "" // (no chord match is major chord)
	chordIntervals := "1P 3M 5P"
	for _, chordType := range dbChords {
		for i, chordPattern := range chordType {
			if i > 1 {
				if len(chordPattern) > len(chordMatch) && strings.ToLower(chordRest) == strings.ToLower(chordPattern) {
					chordMatch = chordPattern
					chordIntervals = chordType[0]
				}
			}
		}
	}
	log.Tracef("chordMatch for %s: %s", chordRest, chordMatch)
	log.Tracef("chordIntervals: %s", chordIntervals)

	// find location of root
	rootPosition := 0
	for i, n := range note.NotesScaleSharp {
		if n == noteMatch {
			rootPosition = i
			break
		}
	}
	log.Tracef("rootPosition: %d", rootPosition)

	/** lua code
		-- find notes from intervals
	  whole_note_semitones={0,2,4,5,7,9,11,12}
	  notes_in_chord={}
	  for interval in string.gmatch(chord_intervals,"%S+") do
	    -- get major note position
	    major_note_position=(string.match(interval,"%d+")-1)%7+1
	    -- find semitones from root
	    semitones=whole_note_semitones[major_note_position]
	    -- adjust semitones based on interval
	    if string.match(interval,"m") then
	      semitones=semitones-1
	    elseif string.match(interval,"A") then
	      semitones=semitones+1
	    end
	    if self.debug then
	      print("interval: "..interval)
	      print("major_note_position: "..major_note_position)
	      print("semitones: "..semitones)
	      print("root_position+semitones: "..root_position+semitones)
	    end
	    -- get note in scale from root position
	    note_in_chord=self.notes_scale_sharp[root_position+semitones]
	    table.insert(notes_in_chord,note_in_chord)
	  end
	  **/

	// go code
	// find notes from intervals
	wholeNoteSemitones := []int{0, 2, 4, 5, 7, 9, 11, 12}
	notesInChord := []string{}
	for _, interval := range strings.Fields(chordIntervals) {
		// get major note position
		majorNotePosition, _ := strconv.Atoi(strings.TrimRight(interval, "mMAP"))
		majorNotePosition = ((majorNotePosition - 1) % 7) + 1
		// find semitones from root
		semitones := wholeNoteSemitones[majorNotePosition]
		// adjust semitones based on interval
		if strings.Contains(interval, "A") || strings.Contains(interval, "M") {
			semitones = semitones + 1
		}
		log.Trace("-------------------------")
		log.Tracef("interval: %s", interval)
		log.Tracef("majorNotePosition: %d", majorNotePosition)
		log.Tracef("semitones: %d", semitones)
		log.Tracef("rootPosition+semitones: %d", rootPosition+semitones)
		// get note in scale from root position
		noteInChord := note.NotesScaleSharp[(rootPosition+semitones)-2+12]
		notesInChord = append(notesInChord, noteInChord)
	}
	log.Tracef("notesInChord: %v", notesInChord)
	log.Tracef("note.NotesScaleSharp: %v", note.NotesScaleSharp)

	// if tranposition, rotate until new root
	if transposeNoteMatch != "" {
		foundNote := false
		for i := 0; i < len(notesInChord); i++ {
			if notesInChord[0] == transposeNoteMatch {
				foundNote = true
				break
			}
			notesInChord = append(notesInChord[1:], notesInChord[0])
		}
		if !foundNote {
			notesInChord = append([]string{transposeNoteMatch}, notesInChord...)
		}
	}
	log.Tracef("notesInChord: %v", notesInChord)

	// go code
	// convert to midi
	midiNotesInChord := []int{}
	lastNote := 0
	for i, n := range notesInChord {
		for _, d := range note.NoteDB {
			if d.MidiValue > lastNote &&
				(d.NameSharp == n+strconv.Itoa(octave) ||
					d.NameSharp == n+strconv.Itoa(octave+1) ||
					d.NameSharp == n+strconv.Itoa(octave+2) ||
					d.NameSharp == n+strconv.Itoa(octave+3)) {
				lastNote = d.MidiValue
				midiNotesInChord = append(midiNotesInChord, d.MidiValue)
				notesInChord[i] = d.NameSharp
				break
			}
		}
	}
	log.Tracef("midiNotesInChord: %v", midiNotesInChord)
	log.Tracef("notesInChord: %v", notesInChord)

	result = make([]note.Note, len(midiNotesInChord))
	for i, m := range midiNotesInChord {
		result[i] = note.Note{Midi: m, Name: strings.ToLower(notesInChord[i])}
	}
	log.Tracef("result: %v", result)

	return
}
