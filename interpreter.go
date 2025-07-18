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
