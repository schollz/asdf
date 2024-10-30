package fileparser

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/schollz/asdf/src/block"
	"github.com/schollz/asdf/src/emitter"
	"github.com/schollz/asdf/src/multiply"
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
	blockIn := false
	blockRead := false
	blockText := ""
	blockName := ""
	outputIn := false
	outputRead := false
	outputText := ""
	outputName := ""
	var currentBlock block.Block
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "**") {
			blockName = strings.TrimPrefix(strings.TrimSuffix(line, "**"), "**")
			blockIn = true
		} else if strings.HasPrefix(line, "*") {
			outputName = strings.TrimPrefix(strings.TrimSuffix(line, "*"), "*")
			outputIn = true
		} else if strings.HasPrefix(line, "```") {
			if blockIn {
				if blockRead {
					currentBlock, err = block.Parse(blockText)
					if err != nil {
						log.Error(err)
						return
					}
					currentBlock.Name = blockName
					sequences.Blocks = append(sequences.Blocks, currentBlock)
					currentBlock = block.Block{}
					blockIn = false
					blockText = ""
				}
				blockRead = !blockRead
			} else if outputIn {
				if outputRead {
					log.Tracef("output: %s", outputName)
					outputText = strings.Join(strings.Split(outputText, "\n"), " ")
					outputText = multiply.Parse(outputText, multiply.Parentheses)
					// remove all parentheses
					outputText = strings.ReplaceAll(outputText, "(", "")
					outputText = strings.ReplaceAll(outputText, ")", "")
					fmt.Println(outputText)

					// copy first block
					out := sprocket.Sprocket{Name: outputName}
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
					for _, name := range strings.Fields(outputName) {
						dotFields := strings.Split(name, ".")
						if dotFields[0] == "midi" && len(dotFields) > 1 {
							channel := 0
							if len(dotFields) > 2 && strings.HasPrefix(dotFields[2], "ch") {
								channel, _ = strconv.Atoi(strings.TrimPrefix(dotFields[2], "ch"))
							}
							midiEmitter, errMidi := emitter.NewMidi(dotFields[1], channel)
							if errMidi != nil {
								log.Error(errMidi)
								continue
							}
							emitters = append(emitters, midiEmitter)
						} else if dotFields[0] == "debug" {
							emitters = append(emitters, emitter.Debugger{})
						} else {
							log.Error(fmt.Errorf("could not find emitter %s", name))
						}
					}
					out.Player = player.New(emitters)

					// add output to list
					sequences.Sprockets = append(sequences.Sprockets, out)

					outputIn = false
					outputText = ""
				}
				outputRead = !outputRead
			}
		} else if blockRead {
			blockText += line + "\n"
		} else if outputRead {
			outputText += line + "\n"
		}
	}

	return
}
