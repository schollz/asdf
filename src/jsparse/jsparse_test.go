package jsparse

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		code         string
		expectedvars []string
		expectedvals []string
	}{
		{"var a = 'hello'; var b = 'world';", []string{"a", "b"}, []string{"hello", "world"}},
		{"var a = 'hello'; var b = 'world'; var c = 'goodbye';", []string{"a", "b", "c"}, []string{"hello", "world", "goodbye"}},
		{"var a = 'hello';", []string{"a"}, []string{"hello"}},
		{`
x=3;		
bar1 = ` + "`" + `
c d e f${x}		
` + "`;", []string{"bar1"}, []string{"c d e f3"}},
	}
	for _, test := range tests {
		variables, values, err := Parse(test.code)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(variables, test.expectedvars) {
			t.Errorf("expected %v but got %v", test.expectedvars, variables)
		}
		if !reflect.DeepEqual(values, test.expectedvals) {
			t.Errorf("expected %v but got %v", test.expectedvals, values)
		}

	}
}
