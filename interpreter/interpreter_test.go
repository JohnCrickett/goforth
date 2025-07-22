package interpreter

import (
	"strings"
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
			input:    "3 2 dup",
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

		// Comments
		"comments": {
			input:    "3 2 1 0 ( -1 -2 ) drop",
			expected: []int{3, 2, 1},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var o strings.Builder
			interpreter := NewInterpreter(&o, test.input)

			for {
				w, err := interpreter.Word()
				if err != nil {
					break
				}
				interpreter.Interpret(w)
			}

			ValidateStack(t, interpreter.stack, test.expected)
		})
	}
}

func TestOutput(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedOutput string
		expectedStack  []int
	}{
		// Output
		".": {
			input:          "1 2 . .",
			expectedOutput: "2 1 ",
			expectedStack:  []int{},
		},
		"emit": {
			input:          "67 74 emit emit",
			expectedOutput: "JC",
			expectedStack:  []int{},
		},
		"cr": {
			input:          "1 2 . cr . cr",
			expectedOutput: "2 \n1 \n",
			expectedStack:  []int{},
		},
		".\"no space": {
			input:          ".\" hello\"",
			expectedOutput: "hello",
			expectedStack:  []int{},
		},
		".\" with space": {
			input:          ".\" hello world\"",
			expectedOutput: "hello world",
			expectedStack:  []int{},
		},
		".S": {
			input:          "1 2 3 4 .S",
			expectedOutput: "<4> 1 2 3 4 ",
			expectedStack:  []int{1, 2, 3, 4},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var o strings.Builder
			interpreter := NewInterpreter(&o, test.input)

			for {
				w, err := interpreter.Word()
				if err != nil {
					break
				}
				interpreter.Interpret(w)
			}

			if o.String() != test.expectedOutput {
				t.Errorf("expected '%v', got '%v'", test.expectedOutput, o.String())
			}

			ValidateStack(t, interpreter.stack, test.expectedStack)
		})
	}
}

func TestDefiningAndRunningWords(t *testing.T) {
	tests := map[string]struct {
		input          string
		expectedOutput string
		expectedStack  []int
	}{
		"add word": {
			input:          "add\n: add 1 + ; 1 add .",
			expectedOutput: "add ?\n2 ",
			expectedStack:  []int{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var o strings.Builder
			interpreter := NewInterpreter(&o, test.input)

			for {
				w, err := interpreter.Word()
				if err != nil {
					break
				}
				interpreter.Interpret(w)
			}

			if test.expectedOutput != o.String() {
				t.Errorf("expected '%v', got '%v'", test.expectedOutput, o.String())
			}

			ValidateStack(t, interpreter.stack, test.expectedStack)
		})
	}
}

func ValidateStack(t *testing.T, stack Stack[int], expected []int) {
	t.Helper()

	if len(stack.items) != len(expected) {
		t.Errorf("expected %v, got %v", expected, stack)
	} else {

		for i := range len(expected) {
			if stack.items[i] != expected[i] {
				t.Errorf("expected %v, got %v", expected, stack)
			}
		}
	}
}
