package line

import (
	"fmt"
	"math"
	"strings"

	"github.com/schollz/asdf/src/arpeggio"
	"github.com/schollz/asdf/src/multiply"
	log "github.com/schollz/logger"
)

const LEFT_GROUP = "("
const RIGHT_GROUP = ")"
const HOLD_GROUP = "-"

func Parse(line string) (result string, err error) {
	result, err = parse(line)
	if err != nil {
		log.Error(err)
		return
	}
	if strings.Contains(result, LEFT_GROUP) || strings.Contains(result, RIGHT_GROUP) {
		result, err = parse(result)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}

func parse(line string) (result string, err error) {
	line = multiply.Parse(line, multiply.Parentheses)
	line = sanitizeLine(line)

	tokens := tokenizeFromGrouping(line)
	err = doParenthesesMatch(tokens)
	if err != nil {
		log.Error(err)
		return
	}
	log.Tracef("tokens: %v", tokens)

	for i, token := range tokens {
		if token != LEFT_GROUP && token != RIGHT_GROUP {
			tokens[i], err = arpeggio.Expand(token)
			if err != nil {
				log.Error(err)
				return
			}
		}
	}
	result = tokenExpandToLine(tokens)
	return
}

func sanitizeLine(line string) string {
	for {
		if !strings.Contains(line, " (") {
			break
		}
		line = strings.Replace(line, " (", "(", -1)
	}
	return line
}

func tokenizeFromGrouping(s string) (tokens []string) {
	var i int
	for i < len(s) {
		switch s[i] {
		case '(':
			tokens = append(tokens, string(LEFT_GROUP))
			i++
		case ')':
			tokens = append(tokens, string(RIGHT_GROUP))
			i++
		case ' ':
			i++
		default:
			// Find the next space or parenthesis to capture the token
			start := i
			for i < len(s) && s[i] != ' ' && s[i] != '(' && s[i] != ')' {
				i++
			}
			tokens = append(tokens, s[start:i])
		}
	}
	return
}

func doParenthesesMatch(tokens []string) (err error) {
	// check to see whether parentheses match
	// and if not, return a string showing where the error is
	depth := 0
	lastDepthIncrease := 0
	for i, token := range tokens {
		if token == LEFT_GROUP {
			depth++
			lastDepthIncrease = i
		} else if token == RIGHT_GROUP {
			depth--
			if depth < 0 {
				sb := strings.Builder{}
				sb.WriteString("\n")
				sb.WriteString(strings.Join(tokens, " ") + "\n")
				for j := 0; j < i; j++ {
					sb.WriteString("  ")
				}
				sb.WriteString("^ parentheses do not match")
				err = fmt.Errorf(sb.String())
				return
			}
		}
	}
	if depth != 0 {
		sb := strings.Builder{}
		sb.WriteString("\n")
		sb.WriteString(strings.Join(tokens, " ") + "\n")
		for i := 0; i < lastDepthIncrease; i++ {
			sb.WriteString(" ")
		}
		sb.WriteString("^ parentheses do not match")
		err = fmt.Errorf(sb.String())
	}

	return
}

// tokenExpandToLine takes a line of notes and/or chords with decorators and expands it
// for example Cm7 (d e f g) will expand to
// Cm7 - - - d e f g
// where the decorators are attached to the enttity and the entities are
// given sustains where nessecary (Cm7 is sustained)
func tokenExpandToLine(tokens []string) (expanded string) {

	currentTokenValues := make([]float64, len(tokens))
	for i := range currentTokenValues {
		currentTokenValues[i] = 1
	}
	tokenValues := determineTokenValue(tokens, currentTokenValues, 0, len(tokens))
	// print values of non-parentheses tokens
	minValue := math.Inf(1)
	for i, token := range tokens {
		if token != LEFT_GROUP && token != RIGHT_GROUP {
			if tokenValues[i] < minValue {
				minValue = tokenValues[i]
			}
		}
	}
	// now expand it by building a string
	sb := strings.Builder{}
	for i, token := range tokens {
		if token != LEFT_GROUP && token != RIGHT_GROUP {
			repetitions := int(math.Round(tokenValues[i] / minValue))
			if strings.HasPrefix(token, ".") {
				repetitions = 1
			}
			sb.WriteString(token + " ")
			for j := 1; j < repetitions; j++ {
				sb.WriteString(HOLD_GROUP + " ")
			}
		}
	}
	expanded = strings.TrimSpace(sb.String())
	expanded = strings.Join(strings.Fields(expanded), " ")
	return
}

func determineTokenValue(tokens []string, currentTokenValues []float64, start int, stop int) (tokenValues []float64) {
	depth := 0
	tokenLocations := [][2]int{} // start, end
	tokenValues = make([]float64, len(tokens))
	copy(tokenValues, currentTokenValues)
	for i := start; i < stop; i++ {
		token := tokens[i]
		if token == LEFT_GROUP {
			if depth == 0 {
				tokenLocations = append(tokenLocations, [2]int{i, -1})
			}
			depth++
		} else if token == RIGHT_GROUP {
			depth--
			if depth == 0 {
				tokenLocations[len(tokenLocations)-1][1] = i
			}
		}
	}
	for _, loc := range tokenLocations {
		tokenValues = determineTokenValue(tokens, tokenValues, loc[0]+1, loc[1])
	}

	numEntities := countEntities(tokens[start:stop])
	valuePerEntity := 1.0 / float64(numEntities)
	for i := start; i < stop; i++ {
		tokenValues[i] = valuePerEntity * tokenValues[i]
	}
	return
}

func countEntities(tokens []string) (entities int) {
	depth := 0
	for _, token := range tokens {
		if token == LEFT_GROUP {
			if depth == 0 {
				entities++
			}
			depth++
		} else if token == RIGHT_GROUP {
			depth--
		} else if strings.HasPrefix(token, ".") {
			// do nothing
		} else if depth == 0 {
			entities++
		}
	}
	return
}
