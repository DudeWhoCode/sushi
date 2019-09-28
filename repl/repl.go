package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dudewhocode/sushi/object"

	"github.com/dudewhocode/sushi/evaluator"
	"github.com/dudewhocode/sushi/lexer"
	"github.com/dudewhocode/sushi/parser"
)

const (
	PROMT   = ">>> "
	WELCOME = `
	すし - v 0.1
`
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	io.WriteString(out, WELCOME)
	io.WriteString(out, "\n")
	for {
		fmt.Printf(PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		// Interpreter creates a new lexer for every new line
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error while parsing: \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
