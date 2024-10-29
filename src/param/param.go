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
const PARAM_UPDOWN = "updown"
const PARAM_RANDOM = "random"

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
		Iterator: 0,
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
	value := p.Current()
	p.Iterator++
	return value
}

func (p *Param) Current() int {
	return p.Values[p.Iterator%len(p.Values)]
}

func Parse(s string) (p Param, err error) {
	p.TextOriginal = s

	// shorthand values
	names := map[string][]string{
		"probability": []string{"p", "prob"},
		"transpose":   []string{"t", "trans"},
		"velocity":    []string{"v", "vel"},
		"gate":        []string{"g"},
		"arpeggio":    []string{"arp"},
		"up":          []string{"u"},
		"down":        []string{"d"},
		"updown":      []string{"ud"},
		"random":      []string{"r"},
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

	// find the name which is prefixing the string
	for k, vs := range names {
		vs = append(vs, k)
		for _, v := range vs {
			if strings.HasPrefix(s, v) {
				p = New(k, values)
				return
			}
		}
	}

	err = fmt.Errorf("could not parse %s", s)
	return
}
