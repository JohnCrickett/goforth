package interpreter

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	scanner    *bufio.Scanner
	out        io.Writer
	stack      Stack[int]
	dictionary map[string]ExecutableToken
}

func NewInterpreter(writer io.Writer, source string) *Interpreter {
	i := Interpreter{
		scanner:    bufio.NewScanner(strings.NewReader(source)),
		out:        writer,
		stack:      Stack[int]{},
		dictionary: make(map[string]ExecutableToken),
	}
	i.scanner.Split(bufio.ScanWords)

	// Quiting
	i.dictionary["bye"] = ExecutableToken{
		name: "bye",
		primitive: func() {
			os.Exit(0)
		},
	}

	// Mathematical Operations
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
			_, err = fmt.Fprintf(i.out, "%d ", a)
			if err != nil {
				log.Fatal(err)
			}
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
			_, err = fmt.Fprintf(i.out, "%c", a)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	i.dictionary["cr"] = ExecutableToken{
		name: "cr",
		primitive: func() {
			_, err := fmt.Fprintln(i.out)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	i.dictionary[".\""] = ExecutableToken{
		name: ".\"",
		primitive: func() {
			first := true
			for i.scanner.Scan() {
				if !first {
					_, err := fmt.Fprint(i.out, " ")
					if err != nil {
						log.Fatal(err)
					}
				}
				w := i.scanner.Text()
				if w[len(w)-1:] == "\"" {
					w = w[:len(w)-1]
					_, err := fmt.Fprint(i.out, w)
					if err != nil {
						log.Fatal(err)
					}
					break
				}
				_, err := fmt.Fprint(i.out, w)
				if err != nil {
					log.Fatal(err)
				}
				first = false
			}
		},
	}
	i.dictionary[".S"] = ExecutableToken{
		name: ".S",
		primitive: func() {
			_, err := fmt.Fprintf(i.out, "<%d> ", len(i.stack.items))
			if err != nil {
				log.Fatal(err)
			}

			for _, v := range i.stack.items {
				_, err := fmt.Fprintf(i.out, "%d ", v)
				if err != nil {
					log.Fatal(err)
				}
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

	// Comparrisions
	i.dictionary["="] = ExecutableToken{
		name: "=",
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
			if a == b {
				i.stack.Push(-1)
			} else {
				i.stack.Push(0)
			}
		},
	}
	i.dictionary["<"] = ExecutableToken{
		name: "<",
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
			if b < a {
				i.stack.Push(-1)
			} else {
				i.stack.Push(0)
			}
		},
	}
	i.dictionary[">"] = ExecutableToken{
		name: ">",
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
			if b > a {
				i.stack.Push(-1)
			} else {
				i.stack.Push(0)
			}
		},
	}
	i.dictionary["<>"] = ExecutableToken{
		name: "<>",
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
			if b != a {
				i.stack.Push(-1)
			} else {
				i.stack.Push(0)
			}
		},
	}

	// Boolean Operators
	i.dictionary["and"] = ExecutableToken{
		name: "and",
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
			if b == -1 && a == -1 {
				i.stack.Push(-1)
			} else {
				i.stack.Push(0)
			}
		},
	}
	i.dictionary["or"] = ExecutableToken{
		name: "or",
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
			if b == -1 || a == -1 {
				i.stack.Push(-1)
			} else {
				i.stack.Push(0)
			}
		},
	}
	i.dictionary["invert"] = ExecutableToken{
		name: "invert",
		primitive: func() {
			a, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			if a == -1 {
				i.stack.Push(0)
			} else {
				i.stack.Push(-1)
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
			_, err := fmt.Fprintf(i.out, "%s ?\n", word)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (i *Interpreter) Prompt() {
	for _, v := range i.stack.items {
		_, err := fmt.Fprintf(i.out, "%d ", v)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Print("ok> ")
}

func (i *Interpreter) Word() (string, error) {
	if i.scanner != nil && i.scanner.Scan() {
		return i.scanner.Text(), nil
	} else {
		return "", errors.New("end of input")
	}
}

func (i *Interpreter) SetScanLine(line string) {
	i.scanner = bufio.NewScanner(strings.NewReader(line))
	i.scanner.Split(bufio.ScanWords)
}
