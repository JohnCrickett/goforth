package interpreter

// A fun one this, using AI to generate the test
// it took several attempts to get it to generate
// against the actual methods of the struct rather
// than ones it imagined a stack would have.

import (
	"testing"
)

func TestStackPush(t *testing.T) {
	stack := Stack[int]{}

	// Test pushing an element onto the stack
	stack.Push(1)

	if len(stack.items) != 1 {
		t.Errorf("expected stack size to be 1, got %d", len(stack.items))
	}

	if stack.items[0] != 1 {
		t.Errorf("expected top element to be 1, got %v", stack.items[0])
	}

	// Test pushing another element
	stack.Push(2)

	if len(stack.items) != 2 {
		t.Errorf("expected stack size to be 2, got %d", len(stack.items))
	}

	if stack.items[1] != 2 {
		t.Errorf("expected top element to be 2, got %v", stack.items[1])
	}
}

func TestStackPop(t *testing.T) {
	stack := Stack[int]{}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	if len(stack.items) != 3 {
		t.Errorf("expected stack size to be 3, got %d", len(stack.items))
	}
	stack.Pop()
	stack.Pop()
	stack.Pop()
	if len(stack.items) != 0 {
		t.Errorf("expected stack size to be 0, got %d", len(stack.items))
	}

	stack.Pop()
	if len(stack.items) != 0 {
		t.Errorf("expected stack size to be 0, got %d", len(stack.items))
	}
}

func TestStackTop(t *testing.T) {
	stack := Stack[int]{}
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)
	if len(stack.items) != 3 {
		t.Errorf("expected stack size to be 3, got %d", len(stack.items))
	}
	v, err := stack.Top()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if v != 3 {
		t.Errorf("expected top element to be 3, got %v", v)
	}
	stack.Pop()
	stack.Pop()
	v, err = stack.Top()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if v != 1 {
		t.Errorf("expected top element to be 3, got %v", v)
	}
	stack.Pop()
	v, err = stack.Top()
	if err == nil {
		t.Errorf("expected an error")
	}
}
