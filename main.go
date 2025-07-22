package main

import (
	"bufio"
	"flag"
	"github.com/JohnCrickett/goforth/interpreter"
	"log"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	filenames := flag.Args()
	if len(filenames) > 1 {
		log.Fatal("only one file can be specified")
	}
	if len(filenames) == 1 {
		sb, err := os.ReadFile(filenames[0])
		if err != nil {
			log.Fatal(err)
		}
		source := string(sb)
		i := interpreter.NewInterpreter(os.Stdout, source)
		for {
			w, err := i.Word()
			if err != nil {
				break
			}
			i.Interpret(w)
		}

	} else {
		reader := bufio.NewReader(os.Stdin)
		i := interpreter.NewInterpreter(os.Stdout, "")

		for {
			word, err := i.Word()
			if err != nil {
				i.Prompt()
				s, err := reader.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
				s = strings.TrimSpace(s)
				i.SetScanLine(s)
			} else {
				i.Interpret(word)
			}
		}
	}
}
