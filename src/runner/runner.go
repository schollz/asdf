package runner

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/fileparser"
	"github.com/schollz/asdf/src/sprocket"
	log "github.com/schollz/logger"
)

func hash256(filename string) string {
	b, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256(b))
}

func Run(filename string) (err error) {
	sequences, err := fileparser.Parse(filename)
	if err != nil {
		log.Error(err)
		return
	}
	filenameHash := hash256(filename)
	sprock := sprocket.New(sequences.Sprockets)

	ctx, cancel := context.WithCancel(context.Background())

	// Run the function in a separate goroutine
	go func() {
		if err := sprock.Run(ctx); err != nil {
			log.Error(err)
		}
	}()

	// check hash for update
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)
			if hash256(filename) != filenameHash {
				log.Debug("file changed, reloading")
				sequences, err = fileparser.Parse(filename)
				if err != nil {
					log.Error(err)
				} else {
					filenameHash = hash256(filename)
					sprock.Update(sequences.Sprockets)
				}
			}
		}
	}()

	// wait for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Tracef("%+v", sig)
			cancel()
		}
	}()

	// wait for context to be done
	<-ctx.Done()
	emitter.CrowClose()
	time.Sleep(1 * time.Second)
	return
}
