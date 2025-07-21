package interpreter

import (
	"bufio"
	"errors"
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
	stack      Stack[int]
	dictionary map[string]ExecutableToken
}

func NewInterpreter(reader *bufio.Reader) *Interpreter {
	i := Interpreter{
		reader:     bufio.NewReader(reader),
		scanner:    nil,
		stack:      Stack[int]{},
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
			fmt.Printf("%d ", a)
		},
	}
	i.dictionary["emit"] = ExecutableToken{
		name: "emit",
		primitive: func() {
			a, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			fmt.Printf("%c\n", a)
		},
	}
	i.dictionary["cr"] = ExecutableToken{
		name: "cr",
		primitive: func() {
			fmt.Println()
		},
	}
	i.dictionary[".\""] = ExecutableToken{
		name: ".\"",
		primitive: func() {
			if i.scanner.Scan() {
				w := i.scanner.Text()
				if w[len(w)-1:] != "\"" {
					log.Fatal("invalid string termination")
				}
				fmt.Println(w[:len(w)-1])
			}
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

	// Defining words
	i.dictionary[":"] = ExecutableToken{
		name: ":",
		primitive: func() {
			name, err := i.Word()
			if err != nil {
				log.Fatal(err)
			}
			words := []string{}
			for i.scanner.Scan() {
				t := i.scanner.Text()
				if t != ";" {
					words = append(words, t)
				} else {
					break
				}
			}
			definition := strings.Join(words, " ")
			xt := ExecutableToken{
				name: name,
				primitive: func() {
					s := bufio.NewScanner(strings.NewReader(definition))
					s.Split(bufio.ScanWords)
					for s.Scan() {
						t := s.Text()
						i.Interpret(t)
					}
				},
			}
			i.dictionary[name] = xt
		},
	}

	// Comments
	i.dictionary["("] = ExecutableToken{
		name: "(",
		primitive: func() {
			// Consume the text until the end of comment ')'
			for i.scanner.Scan() {
				t := i.scanner.Text()
				if t == ")" {
					break
				}
			}
		},
	}

	return &i
}

func (i *Interpreter) Interpret(word string) {
	if xt, ok := i.dictionary[word]; ok {
		xt.primitive()
	} else {
		v, err := strconv.ParseInt(word, 10, 64)
		if err == nil {
			i.stack.Push(int(v))
		} else {
			fmt.Printf("%s ?\n", word)
		}
	}
}

func (i *Interpreter) prompt() {
	for _, v := range i.stack.items {
		fmt.Printf("%d ", v)
	}
	fmt.Print("ok> ")
}

func (i *Interpreter) Word() (string, error) {
	if i.scanner != nil && i.scanner.Scan() {
		return i.scanner.Text(), nil
	} else {
		i.prompt()
		s, err := i.reader.ReadString('\n')
		if err != nil {
			return "", errors.New("end of input")
		}
		s = strings.TrimSpace(s)
		i.scanner = bufio.NewScanner(strings.NewReader(s))
		i.scanner.Split(bufio.ScanWords)
		i.scanner.Scan()
		return i.scanner.Text(), nil
	}
}
