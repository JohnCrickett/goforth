package interpreter

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestNumbers(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []int
	}{
		"1": {
			input:    "1\n",
			expected: []int{1},
		},
		"1 2 3": {
			input:    "1 2 3",
			expected: []int{1, 2, 3},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			interpreter := NewInterpreter(bufio.NewReader(os.Stdin))
			interpreter.scanner = bufio.NewScanner(strings.NewReader(test.input))
			interpreter.scanner.Split(bufio.ScanWords)
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

func TestMathOperations(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []int
	}{
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			interpreter := NewInterpreter(bufio.NewReader(os.Stdin))
			interpreter.scanner = bufio.NewScanner(strings.NewReader(test.input))
			interpreter.scanner.Split(bufio.ScanWords)
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
