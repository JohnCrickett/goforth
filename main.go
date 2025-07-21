package main

import (
	"bufio"
	"github.com/JohnCrickett/goforth/interpreter"
	"log"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	i := interpreter.NewInterpreter("")

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
