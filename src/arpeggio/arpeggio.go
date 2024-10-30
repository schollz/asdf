package arpeggio

import (
	"strings"

	"github.com/schollz/asdf/src/note"
	"github.com/schollz/asdf/src/param"
	"github.com/schollz/asdf/src/step"

	log "github.com/schollz/logger"
)

func Expand(s string) (result string, err error) {
	stepOriginal, err := step.Parse(s)
	if err != nil {
		log.Error(err)
		return
	}
	result = stepOriginal.TextOriginal
	if !stepOriginal.HasParam(param.PARAM_ARPEGGIO) {
		return
	}
	stepOriginal.RemoveParam(param.PARAM_ARPEGGIO)

	// go through the arpeggio parameters
	paramsNew := make([]param.Param, 0)

	notes := make([]note.Note, 0)
	arpI := 0
	first := true
	for _, p := range stepOriginal.Params {
		if p.Name == param.PARAM_UP || p.Name == param.PARAM_DOWN || p.Name == param.PARAM_THUMB {
			// go up
			arpLength := p.Current()
			for i := 0; i < arpLength; i++ {
				if first {
					first = false
				} else if p.Name == param.PARAM_UP {
					arpI++
				} else if p.Name == param.PARAM_DOWN {
					arpI--
				} else if p.Name == param.PARAM_THUMB {
					arpI = 0
				}
				arpIPos := arpI
				if arpIPos < 0 {
					arpIPos = arpIPos * -1
				}
				noteMidi := stepOriginal.Notes[arpIPos%len(stepOriginal.Notes)].Midi + (int(arpI)/len(stepOriginal.Notes))*12
				if arpI < 0 {
					noteMidi -= 12
				}
				var noteNext note.Note
				noteNext, err = note.FromMidi(noteMidi)
				if err != nil {
					log.Error(err)
					return
				}
				log.Tracef("arpI: %d, noteMidi: %d, noteNext: %+v", arpI, noteMidi, noteNext)
				notes = append(notes, noteNext)
			}
		} else {
			paramsNew = append(paramsNew, p)

		}
	}

	newSteps := make([]step.Step, len(notes))
	for i, v := range notes {
		newSteps[i] = step.Step{
			Notes:  []note.Note{v},
			Params: paramsNew,
		}
	}

	stepString := strings.Builder{}
	for _, s := range newSteps {
		stepString.WriteString(s.String())
		stepString.WriteString(" ")
	}
	result = strings.TrimSpace(stepString.String())
	if len(result) > 0 {
		result = "(" + result + ")"
	}
	return
}
