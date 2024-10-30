package block

import (
	"fmt"
	"strings"

	"github.com/schollz/asdf/src/line"
	"github.com/schollz/asdf/src/step"
	log "github.com/schollz/logger"
)

type Block struct {
	Name      string
	Steps     []step.Step
	TotalTime float64
}

func Parse(block string) (b Block, err error) {
	// first expand block
	expanded, err := expand(block)
	if err != nil {
		log.Error(err)
		return
	}
	log.Tracef("expanded block: %s", expanded)
	steps := make([]step.Step, 0)
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
	totalTime := 0.0
	totalBeats := 0.0
	for _, s := range steps {
		log.Tracef("step: %+v", s.Info())
		totalTime += s.Beats * 60.0 / s.BPM
		totalBeats += s.Beats
	}
	log.Tracef("totalTime: %2.3f", totalTime)
	log.Tracef("totalBeats: %2.3f", totalBeats)

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
	// calculate time start/stop
	for i := range steps {
		steps[i].TimeStart = steps[i].BeatStart * 60.0 / steps[i].BPM
		steps[i].TimeEnd = steps[i].BeatEnd * 60.0 / steps[i].BPM
	}

	log.Trace("consolidated steps:")
	for _, s := range steps {
		fmt.Printf("step: %+v\n", s.Info())
	}

	b.Steps = steps
	b.TotalTime = totalTime

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
