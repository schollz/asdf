package sprocket

import (
	"context"
	"sync"
	"time"

	"github.com/loov/hrtime"
	"github.com/schollz/asdf/src/block"
	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/note"
	"github.com/schollz/asdf/src/player"
	log "github.com/schollz/logger"
)

type Sprocket struct {
	Name   string
	Block  block.Block
	Player player.Player
}

type Sprockets struct {
	Sprockets   []Sprocket
	Playing     bool
	StartTime   time.Duration
	mu          sync.Mutex
	noteMarquee []string
}

func New(sprockets []Sprocket) Sprockets {
	return Sprockets{
		Playing:   true,
		Sprockets: sprockets,
	}
}

func (s *Sprockets) Update(sprockets []Sprocket) {
	s.mu.Lock()
	// reset current players
	for _, sp := range s.Sprockets {
		sp.Player.Reset()
	}
	s.Sprockets = sprockets
	s.mu.Unlock()
}

func (s *Sprockets) NotesOn() []string {
	return s.noteMarquee
}

func (s *Sprockets) AddToMarquee(midis []int) {
	for _, midi := range midis {
		n, _ := note.FromMidi(midi)
		s.noteMarquee = append(s.noteMarquee, n.Name)
		if len(s.noteMarquee) > 10 {
			s.noteMarquee = s.noteMarquee[1:]
		}
	}
}

func (s *Sprockets) Run(ctx context.Context) (err error) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop() // Ensure ticker is stopped when function exits
	s.StartTime = hrtime.Now()
	lastTime := -1.0
	for {
		select {
		case <-ctx.Done():
			// reset players
			for _, sp := range s.Sprockets {
				sp.Player.Reset()
			}
			log.Debugf("sprocket received done signal")
			return nil
		case <-ticker.C:
			if !s.Playing {
				continue
			}
			currentTime := hrtime.Since(s.StartTime).Seconds()
			s.update(lastTime, currentTime)
			if err != nil {
				log.Error(err)
				return
			}
			lastTime = currentTime
			emitter.CrowFlush()
		}
	}
}

func (s *Sprockets) Toggle(play ...bool) {
	wasPlaying := s.Playing
	var nowPlaying bool
	if len(play) > 0 {
		nowPlaying = play[0]
	} else {
		nowPlaying = !s.Playing
	}
	if !wasPlaying && nowPlaying {
		// reset start time
		s.StartTime = hrtime.Now()
	} else if wasPlaying && !nowPlaying {
		// reset players
		for _, sp := range s.Sprockets {
			sp.Player.Reset()
		}
	}
	s.Playing = nowPlaying
	log.Tracef("sprocket is playing: %v", s.Playing)
}

func (s *Sprockets) update(totalLast, totalTime float64) (err error) {
	s.mu.Lock()
	for _, sp := range s.Sprockets {
		currentTime := totalTime
		currentTimeLast := totalLast
		for {
			if currentTime > sp.Block.TotalTime {
				currentTime -= sp.Block.TotalTime
			} else {
				break
			}
		}
		for {
			if currentTimeLast > sp.Block.TotalTime {
				currentTimeLast -= sp.Block.TotalTime
			} else {
				break
			}
		}
		if currentTimeLast > currentTime {
			for _, step := range sp.Block.Steps {
				notesOn, _ := step.Play(currentTimeLast, currentTime+sp.Block.TotalTime, &sp.Player)
				s.AddToMarquee(notesOn)
			}
			for _, step := range sp.Block.Steps {
				notesOn, _ := step.Play(currentTimeLast-sp.Block.TotalTime, currentTime, &sp.Player)
				s.AddToMarquee(notesOn)
			}
		} else {
			for _, step := range sp.Block.Steps {
				notesOn, _ := step.Play(currentTimeLast, currentTime, &sp.Player)
				s.AddToMarquee(notesOn)
			}

		}

	}
	s.mu.Unlock()
	return
}
