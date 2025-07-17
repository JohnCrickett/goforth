package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type ExecutableToken struct {
	name      string
	primitive func()
}

type Interpreter struct {
	reader     *bufio.Reader
	stack      Stack
	dictionary map[string]ExecutableToken
}

func NewInterpreter(reader *bufio.Reader) *Interpreter {
	i := Interpreter{
		reader:     bufio.NewReader(reader),
		stack:      Stack{},
		dictionary: make(map[string]ExecutableToken),
	}

	i.dictionary["exit"] = ExecutableToken{
		name: "exit",
		primitive: func() {
			os.Exit(0)
		},
	}
	i.dictionary["+"] = ExecutableToken{
		name: "+",
		primitive: func() {
			a, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			b, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			i.stack.Push(a + b)
		},
	}
	i.dictionary["*"] = ExecutableToken{
		name: "*",
		primitive: func() {
			a, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			b, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			i.stack.Push(a * b)
		},
	}
	i.dictionary["drop"] = ExecutableToken{
		name: "drop",
		primitive: func() {
			i.stack.Pop()
		},
	}
	i.dictionary["."] = ExecutableToken{
		name: ".",
		primitive: func() {
			a, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			fmt.Printf("%d\n", a)
		},
	}

	return &i
}

func (i *Interpreter) interpret(word string) {
	if xt, ok := i.dictionary[word]; ok {
		xt.primitive()
	} else {
		v, err := strconv.ParseInt(word, 10, 64)
		if err == nil {
			i.stack.Push(int(v))
		} else {
			log.Fatal("unknown word")
		}
	}
}

func (i *Interpreter) word() string {
	word, _ := i.reader.ReadString('\n')
	return strings.TrimSpace(word)
}

func (i *Interpreter) prompt() {
	fmt.Printf("%v ok> ", i.stack.items)
}
