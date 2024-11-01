package param

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/schollz/logger"
)

const PARAM_PROBABILITY = "probability"
const PARAM_TRANSPOSE = "transpose"
const PARAM_VELOCITY = "velocity"
const PARAM_GATE = "gate"
const PARAM_ARPEGGIO = "arpeggio"
const PARAM_UP = "up"
const PARAM_DOWN = "down"
const PARAM_THUMB = "thumb"
const PARAM_RANDOM = "random"
const PARAM_BEATS = "beats"
const PARAM_BPM = "bpm"

var names map[string][]string

func init() {
	// shorthand values
	names = map[string][]string{
		"probability": []string{"pr", "prob"},
		"transpose":   []string{"trans", "tr"},
		"velocity":    []string{"vel"},
		"gate":        []string{"gat"},
		"arpeggio":    []string{"arp"},
		"up":          []string{"up"},
		"down":        []string{"dow"},
		"thumb":       []string{"thumb"},
		"random":      []string{"rand"},
		"beats":       []string{"beats"},
		"bpm":         []string{"bpm"},
		"attack":      []string{"atk"},
		"decay":       []string{"dec"},
		"sustain":     []string{"sus"},
		"release":     []string{"rel"},
	}
}

type Param struct {
	TextOriginal string
	Name         string
	Values       []int
	Iterator     int
}

func New(name string, values []int) Param {
	return Param{
		Name:     name,
		Values:   values,
		Iterator: -1,
	}
}

func (p Param) String() string {
	sb := strings.Builder{}
	sb.WriteString(p.Name)
	for i, v := range p.Values {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func (p *Param) Next() int {
	p.Iterator++
	value := p.Current()
	return value
}

func (p *Param) Current() int {
	if p.Iterator < 0 {
		return p.Values[0]
	}
	return p.Values[p.Iterator%len(p.Values)]
}

func (p *Param) Rotate() {
	p.Values = append(p.Values[1:], p.Values[0])
}

func (p Param) Copy() Param {
	values := make([]int, len(p.Values))
	copy(values, p.Values)
	return Param{
		Name:     p.Name,
		Values:   values,
		Iterator: p.Iterator,
	}
}

func IsSpecialParameter(name string) bool {
	_, ok := names[name]
	return ok
}

func Parse(s string) (p Param, err error) {
	p.TextOriginal = s

	// extract all whole numbers (positive or negative) using regex
	re := regexp.MustCompile(`-?\d+`)
	valueStrings := re.FindAllString(s, -1)
	if len(valueStrings) == 0 {
		valueStrings = []string{"1"}
	}
	values := make([]int, len(valueStrings))
	for i, v := range valueStrings {
		values[i], err = strconv.Atoi(v)
		if err != nil {
			log.Errorf("could not convert %s to int: %v", v, err)
			return
		}
	}

	// find which name has the longest prefix
	longestPrefix := 0
	bestK := ""
	for k, vs := range names {
		vs = append(vs, k)
		for _, v := range vs {
			if strings.HasPrefix(s, v) && len(v) > longestPrefix {
				longestPrefix = len(v)
				bestK = k
			}
		}
	}
	if bestK == "" {
		// capture everything before the first number
		re = regexp.MustCompile(`\D+`)
		bestK = re.FindString(s)
		if bestK == "" {
			err = fmt.Errorf("could not find parameter name in %s", s)
			return
		}
	}
	p = New(bestK, values)

	return
}
