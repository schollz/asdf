package arpeggio

import (
	"github.com/schollz/asdf/src/step"

	log "github.com/schollz/logger"
)

func Expand(s string) (result string, err error) {
	step, err := step.Parse(s)
	if err != nil {
		log.Error(err)
		return
	}
	result = step.Text
	return
}
