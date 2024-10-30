package sprocket

import (
	"context"
	"sync"
	"time"

	"github.com/loov/hrtime"
	"github.com/schollz/asdf/src/block"
	"github.com/schollz/asdf/src/player"
	log "github.com/schollz/logger"
)

type Sprocket struct {
	Name   string
	Block  block.Block
	Player player.Player
}

type Sprockets struct {
	Sprockets []Sprocket
	mu        sync.Mutex
}

func New(sprockets []Sprocket) Sprockets {
	return Sprockets{
		Sprockets: sprockets,
	}
}

func (s *Sprockets) Update(sprockets []Sprocket) {
	s.mu.Lock()
	s.Sprockets = sprockets
	s.mu.Unlock()
}

func (s *Sprockets) Run(ctx context.Context) (err error) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop() // Ensure ticker is stopped when function exits
	startTime := hrtime.Now()
	lastTime := -1.0
	for {
		select {
		case <-ctx.Done():
			log.Debugf("sprocket received done signal")
			return nil
		case <-ticker.C:
			currentTime := hrtime.Since(startTime).Seconds()
			err = s.update(lastTime, currentTime)
			if err != nil {
				log.Error(err)
				return
			}
			lastTime = currentTime
		}
	}
}

func (s *Sprockets) update(totalLast, totalTime float64) (err error) {
	s.mu.Lock()
	for i, sp := range s.Sprockets {
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
		for _, step := range sp.Block.Steps {
			if step.TimeStart > currentTimeLast && step.TimeStart <= currentTime {
				log.Debugf("%s[%d] playing step: %s", sp.Name, i, step.Info())
			}
		}
	}
	s.mu.Unlock()
	return
}