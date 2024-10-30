package player

import (
	"sync"

	"github.com/schollz/asdf/src/emitter"
)

type Player struct {
	Emitters []emitter.Emitter
	notesOn  map[int]bool
	mu       sync.Mutex
}

func New(emitters []emitter.Emitter) Player {
	return Player{
		Emitters: emitters,
		notesOn:  make(map[int]bool),
	}
}

func (p *Player) NoteOn(note int, velocity int) {
	for _, e := range p.Emitters {
		e.NoteOn(note, velocity)
	}
	p.mu.Lock()
	p.notesOn[note] = true
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

func (p *Player) Reset() {
	for note := range p.notesOn {
		p.NoteOff(note)
	}
	p.mu.Lock()
	p.notesOn = make(map[int]bool)
	p.mu.Unlock()
}
