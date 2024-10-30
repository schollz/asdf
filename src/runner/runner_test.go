package runner

import (
	"context"
	"testing"
	"time"

	"github.com/schollz/asdf/src/fileparser"
	"github.com/schollz/asdf/src/sprocket"
)

func TestRun(t *testing.T) {
	filename := "../fileparser/test.txt"
	sequences, err := fileparser.Parse(filename)
	if err != nil {
		t.Errorf("error: %s", err)
	}

	sprock := sprocket.New(sequences.Sprockets)

	ctx, cancel := context.WithCancel(context.Background())

	// Run the function in a separate goroutine
	go func() {
		if err := sprock.Run(ctx); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	// Allow the function to run for 5 seconds, then cancel
	time.Sleep(3 * time.Second)
	cancel()

	// Wait a moment to ensure the goroutine finishes
	time.Sleep(1 * time.Second)

}
