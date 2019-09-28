package lexer

import (
	"github.com/dudewhocode/sushi/token"
)

type Lexer struct {
	input        string
	position     int  // current char
	readPosition int  // after current char
	ch           byte // current charecter under analysis
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

func (l *Lexer) readNumber() (string, bool) {
	startPosition := l.position
	var isFloat bool
	for isDigit(l.ch) {
		l.readChar()
	}
	if isPeriod(l.ch) {
		l.readChar() // read the floating point .
		for isDigit(l.ch) {
			l.readChar()
		}
		isFloat = true
	}
	return l.input[startPosition:l.position], isFloat
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isPeriod(ch byte) bool {
	return ch == '.'
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	startPosition := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[startPosition:l.position]
}

func (l *Lexer) NextToken() *token.Token {
	var tok *token.Token
	l.skipWhiteSpace()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.NewToken(token.EQ, literal)
		} else {
			tok = token.NewToken(token.ASSIGN, string(l.ch))
		}
	case ';':
		tok = token.NewToken(token.SEMICOLON, string(l.ch))
	case '(':
		tok = token.NewToken(token.LPAREN, string(l.ch))
	case ')':
		tok = token.NewToken(token.RPAREN, string(l.ch))
	case ',':
		tok = token.NewToken(token.COMMA, string(l.ch))
	case '+':
		tok = token.NewToken(token.PLUS, string(l.ch))
	case '-':
		tok = token.NewToken(token.MINUS, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.NewToken(token.NOTEQ, literal)
		} else {
			tok = token.NewToken(token.BANG, string(l.ch))
		}
	case '*':
		tok = token.NewToken(token.ASTERISK, string(l.ch))
	case '/':
		tok = token.NewToken(token.SLASH, string(l.ch))
	case '<':
		tok = token.NewToken(token.LT, string(l.ch))
	case '>':
		tok = token.NewToken(token.GT, string(l.ch))
	case '{':
		tok = token.NewToken(token.LBRACE, string(l.ch))
	case '}':
		tok = token.NewToken(token.RBRACE, string(l.ch))
	case 0:
		tok = token.NewToken(token.EOF, "")
	case '"':
		tok = token.NewToken(token.STRING, "")
		tok.Literal = l.readString()
	case '[':
		tok = token.NewToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = token.NewToken(token.RBRACKET, string(l.ch))
	case ':':
		tok = token.NewToken(token.COLON, string(l.ch))
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokenType := token.LookupIdent(literal)
			tok = token.NewToken(tokenType, literal)
		} else if isDigit(l.ch) {
			literal, isFloat := l.readNumber()
			if isFloat {
				tok = token.NewToken(token.FLOAT, literal)
			} else {
				tok = token.NewToken(token.INT, literal)
			}
		} else {
			tok = token.NewToken(token.ILLEGAL, string(l.ch))
		}
		return tok // If you dont return the l.readChar will be called and it advances the index once more
	}
	l.readChar()
	return tok
}
