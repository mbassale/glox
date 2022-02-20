package glox

import "fmt"

const (
	_ = iota

	// Single-character tokens.
	TOKEN_LEFT_PAREN
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE
	TOKEN_RIGHT_BRACE
	TOKEN_COMMA
	TOKEN_DOT
	TOKEN_MINUS
	TOKEN_PLUS
	TOKEN_SEMICOLON
	TOKEN_SLASH
	TOKEN_STAR

	// One or two character tokens.
	TOKEN_BANG
	TOKEN_BANG_EQUAL
	TOKEN_EQUAL
	TOKEN_EQUAL_EQUAL
	TOKEN_GREATER
	TOKEN_GREATER_EQUAL
	TOKEN_LESS
	TOKEN_LESS_EQUAL

	// Literals.
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	// Keywords.
	TOKEN_AND
	TOKEN_CLASS
	TOKEN_ELSE
	TOKEN_FALSE
	TOKEN_FUN
	TOKEN_FOR
	TOKEN_IF
	TOKEN_NIL
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_SUPER
	TOKEN_THIS
	TOKEN_TRUE
	TOKEN_VAR
	TOKEN_WHILE

	TOKEN_EOF
)

type Token struct {
	Type    int
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(tokenType int, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%v %v %v", TokenTypeToString(t.Type), t.Lexeme, TokenLiteralToString(t.Type, t.Literal))
}

func TokenTypeToString(tokenType int) string {
	switch tokenType {
	case TOKEN_LEFT_PAREN:
		return "("
	case TOKEN_RIGHT_PAREN:
		return ")"
	case TOKEN_LEFT_BRACE:
		return "{"
	case TOKEN_RIGHT_BRACE:
		return "}"
	case TOKEN_COMMA:
		return ","
	case TOKEN_MINUS:
		return "-"
	case TOKEN_PLUS:
		return "+"
	case TOKEN_SEMICOLON:
		return ";"
	case TOKEN_STAR:
		return "*"
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_BANG_EQUAL:
		return "!="
	case TOKEN_BANG:
		return "!"
	case TOKEN_EQUAL_EQUAL:
		return "=="
	case TOKEN_EQUAL:
		return "="
	case TOKEN_LESS_EQUAL:
		return "<="
	case TOKEN_LESS:
		return "<"
	case TOKEN_GREATER_EQUAL:
		return ">="
	case TOKEN_GREATER:
		return ">"
	case TOKEN_SLASH:
		return "/"
	case TOKEN_STRING:
		return "TOKEN_STRING"
	default:
		return "N/A"
	}
}

func TokenLiteralToString(tokenType int, literal interface{}) string {
	switch literal := literal.(type) {
	case string:
		return literal
	case []rune:
		return string(literal)
	default:
		return fmt.Sprintf("%v", literal)
	}
}
