package main

import (
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	interpreter := NewInterpreter(reader)

	for {
		interpreter.prompt()
		word := interpreter.word()
		interpreter.interpret(word)
	}
}
