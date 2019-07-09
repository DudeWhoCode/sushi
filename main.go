package main

import (
	"os"

	"github.com/dudewhocode/interpreter/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
