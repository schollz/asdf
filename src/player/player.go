package player

import (
	"sort"
	"sync"
	"time"

	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/note"
)

type Player struct {
	Emitters []emitter.Emitter
	notesOn  map[int]time.Time
	mu       sync.Mutex
}

func New(emitters []emitter.Emitter) Player {
	return Player{
		Emitters: emitters,
		notesOn:  make(map[int]time.Time),
	}
}

func (p *Player) NoteOn(note int, velocity int) {
	for _, e := range p.Emitters {
		e.NoteOn(note, velocity)
	}
	p.mu.Lock()
	p.notesOn[note] = time.Now()
	p.mu.Unlock()
}

func (p *Player) NoteOff(note int) {
	for _, e := range p.Emitters {
		e.NoteOff(note)
	}
	p.mu.Lock()
	delete(p.notesOn, note)
	p.mu.Unlock()
}

func (p *Player) Set(param string, value float64) {
	for _, e := range p.Emitters {
		e.Set(param, value)
	}
}

func (p *Player) Reset() {
	for note := range p.notesOn {
		p.NoteOff(note)
	}
	p.mu.Lock()
	p.notesOn = make(map[int]time.Time)
	p.mu.Unlock()
}

type NoteSort struct {
	note note.Note
	t    time.Time
}

func (p *Player) NotesOn() []string {
	// sort notesOn by time
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.notesOn) == 0 {
		return []string{}
	}
	notes := make([]NoteSort, len(p.notesOn))
	i := 0
	for midi, t := range p.notesOn {
		n, errNote := note.FromMidi(midi)
		if errNote == nil {
			notes[i] = NoteSort{n, t}
			i++
		}
	}
	// sort notes by time
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].t.Before(notes[j].t)
	})
	// convert to strings
	strs := make([]string, len(notes))
	for i, n := range notes {
		strs[i] = n.note.Name
	}
	return strs
}
