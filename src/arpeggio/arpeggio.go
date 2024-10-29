package arpeggio

import (
	"github.com/schollz/asdf/src/param"
	"github.com/schollz/asdf/src/step"

	log "github.com/schollz/logger"
)

func Expand(s string) (result string, err error) {
	step, err := step.Parse(s)
	if err != nil {
		log.Error(err)
		return
	}
	result = step.TextOriginal
	if !step.HasParam(param.PARAM_ARPEGGIO) {
		return
	}

	step.RemoveParam(param.PARAM_ARPEGGIO)
	stepString := step.String()

	result = stepString

	return
}
