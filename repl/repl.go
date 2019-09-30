package repl

import (
	"bufio"
	"fmt"
	"io"
	"log"

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

func pushStartBlock(line []byte, stack *stack) {
	startBlocks := map[byte]bool{
		'{': true,
		'(': true,
		'[': true,
	}
	for _, b := range line {
		if _, ok := startBlocks[b]; ok {
			stack.push(b)
		}
	}
}

func popEndBlock(line []byte, stack *stack) {
	endBlocks := map[byte]bool{
		'}': true,
		')': true,
		']': true,
	}
	for _, b := range line {
		if _, ok := endBlocks[b]; ok {
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
		log.Println("current line length: ", len(currentLine))
		if len(currentLine) == 0 {
			continue
		}

		lastCh := currentLine[len(currentLine)-1]
		if isStartBlock(lastCh) {
			stack.push(lastCh)
			line = append(line, currentLine...)
			log.Println("Pushed to stack and continuing: ", string(lastCh))
			continue
		}
		if isEndBlock(lastCh) {
			// you don't care about what char is getting popped out, you leave it to the parser to handle the closing block
			ch := stack.pop()
			log.Println("Popped from stack: ", string(ch))
		}
		if stack.count != 0 {
			log.Println("Stack is not empty, continuing")
			continue
		}

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
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Error while parsing: \n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
