package line

import (
	"fmt"
	"testing"

	log "github.com/schollz/logger"
)

func TestLine(t *testing.T) {
	log.SetLevel("trace")
	tests := []struct {
		line     string
		expected string
	}{
		{"a*4 b c", "a a a a b - - - c - - -"},
		{"(a a a)*2 b c", "a a a a a a b - - - - - c - - - - -"},
		{"- a*2 b c", "- - a a b - c -"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("line(%s)", test.line), func(t *testing.T) {
			result, err := Parse(test.line)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else if result != test.expected {
				t.Fatalf("\n\t%s -->\n\t%v != %v", test.line, result, test.expected)
			}
		})
	}

}
