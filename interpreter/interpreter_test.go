package interpreter

import (
	"testing"
)

func TestInterpret(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []int
	}{
		// Numbers
		"1": {
			input:    "1\n",
			expected: []int{1},
		},
		"1 2 3": {
			input:    "1 2 3",
			expected: []int{1, 2, 3},
		},
		// Mathematical ops
		"add": {
			input:    "1 2 +",
			expected: []int{3},
		},
		"add2": {
			input:    "1 2 2 + +",
			expected: []int{5},
		},
		"multiply": {
			input:    "2 3 *",
			expected: []int{6},
		},
		"divide": {
			input:    "6 2 /",
			expected: []int{3},
		},
		"subtract": {
			input:    "2 1 -",
			expected: []int{1},
		},
		"mod": {
			input:    "3 2 mod",
			expected: []int{1},
		},
		// Stack ops
		"swap": {
			input:    "3 2 swap",
			expected: []int{2, 3},
		},
		"dup": {
			input:    "3 2 2 dup",
			expected: []int{3, 2, 2},
		},
		"over": {
			input:    "3 2 over",
			expected: []int{3, 2, 3},
		},
		"rot": {
			input:    "3 2 1 rot",
			expected: []int{2, 1, 3},
		},
		"drop": {
			input:    "3 2 1 0 drop",
			expected: []int{3, 2, 1},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			interpreter := NewInterpreter(test.input)

			for {
				w, err := interpreter.Word()
				if err != nil {
					break
				}
				interpreter.Interpret(w)
			}

			for i := range len(test.expected) {
				if interpreter.stack.items[i] != test.expected[i] {
					t.Errorf("expected %v, got %v", test.expected, interpreter.stack)
				}
			}
		})
	}
}
