package main

import (
	"fmt"
	"github.com/sean-d/sloth/repl"
	"os"
	"os/user"
)

func main() {
	usr, err := user.Current()

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n\n\n", repl.WELCOME_SLOTH)
	fmt.Printf("welcom %s to sloth.0\n\n", usr.Username)

	repl.Start(os.Stdin, os.Stdout)
}
