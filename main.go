package main

import (
	"fmt"
	"os"

	"github.com/dudewhocode/interpreter/repl"
)

func main() {
	fmt.Println("boa v0.0.1")
	repl.Start(os.Stdin, os.Stdout)
}
