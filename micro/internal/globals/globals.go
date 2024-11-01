package globals

import (
	"context"
	"time"

	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/fileparser"
	"github.com/schollz/asdf/src/sprocket"
	log "github.com/schollz/logger"
)

var Sprock sprocket.Sprockets
var SprockRunning bool
var Cancel context.CancelFunc
var ctx context.Context

func DoCancel() {
	if Cancel != nil {
		Cancel()
		<-ctx.Done()
		emitter.CrowClose()
		time.Sleep(1 * time.Second)
	}
}

func ProcessFilename(filename string) (err error) {
	sequences, err := fileparser.Parse(filename)
	if err != nil {
		log.Error(err)
		return
	} else {
		Sprock.Update(sequences.Sprockets)
		if !SprockRunning {
			SprockRunning = true
			ctx, Cancel = context.WithCancel(context.Background())

			// Run the function in a separate goroutine
			go func() {
				if err := Sprock.Run(ctx); err != nil {
					log.Error(err)
				}
			}()

		}
	}
	return
}
