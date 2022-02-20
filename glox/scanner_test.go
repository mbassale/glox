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
		{"empty file", "", []Token{NewToken(TOKEN_EOF, "", nil, 1)}},
		{"(){}", "(){}", []Token{
			NewToken(TOKEN_LEFT_PAREN, "(", nil, 1),
			NewToken(TOKEN_RIGHT_PAREN, ")", nil, 1),
			NewToken(TOKEN_LEFT_BRACE, "{", nil, 1),
			NewToken(TOKEN_RIGHT_BRACE, "}", nil, 1),
			NewToken(TOKEN_EOF, "", nil, 1),
		},
		},
	}
	for _, testCase := range testCases {
		errorReporter := NewConsoleErrorReporter()
		scanner := NewScanner(testCase.source, errorReporter)
		currentTokens := scanner.ScanTokens()
		assert.Equal(t, testCase.expectedTokens, currentTokens)
	}
}
