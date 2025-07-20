package main

import (
	"bufio"
	"github.com/JohnCrickett/goforth/interpreter"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	interpreter := interpreter.NewInterpreter(reader)

	for {
		word := interpreter.Word()
		interpreter.Interpret(word)
	}
}
