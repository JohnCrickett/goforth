package main

import (
	"bufio"
	"github.com/JohnCrickett/goforth/interpreter"
	"log"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	interpreter := interpreter.NewInterpreter(reader)

	for {
		word, err := interpreter.Word()
		if err != nil {
			log.Fatal(err)
		} else {
			interpreter.Interpret(word)
		}
	}
}
