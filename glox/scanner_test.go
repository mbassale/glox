package glox

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanTokens(t *testing.T) {
	testCases := []struct {
		name           string
		source         string
		expectedTokens []Token
	}{
		{"empty", "", []Token{NewToken(TOKEN_EOF, "", nil, 1)}},
		{"single and two-character tokens", "(){}<><=>====*,;-+", []Token{
			NewToken(TOKEN_LEFT_PAREN, "(", nil, 1),
			NewToken(TOKEN_RIGHT_PAREN, ")", nil, 1),
			NewToken(TOKEN_LEFT_BRACE, "{", nil, 1),
			NewToken(TOKEN_RIGHT_BRACE, "}", nil, 1),
			NewToken(TOKEN_LESS, "<", nil, 1),
			NewToken(TOKEN_GREATER, ">", nil, 1),
			NewToken(TOKEN_LESS_EQUAL, "<=", nil, 1),
			NewToken(TOKEN_GREATER_EQUAL, ">=", nil, 1),
			NewToken(TOKEN_EQUAL_EQUAL, "==", nil, 1),
			NewToken(TOKEN_EQUAL, "=", nil, 1),
			NewToken(TOKEN_STAR, "*", nil, 1),
			NewToken(TOKEN_COMMA, ",", nil, 1),
			NewToken(TOKEN_SEMICOLON, ";", nil, 1),
			NewToken(TOKEN_MINUS, "-", nil, 1),
			NewToken(TOKEN_PLUS, "+", nil, 1),
			NewToken(TOKEN_EOF, "", nil, 1),
		},
		},
		{"numbers", "+1.23-12345.67890*1/2", []Token{
			NewToken(TOKEN_PLUS, "+", nil, 1),
			NewToken(TOKEN_NUMBER, "1.23", 1.23, 1),
			NewToken(TOKEN_MINUS, "-", nil, 1),
			NewToken(TOKEN_NUMBER, "12345.67890", 12345.6789, 1),
			NewToken(TOKEN_STAR, "*", nil, 1),
			NewToken(TOKEN_NUMBER, "1", 1.0, 1),
			NewToken(TOKEN_SLASH, "/", nil, 1),
			NewToken(TOKEN_NUMBER, "2", 2.0, 1),
			NewToken(TOKEN_EOF, "", nil, 1),
		},
		},
	}
	for _, testCase := range testCases {
		errorReporter := NewConsoleErrorReporter()
		scanner := NewScanner(testCase.source, errorReporter)
		currentTokens := scanner.ScanTokens()
		assert.Equal(t, testCase.expectedTokens, currentTokens, "Failed %s: Source: %s", testCase.name, testCase.source)
	}
}
