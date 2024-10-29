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
