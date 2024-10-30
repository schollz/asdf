package block

import (
	"strings"

	"github.com/schollz/asdf/src/line"
	"github.com/schollz/asdf/src/step"
	log "github.com/schollz/logger"
)

func Parse(block string) (steps []step.Step, err error) {
	// first expand block
	expanded, err := expand(block)
	if err != nil {
		log.Error(err)
		return
	}
	log.Tracef("expanded block: %s", expanded)
	steps = make([]step.Step, 0)
	lines := strings.Split(expanded, "\n")
	bpm := 60.0
	beatsInLine := 4

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		log.Tracef("parsing line: %s", l)
		var lineSteps []step.Step
		entitiesInLine := 0
		for _, v := range strings.Fields(l) {
			var s step.Step
			s, err = step.Parse(v)
			if err != nil {
				log.Error(err)
				return
			}
			if s.IsLegato || s.IsRest || s.IsNote {
				entitiesInLine++
			}
			if s.HasParam("beats") {
				beatsInLine = s.GetParamNext("beats", beatsInLine)
			}
			if s.HasParam("bpm") {
				bpm = float64(s.GetParamNext("bpm", int(bpm)))
			}
			s.BPM = bpm
			lineSteps = append(lineSteps, s)
		}
		log.Tracef("line %s has %d entitiesInLine with %d beats", l, entitiesInLine, beatsInLine)
		if entitiesInLine > 0 {
			for i := range lineSteps {
				lineSteps[i].Beats = float64(beatsInLine) / float64(entitiesInLine)
			}
		}
		steps = append(steps, lineSteps...)
	}
	for _, s := range steps {
		log.Tracef("step: %+v", s.Info())
	}

	// consolidate steps (removing rests and legatos)
	stepsConsolidated := []step.Step{}
	currentBeat := 0.0
	for i := 0; i < len(steps); i++ {
		if !steps[i].IsNote {
			currentBeat += steps[i].Beats
			continue
		}
		s := steps[i]
		s.BeatStart = currentBeat
		beatDuration := s.Beats
		for j := 1; j < len(steps); j++ {
			// find where it stops
			k := (i + j) % len(steps)
			if steps[k].IsNote || steps[k].IsRest {
				s.BeatEnd = currentBeat + beatDuration
				break
			}
			beatDuration += steps[k].Beats
		}
		currentBeat += s.Beats
		stepsConsolidated = append(stepsConsolidated, s)
	}

	steps = stepsConsolidated
	log.Trace("consolidated steps:")
	for _, s := range steps {
		log.Tracef("step: %+v", s.Info())
	}

	return
}

func expand(block string) (result string, err error) {
	lines := strings.Split(block, "\n")
	log.Tracef("parsing block of %d lines", len(lines))

	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		log.Tracef("parsing line: %s", l)
		var resultLine string
		resultLine, err = line.Parse(l)
		if err != nil {
			log.Error(err)
			return
		}
		result += resultLine + "\n"
	}
	return
}
