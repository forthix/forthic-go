package forthic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizerBasicTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "single word",
			input:    "WORD",
			expected: []TokenType{TOKEN_WORD, TOKEN_EOS},
		},
		{
			name:     "multiple words",
			input:    "WORD1 WORD2 WORD3",
			expected: []TokenType{TOKEN_WORD, TOKEN_WORD, TOKEN_WORD, TOKEN_EOS},
		},
		{
			name:     "array tokens",
			input:    "[ 1 2 3 ]",
			expected: []TokenType{TOKEN_START_ARRAY, TOKEN_WORD, TOKEN_WORD, TOKEN_WORD, TOKEN_END_ARRAY, TOKEN_EOS},
		},
		{
			name:     "module tokens",
			input:    "{module}",
			expected: []TokenType{TOKEN_START_MODULE, TOKEN_END_MODULE, TOKEN_EOS},
		},
		{
			name:     "definition tokens",
			input:    ": DOUBLE 2 * ;",
			expected: []TokenType{TOKEN_START_DEF, TOKEN_WORD, TOKEN_WORD, TOKEN_END_DEF, TOKEN_EOS},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input, nil, false)
			var tokens []TokenType

			for {
				token, err := tokenizer.NextToken()
				assert.NoError(t, err)
				tokens = append(tokens, token.Type)
				if token.Type == TOKEN_EOS {
					break
				}
			}

			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestTokenizerStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "double quote string",
			input:    `"hello world"`,
			expected: "hello world",
		},
		{
			name:     "single quote string",
			input:    `'hello world'`,
			expected: "hello world",
		},
		{
			name:     "caret quote string",
			input:    `^hello world^`,
			expected: "hello world",
		},
		{
			name:     "triple quote string",
			input:    "\"\"\"multi\nline\nstring\"\"\"",
			expected: "multi\nline\nstring",
		},
		{
			name:     "empty string",
			input:    `""`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input, nil, false)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, TOKEN_STRING, token.Type)
			assert.Equal(t, tt.expected, token.String)
		})
	}
}

func TestTokenizerComments(t *testing.T) {
	input := "WORD1 # this is a comment\nWORD2"
	tokenizer := NewTokenizer(input, nil, false)

	// First token: WORD1
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_WORD, token.Type)
	assert.Equal(t, "WORD1", token.String)

	// Second token: comment
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_COMMENT, token.Type)
	assert.Contains(t, token.String, "this is a comment")

	// Third token: WORD2
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_WORD, token.Type)
	assert.Equal(t, "WORD2", token.String)
}

func TestTokenizerDotSymbol(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType TokenType
		expectedStr  string
	}{
		{
			name:         "simple dot symbol",
			input:        ".field",
			expectedType: TOKEN_DOT_SYMBOL,
			expectedStr:  "field",
		},
		{
			name:         "dot symbol with hyphen",
			input:        ".field-name",
			expectedType: TOKEN_DOT_SYMBOL,
			expectedStr:  "field-name",
		},
		{
			name:         "lone dot is word",
			input:        ".",
			expectedType: TOKEN_WORD,
			expectedStr:  ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizer := NewTokenizer(tt.input, nil, false)
			token, err := tokenizer.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedType, token.Type)
			assert.Equal(t, tt.expectedStr, token.String)
		})
	}
}

func TestTokenizerMemo(t *testing.T) {
	input := "@: MEMOIZED 2 * ;"
	tokenizer := NewTokenizer(input, nil, false)

	// First token: START_MEMO with name
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_START_MEMO, token.Type)
	assert.Equal(t, "MEMOIZED", token.String)

	// Second token: WORD (2)
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_WORD, token.Type)
	assert.Equal(t, "2", token.String)
}

func TestTokenizerRFC9557DateTime(t *testing.T) {
	input := "2025-05-20T08:00:00[America/Los_Angeles]"
	tokenizer := NewTokenizer(input, nil, false)

	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, TOKEN_WORD, token.Type)
	assert.Equal(t, "2025-05-20T08:00:00[America/Los_Angeles]", token.String)
}

func TestTokenizerWhitespace(t *testing.T) {
	input := "WORD1\t\tWORD2\n\nWORD3"
	tokenizer := NewTokenizer(input, nil, false)

	expected := []string{"WORD1", "WORD2", "WORD3"}
	for _, exp := range expected {
		token, err := tokenizer.NextToken()
		assert.NoError(t, err)
		assert.Equal(t, TOKEN_WORD, token.Type)
		assert.Equal(t, exp, token.String)
	}
}

func TestTokenizerLocationTracking(t *testing.T) {
	input := "WORD1\nWORD2"
	tokenizer := NewTokenizer(input, nil, false)

	// First token on line 1
	token, err := tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, 1, token.Location.Line)
	assert.Equal(t, 1, token.Location.Column)

	// Second token on line 2
	token, err = tokenizer.NextToken()
	assert.NoError(t, err)
	assert.Equal(t, 2, token.Location.Line)
	assert.Equal(t, 1, token.Location.Column)
}
