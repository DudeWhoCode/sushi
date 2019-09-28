package main

import (
	"os"

	"github.com/dudewhocode/sushi/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
