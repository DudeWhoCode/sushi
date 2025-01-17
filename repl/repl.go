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

var startBlocks = map[byte]bool{
	'{': true,
	'(': true,
	'[': true,
}

var endBlocks = map[byte]bool{
	'}': true,
	')': true,
	']': true,
}

func pushPopBlocks(line []byte, stack *stack) {
	for _, ch := range line {
		if _, ok := startBlocks[ch]; ok {
			stack.push(ch)
		}
		if _, ok := endBlocks[ch]; ok {
			stack.pop()
		}
	}
}

func Start(in io.Reader, out io.Writer) {
	stack := NewStack()
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	io.WriteString(out, WELCOME)
	io.WriteString(out, "\n")
	var line []byte
	for {
		fmt.Printf(PROMT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		currentLine := scanner.Bytes()
		pushPopBlocks(currentLine, stack)

		if stack.count != 0 {
			line = append(line, currentLine...)
			continue
		}
		line = append(line, currentLine...)

		// Interpreter creates a new lexer for every new line
		l := lexer.New(string(line))
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
		line = []byte{}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error while parsing: \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
