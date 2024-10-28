package multiply

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		line     string
		expected string
	}{
		{"a b c", "a b c"},
		{"a*3 b c", "[a a a] b c"},
		{"a*3 b*2 c", "[a a a] [b b] c"},
		{"[a b] c", "[a b] c"},
		{"[a b]*3 c", "[[a b] [a b] [a b]] c"},
		{"[[a b] * 2]*2 c", "[[[[a b] [a b]]] [[[a b] [a b]]]] c"},
		{"a  *3 b c", "[a a a] b c"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("line(%s)", test.line), func(t *testing.T) {
			result := Parse(test.line, Brackets)
			if result != test.expected {
				t.Fatalf("\n\t%s -->\n\t%v != %v", test.line, result, test.expected)
			}
		})
	}

	tests = []struct {
		line     string
		expected string
	}{
		{"a.123 b c", "a.123 b c"},
		{"a *3 b c", "(a a a) b c"},
		{"a*3 b * 2 c", "(a a a) (b b) c"},
		{"(a b) c", "(a b) c"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("line(%s)", test.line), func(t *testing.T) {
			result := Parse(test.line, Parentheses)
			if result != test.expected {
				t.Fatalf("\n\t%s -->\n\t%v != %v", test.line, result, test.expected)
			}
		})
	}
}
