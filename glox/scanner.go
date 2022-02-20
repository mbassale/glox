package glox

import "fmt"

type Scanner interface {
	ScanTokens() []Token
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
	default:
		s.errorReporter.Error(s.line, fmt.Sprintf("Unexpected character: '%c'", c))
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

func (s *SimpleScanner) addToken(tokenType int) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *SimpleScanner) addTokenWithLiteral(tokenType int, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, string(text), literal, s.line))
}
