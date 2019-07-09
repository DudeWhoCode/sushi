package main

import (
	"os"

	"github.com/dudewhocode/boa/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
