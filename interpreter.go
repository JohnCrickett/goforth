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
	scanner    *bufio.Scanner
	stack      Stack
	dictionary map[string]ExecutableToken
}

func NewInterpreter(reader *bufio.Reader) *Interpreter {
	i := Interpreter{
		reader:     bufio.NewReader(reader),
		scanner:    nil,
		stack:      Stack{},
		dictionary: make(map[string]ExecutableToken),
	}
	// Quiting
	i.dictionary["bye"] = ExecutableToken{
		name: "bye",
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

	// Mathematical Operations
	i.dictionary["-"] = ExecutableToken{
		name: "-",
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
			i.stack.Push(b - a)
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
	i.dictionary["/"] = ExecutableToken{
		name: "/",
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
			i.stack.Push(b / a)
		},
	}
	i.dictionary["mod"] = ExecutableToken{
		name: "mod",
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
			i.stack.Push(b % a)
		},
	}

	// Stack manipulation
	i.dictionary["swap"] = ExecutableToken{
		name: "swap",
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
			i.stack.Push(a)
			i.stack.Push(b)
		},
	}
	i.dictionary["dup"] = ExecutableToken{
		name: "dup",
		primitive: func() {
			a, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Push(a)
		},
	}
	i.dictionary["over"] = ExecutableToken{
		name: "over",
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
			i.stack.Push(a)
			i.stack.Push(b)
		},
	}
	i.dictionary["rot"] = ExecutableToken{
		name: "rot",
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
			c, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			i.stack.Push(b)
			i.stack.Push(a)
			i.stack.Push(c)
		},
	}
	i.dictionary["drop"] = ExecutableToken{
		name: "drop",
		primitive: func() {
			i.stack.Pop()
		},
	}

	// Output
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
	i.dictionary[".S"] = ExecutableToken{
		name: ".S",
		primitive: func() {
			fmt.Printf("<%d> ", len(i.stack.items))
			for _, v := range i.stack.items {
				fmt.Printf("%d ", v)
			}
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
			log.Fatal("unknown word ", word)
		}
	}
}

func (i *Interpreter) prompt() {
	for _, v := range i.stack.items {
		fmt.Printf("%d ", v)
	}
	fmt.Print("ok> ")
}

func (i *Interpreter) word() string {
	if i.scanner != nil && i.scanner.Scan() {
		return i.scanner.Text()
	} else {
		i.prompt()
		s, _ := i.reader.ReadString('\n')
		s = strings.TrimSpace(s)
		i.scanner = bufio.NewScanner(strings.NewReader(s))
		i.scanner.Split(bufio.ScanWords)
		i.scanner.Scan()
		return i.scanner.Text()
	}
}
