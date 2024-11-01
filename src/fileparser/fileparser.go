package fileparser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/asdf/src/block"
	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/jsparse"
	"github.com/schollz/asdf/src/multiply"
	"github.com/schollz/asdf/src/noteorchord"
	"github.com/schollz/asdf/src/player"
	"github.com/schollz/asdf/src/sprocket"
	log "github.com/schollz/logger"
)

type Sequences struct {
	Blocks    []block.Block
	Sprockets []sprocket.Sprocket
}

func (s Sequences) GetBlock(name string) (b block.Block, err error) {
	for _, block := range s.Blocks {
		if block.Name == name {
			b = block
			return
		}
	}
	err = fmt.Errorf("could not find block %s", name)
	return
}

func Parse(filename string) (sequences Sequences, err error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return
	}
	variables, values, err := jsparse.Parse(string(b))
	if err != nil {
		log.Error(err)
	}
	for i, v := range variables {
		// check if block or if output
		lines := strings.Split(values[i], "\n")
		// try to parse first element
		el := strings.Fields(strings.Replace(lines[0], "(", "", -1))[0]
		_, errNote := noteorchord.Parse(strings.Split(el, ".")[0])
		if errNote == nil || strings.HasPrefix(el, ".") {
			// is block
			var currentBlock block.Block
			currentBlock, err = block.Parse(values[i])
			if err != nil {
				log.Error(err)
				return
			}
			currentBlock.Name = v
			sequences.Blocks = append(sequences.Blocks, currentBlock)
		} else {
			// is output
			outputText := strings.Join(strings.Split(values[i], "\n"), " ")
			outputText = multiply.Parse(outputText, multiply.Parentheses)
			// remove all parentheses
			outputText = strings.ReplaceAll(outputText, "(", "")
			outputText = strings.ReplaceAll(outputText, ")", "")

			// copy first block
			out := sprocket.Sprocket{Name: v}
			for i, name := range strings.Fields(outputText) {
				var bl block.Block
				bl, err = sequences.GetBlock(name)
				if err != nil {
					log.Error(err)
					return
				}
				if i == 0 {
					out.Block = bl
				} else {
					out.Block.Add(bl)
				}
			}

			// compute emitter + player
			emitters := []emitter.Emitter{}
			for _, name := range strings.Fields(v) {
				dashFields := strings.Split(name, "_")
				if dashFields[0] == "midi" && len(dashFields) > 1 {
					channel := 0
					if len(dashFields) > 2 && strings.HasPrefix(dashFields[2], "ch") {
						channel, _ = strconv.Atoi(strings.TrimPrefix(dashFields[2], "ch"))
					}
					midiEmitter, errMidi := emitter.NewMidi(dashFields[1], channel)
					if errMidi != nil {
						log.Error(errMidi)
						continue
					}
					emitters = append(emitters, &midiEmitter)
				} else if dashFields[0] == "crow" && len(dashFields) > 1 {
					pitch := 0
					env := 0
					for _, field := range dashFields[1:] {
						if strings.HasPrefix(field, "voct") {
							pitch, _ = strconv.Atoi(strings.TrimPrefix(field, "voct"))
						} else if strings.HasPrefix(field, "env") {
							env, _ = strconv.Atoi(strings.TrimPrefix(field, "env"))
						}
					}
					crowEmitter, errCrow := emitter.NewCrow(pitch, env)
					if errCrow != nil {
						log.Error(errCrow)
						continue
					}
					emitters = append(emitters, &crowEmitter)
				} else if dashFields[0] == "debug" {
					emitters = append(emitters, &emitter.Debugger{})
				} else {
					log.Error(fmt.Errorf("could not find emitter %s", name))
				}
			}
			out.Player = player.New(emitters)

			// add output to list
			sequences.Sprockets = append(sequences.Sprockets, out)

		}
	}

	return
}
