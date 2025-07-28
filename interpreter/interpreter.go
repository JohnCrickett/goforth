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
	environments []*bufio.Scanner
	out          io.Writer
	stack        Stack[int]
	loopStack    Stack[int]
	dictionary   map[string]ExecutableToken
}

func NewInterpreter(writer io.Writer, source string) *Interpreter {
	i := Interpreter{
		out:        writer,
		stack:      Stack[int]{},
		dictionary: make(map[string]ExecutableToken),
	}
	i.environments = append(i.environments, bufio.NewScanner(strings.NewReader(source)))
	i.environments[0].Split(bufio.ScanWords)

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
			for i.environments[len(i.environments)-1].Scan() {
				if !first {
					_, err := fmt.Fprint(i.out, " ")
					if err != nil {
						log.Fatal(err)
					}
				}
				w := i.environments[len(i.environments)-1].Text()
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
			for i.environments[len(i.environments)-1].Scan() {
				t := i.environments[len(i.environments)-1].Text()
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
					i.environments = append(i.environments, s)
					for i.environments[len(i.environments)-1].Scan() {
						t := i.environments[len(i.environments)-1].Text()
						i.Interpret(t)
					}
					i.environments = i.environments[:len(i.environments)-1]
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
			for i.environments[len(i.environments)-1].Scan() {
				t := i.environments[len(i.environments)-1].Text()
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

	// if
	i.dictionary["if"] = ExecutableToken{
		name: "if",
		primitive: func() {
			a, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()
			if a == -1 {
				interpret := true
				// true code
				// get the next word and process it,
				for i.environments[len(i.environments)-1].Scan() {
					w := i.environments[len(i.environments)-1].Text()
					if w == "then" {
						break
					} else if w == "else" {
						interpret = false
					} else if interpret {
						i.Interpret(w)
					}
				}
			} else {
				// if there is an else
				// skip everything until the else, then interpret
				for i.environments[len(i.environments)-1].Scan() {
					w := i.environments[len(i.environments)-1].Text()
					if w == "else" {
						break
					} else if w == "then" {
						// if no else, return early
						return
					}
				}
				var foundThen bool
				for i.environments[len(i.environments)-1].Scan() && !foundThen {
					w := i.environments[len(i.environments)-1].Text()
					if w == "then" {
						foundThen = true
					} else {
						i.Interpret(w)
					}
				}
				if !foundThen {
					panic("missing 'then'")
				}
			}
		},
	}

	// do loop
	i.dictionary["do"] = ExecutableToken{
		name: "do",
		primitive: func() {
			var words []string

			// grab string to 'loop'
			for i.environments[len(i.environments)-1].Scan() {
				w := i.environments[len(i.environments)-1].Text()
				if w == "loop" {
					break
				} else {
					words = append(words, w)
				}
			}
			definition := strings.Join(words, " ")

			// get the index and limit from the data stack
			start, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()

			end, err := i.stack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Pop()

			for index := start; index < end; index++ {
				i.loopStack.Push(index)
				s := bufio.NewScanner(strings.NewReader(definition))
				s.Split(bufio.ScanWords)
				i.environments = append(i.environments, s)

				for i.environments[len(i.environments)-1].Scan() {
					t := i.environments[len(i.environments)-1].Text()
					i.Interpret(t)
				}

				i.environments = i.environments[:len(i.environments)-1]
				i.loopStack.Pop()
			}
		},
	}
	i.dictionary["i"] = ExecutableToken{
		name: "i",
		primitive: func() {
			v, err := i.loopStack.Top()
			if err != nil {
				log.Fatal(err)
			}
			i.stack.Push(v)
		},
	}

	return &i
}

func (i *Interpreter) Interpret(word string) {
	defer func() {
		if r := recover(); r != nil {
			_, err := fmt.Fprintf(i.out, "%s", r)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()

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
	if i.environments != nil && i.environments[len(i.environments)-1].Scan() {
		return i.environments[len(i.environments)-1].Text(), nil
	} else {
		return "", errors.New("end of input")
	}
}

func (i *Interpreter) SetScanLine(line string) {
	i.environments[len(i.environments)-1] = bufio.NewScanner(strings.NewReader(line))
	i.environments[len(i.environments)-1].Split(bufio.ScanWords)
}
