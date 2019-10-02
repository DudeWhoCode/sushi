package repl

import "testing"

func TestPushPopBlocks(t *testing.T) {
	tests := []struct {
		name        string
		line        []byte
		s           *stack
		expectedCnt int
	}{
		{"OpeningBlocks",
			[]byte(
				"{(["),
			NewStack(),
			3,
		},
		{"ClosingBlocks",
			[]byte(
				"})]"),
			NewStack(),
			0,
		},
		{"LetStatement",
			[]byte(
				`let a = [
				1, 2, 3
				]`),
			NewStack(),
			0,
		},
		{"Function",
			[]byte(
				`let a = fn(
					x, y) {
						return x + y
					} `),
			NewStack(),
			0,
		},
		{"Array",
			[]byte(
				`let a = [
					1, 2, 
					3] `),
			NewStack(),
			0,
		},
		{"ArrayIncomplete",
			[]byte(
				`let a = [
					1, 2, 3 `),
			NewStack(),
			1,
		},
		{"FunctionIncomplete",
			[]byte(
				`let a = fn(
					x, y {
						return x + y
					} `),
			NewStack(),
			1,
		},
		{"FunctionWrongBraces",
			[]byte(
				`let a = fn(
					x, y} {
						return x + y
					) `),
			NewStack(),
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pushPopBlocks(tt.line, tt.s)
			if tt.s.count != tt.expectedCnt {
				t.Errorf("want stack count = %d, got = %d", tt.expectedCnt, tt.s.count)
			}
		})
	}
}
