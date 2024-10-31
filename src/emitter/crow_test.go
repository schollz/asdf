package emitter

import (
	"testing"
	"time"
)

func TestCrow(t *testing.T) {
	setupCrows()
	if !crowsSetup {
		return
	}
	c, err := NewCrow(1, 2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	c.Set("attack", 1000)
	c.Set("release", 1000)
	c.NoteOn(60, 100)
	err = CrowFlush()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	time.Sleep(3 * time.Second)
	c.NoteOff(60)
	err = CrowFlush()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
