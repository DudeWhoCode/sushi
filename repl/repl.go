package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dudewhocode/interpreter/object"

	"github.com/dudewhocode/interpreter/evaluator"
	"github.com/dudewhocode/interpreter/lexer"
	"github.com/dudewhocode/interpreter/parser"
)

const (
	PROMT   = ">>> "
	WELCOME = `
██████╗  ██████╗  █████╗         ██╗   ██╗     ██████╗    ██╗
██╔══██╗██╔═══██╗██╔══██╗        ██║   ██║    ██╔═████╗  ███║
██████╔╝██║   ██║███████║        ██║   ██║    ██║██╔██║  ╚██║
██╔══██╗██║   ██║██╔══██║        ╚██╗ ██╔╝    ████╔╝██║   ██║
██████╔╝╚██████╔╝██║  ██║         ╚████╔╝     ╚██████╔╝██╗██║
╚═════╝  ╚═════╝ ╚═╝  ╚═╝          ╚═══╝       ╚═════╝ ╚═╝╚═╝
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
