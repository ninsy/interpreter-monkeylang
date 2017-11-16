package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

// TODO: negative test cases
// TODO: handle postfix operators

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands.\n")
	repl.Start(os.Stdin, os.Stdout)
}
