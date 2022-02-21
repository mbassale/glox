package glox

import (
	"fmt"
	"strconv"
)

type Scanner interface {
	ScanTokens() []Token
}

var Keywords = map[string]int{
	"and":    TOKEN_AND,
	"class":  TOKEN_CLASS,
	"else":   TOKEN_ELSE,
	"false":  TOKEN_FALSE,
	"for":    TOKEN_FOR,
	"fun":    TOKEN_FUN,
	"if":     TOKEN_IF,
	"nil":    TOKEN_NIL,
	"or":     TOKEN_OR,
	"print":  TOKEN_PRINT,
	"return": TOKEN_RETURN,
	"super":  TOKEN_SUPER,
	"this":   TOKEN_THIS,
	"true":   TOKEN_TRUE,
	"var":    TOKEN_VAR,
	"while":  TOKEN_WHILE,
}

type SimpleScanner struct {
	source        []rune
	errorReporter ErrorReporter
	tokens        []Token
	start         int
	current       int
	line          int
}

func NewScanner(source string, errorReporter ErrorReporter) *SimpleScanner {
	return &SimpleScanner{
		source:        []rune(source),
		errorReporter: errorReporter,
		tokens:        []Token{},
		start:         0,
		current:       0,
		line:          1,
	}
}

func (s *SimpleScanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, NewToken(TOKEN_EOF, "", nil, s.line))
	return s.tokens
}

func (s *SimpleScanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(TOKEN_LEFT_PAREN)
	case ')':
		s.addToken(TOKEN_RIGHT_PAREN)
	case '{':
		s.addToken(TOKEN_LEFT_BRACE)
	case '}':
		s.addToken(TOKEN_RIGHT_BRACE)
	case ',':
		s.addToken(TOKEN_COMMA)
	case '-':
		s.addToken(TOKEN_MINUS)
	case '+':
		s.addToken(TOKEN_PLUS)
	case ';':
		s.addToken(TOKEN_SEMICOLON)
	case '*':
		s.addToken(TOKEN_STAR)
	case '!':
		if s.match('=') {
			s.addToken(TOKEN_BANG_EQUAL)
		} else {
			s.addToken(TOKEN_BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(TOKEN_EQUAL_EQUAL)
		} else {
			s.addToken(TOKEN_EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(TOKEN_LESS_EQUAL)
		} else {
			s.addToken(TOKEN_LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(TOKEN_GREATER_EQUAL)
		} else {
			s.addToken(TOKEN_GREATER)
		}
	case '/':
		if s.match('/') {
			// comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TOKEN_SLASH)
		}
	case ' ':
	case '\r':
	case '\t':
		// ignore whitespace
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.errorReporter.Error(s.line, fmt.Sprintf("Unexpected character: '%c'", c))
		}
	}
}

func (s *SimpleScanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *SimpleScanner) advance() rune {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *SimpleScanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *SimpleScanner) peek() rune {
	if s.isAtEnd() {
		return 0x0
	}
	return s.source[s.current]
}

func (s *SimpleScanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0x0
	}
	return s.source[s.current+1]
}

func (s *SimpleScanner) addToken(tokenType int) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *SimpleScanner) addTokenWithLiteral(tokenType int, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, string(text), literal, s.line))
}

func (s *SimpleScanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errorReporter.Error(s.line, "Unterminated string.")
		return
	}

	// the closing ".
	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(TOKEN_STRING, value)
}

func (s *SimpleScanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	numberStr := string(s.source[s.start:s.current])
	if number, err := strconv.ParseFloat(numberStr, 64); err == nil {
		s.addTokenWithLiteral(TOKEN_NUMBER, number)
	} else {
		s.errorReporter.Error(s.line, fmt.Sprintf("Invalid float number: %v", numberStr))
	}
}

func (s *SimpleScanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := string(s.source[s.start:s.current])
	tokenType, isReservedWord := Keywords[text]
	if isReservedWord {
		s.addToken(tokenType)
	} else {
		s.addTokenWithLiteral(TOKEN_IDENTIFIER, text)
	}
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}
