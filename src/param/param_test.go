package param

import "testing"

func TestParam(t *testing.T) {
	p := New("test", []int{1, 2, 3})
	if p.Next() != 1 {
		t.Error("Expected 1")
	}
	if p.Next() != 2 {
		t.Error("Expected 2")
	}
	if p.Next() != 3 {
		t.Error("Expected 3")
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		s        string
		expected Param
	}{
		{"p50", Param{Name: "probability", Values: []int{50}}},
		{"p50,20", Param{Name: "probability", Values: []int{50, 20}}},
		{"vel120,60,30", Param{Name: "velocity", Values: []int{120, 60, 30}}},
	}
	for _, test := range tests {
		p, err := Parse(test.s)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if p.Name != test.expected.Name {
			t.Errorf("expected %s but got %s", test.expected.Name, p.Name)
		}
		if len(p.Values) != len(test.expected.Values) {
			t.Errorf("expected %v but got %v", test.expected.Values, p.Values)
		}
	}
}
