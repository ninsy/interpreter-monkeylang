package main

import (
	"monkey/repl"
	"os"
)

// TODO: negative test cases
// TODO: handle postfix operators
// TODO: record latest lines, conjure them up with upper arrow
// TODO: rewrite to C, write own GC
// TODO: add char espacing for strings: \n \t
// TODO: type coercion
// TODO: parseInt, parseFloat, isNan impl
// TODO: support floats
// TODO: prototypes
// TODO: make reduce to accepts other types ( now only int, array) as an initial value
// TODO: reassigments
// TODO: loops

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
