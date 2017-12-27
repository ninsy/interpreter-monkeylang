package main

import (
	"monkey/repl"
	"os"
)

// TODO: negative test cases
// TODO: handle postfix operators
// TODO: record latest lines, conjure them up with upper arrow

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
