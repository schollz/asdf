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

func Parse(s string) (p Param, err error) {
	p.TextOriginal = s

	// shorthand values
	names := map[string][]string{
		"probability": []string{"p", "prob"},
		"transpose":   []string{"trans", "t", "tr"},
		"velocity":    []string{"v", "vel"},
		"gate":        []string{"g"},
		"arpeggio":    []string{"arp"},
		"up":          []string{"u"},
		"down":        []string{"d"},
		"thumb":       []string{"thumb"},
		"random":      []string{"r"},
		"beats":       []string{"beats"},
		"bpm":         []string{"bpm"},
	}
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
		err = fmt.Errorf("could not parse %s", s)
		return
	}
	p = New(bestK, values)

	return
}
