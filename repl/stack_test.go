package repl

import "testing"

func TestStringStack(t *testing.T) {
	s := NewStack()
	tests := []struct {
		input byte
	}{
		{'{'},
		{'('},
		{')'},
		{'}'},
	}

	for _, tt := range tests {
		s.push(tt.input)
	}
	for _, tt := range tests {
		expected := tt.input
		item := s.pop()
		if item != expected {
			t.Errorf("expected=%c. got=%c", expected, item)
		}
	}
}
