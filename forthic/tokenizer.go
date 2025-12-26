package forthic

import (
	"strings"
)

// ============================================================================
// Token Types
// ============================================================================

type TokenType int

const (
	TOKEN_STRING TokenType = iota + 1
	TOKEN_COMMENT
	TOKEN_START_ARRAY
	TOKEN_END_ARRAY
	TOKEN_START_MODULE
	TOKEN_END_MODULE
	TOKEN_START_DEF
	TOKEN_END_DEF
	TOKEN_START_MEMO
	TOKEN_WORD
	TOKEN_DOT_SYMBOL
	TOKEN_EOS
)

// ============================================================================
// Token
// ============================================================================

type Token struct {
	Type     TokenType
	String   string
	Location *CodeLocation
}

func NewToken(tokenType TokenType, str string, location *CodeLocation) *Token {
	return &Token{
		Type:     tokenType,
		String:   str,
		Location: location,
	}
}

// ============================================================================
// String Delta (for string tokens)
// ============================================================================

type stringDelta struct {
	start int
	end   int
}

// ============================================================================
// Tokenizer
// ============================================================================

type Tokenizer struct {
	referenceLocation *CodeLocation
	line              int
	column            int
	inputString       string
	inputPos          int
	whitespace        []rune
	quoteChars        []rune
	tokenStartPos     int
	tokenEndPos       int
	tokenLine         int
	tokenColumn       int
	tokenString       strings.Builder
	stringDelta       *stringDelta
	streaming         bool
}

func NewTokenizer(inputString string, referenceLocation *CodeLocation, streaming bool) *Tokenizer {
	if referenceLocation == nil {
		referenceLocation = &CodeLocation{Source: "", Line: 1, Column: 1, StartPos: 0, EndPos: 0}
	}

	return &Tokenizer{
		referenceLocation: referenceLocation,
		line:              referenceLocation.Line,
		column:            referenceLocation.Column,
		inputString:       unescapeString(inputString),
		inputPos:          0,
		whitespace:        []rune{' ', '\t', '\n', '\r', '(', ')', ','},
		quoteChars:        []rune{'"', '\'', '^'},
		tokenStartPos:     0,
		tokenEndPos:       0,
		tokenLine:         0,
		tokenColumn:       0,
		stringDelta:       nil,
		streaming:         streaming,
	}
}

// ============================================================================
// Helper Functions
// ============================================================================

func unescapeString(s string) string {
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	return s
}

func (t *Tokenizer) clearTokenString() {
	t.tokenString.Reset()
}

func (t *Tokenizer) noteStartToken() {
	t.tokenStartPos = t.inputPos + t.referenceLocation.StartPos
	t.tokenLine = t.line
	t.tokenColumn = t.column
}

func (t *Tokenizer) isWhitespace(ch rune) bool {
	for _, ws := range t.whitespace {
		if ch == ws {
			return true
		}
	}
	return false
}

func (t *Tokenizer) isQuote(ch rune) bool {
	for _, qc := range t.quoteChars {
		if ch == qc {
			return true
		}
	}
	return false
}

func (t *Tokenizer) isTripleQuote(index int, ch rune) bool {
	if !t.isQuote(ch) {
		return false
	}
	if index+2 >= len(t.inputString) {
		return false
	}
	return rune(t.inputString[index+1]) == ch && rune(t.inputString[index+2]) == ch
}

func (t *Tokenizer) isStartMemo(index int) bool {
	if index+1 >= len(t.inputString) {
		return false
	}
	return t.inputString[index] == '@' && t.inputString[index+1] == ':'
}

func (t *Tokenizer) advancePosition(numChars int) {
	if numChars >= 0 {
		for i := 0; i < numChars; i++ {
			if t.inputPos < len(t.inputString) && t.inputString[t.inputPos] == '\n' {
				t.line++
				t.column = 1
			} else {
				t.column++
			}
			t.inputPos++
		}
	} else {
		for i := 0; i < -numChars; i++ {
			t.inputPos--
			if t.inputPos < 0 {
				panic("Invalid input position")
			}
			if t.inputString[t.inputPos] == '\n' {
				t.line--
				t.column = 1
			} else {
				t.column--
			}
		}
	}
}

func (t *Tokenizer) getTokenLocation() *CodeLocation {
	return &CodeLocation{
		Source:   t.referenceLocation.Source,
		Line:     t.tokenLine,
		Column:   t.tokenColumn,
		StartPos: t.tokenStartPos,
		EndPos:   t.tokenStartPos + t.tokenString.Len(),
	}
}

// ============================================================================
// Public API
// ============================================================================

func (t *Tokenizer) NextToken() (*Token, error) {
	t.clearTokenString()
	return t.transitionFromSTART()
}

// ============================================================================
// State Transitions
// ============================================================================

func (t *Tokenizer) transitionFromSTART() (*Token, error) {
	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.noteStartToken()
		t.advancePosition(1)

		if t.isWhitespace(ch) {
			continue
		} else if ch == '#' {
			return t.transitionFromCOMMENT()
		} else if ch == ':' {
			return t.transitionFromSTART_DEFINITION()
		} else if t.isStartMemo(t.inputPos - 1) {
			t.advancePosition(1) // Skip over ":" in "@:"
			return t.transitionFromSTART_MEMO()
		} else if ch == ';' {
			t.tokenString.WriteRune(ch)
			return NewToken(TOKEN_END_DEF, string(ch), t.getTokenLocation()), nil
		} else if ch == '[' {
			t.tokenString.WriteRune(ch)
			return NewToken(TOKEN_START_ARRAY, string(ch), t.getTokenLocation()), nil
		} else if ch == ']' {
			t.tokenString.WriteRune(ch)
			return NewToken(TOKEN_END_ARRAY, string(ch), t.getTokenLocation()), nil
		} else if ch == '{' {
			return t.transitionFromGATHER_MODULE()
		} else if ch == '}' {
			t.tokenString.WriteRune(ch)
			return NewToken(TOKEN_END_MODULE, string(ch), t.getTokenLocation()), nil
		} else if t.isTripleQuote(t.inputPos-1, ch) {
			t.advancePosition(2) // Skip over 2nd and 3rd quote chars
			return t.transitionFromGATHER_TRIPLE_QUOTE_STRING(ch)
		} else if t.isQuote(ch) {
			return t.transitionFromGATHER_STRING(ch)
		} else if ch == '.' {
			t.advancePosition(-1) // Back up to beginning of dot symbol
			return t.transitionFromGATHER_DOT_SYMBOL()
		} else {
			t.advancePosition(-1) // Back up to beginning of word
			return t.transitionFromGATHER_WORD()
		}
	}
	return NewToken(TOKEN_EOS, "", t.getTokenLocation()), nil
}

func (t *Tokenizer) transitionFromCOMMENT() (*Token, error) {
	t.noteStartToken()
	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.tokenString.WriteRune(ch)
		t.advancePosition(1)
		if ch == '\n' {
			t.advancePosition(-1)
			break
		}
	}
	return NewToken(TOKEN_COMMENT, t.tokenString.String(), t.getTokenLocation()), nil
}

func (t *Tokenizer) transitionFromSTART_DEFINITION() (*Token, error) {
	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.advancePosition(1)

		if t.isWhitespace(ch) {
			continue
		} else if t.isQuote(ch) {
			return nil, NewForthicError("Definition names can't have quotes in them").
				WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
		} else {
			t.advancePosition(-1)
			return t.transitionFromGATHER_DEFINITION_NAME()
		}
	}

	return nil, NewForthicError("Got EOS in START_DEFINITION").
		WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
}

func (t *Tokenizer) transitionFromSTART_MEMO() (*Token, error) {
	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.advancePosition(1)

		if t.isWhitespace(ch) {
			continue
		} else if t.isQuote(ch) {
			return nil, NewForthicError("Memo names can't have quotes in them").
				WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
		} else {
			t.advancePosition(-1)
			return t.transitionFromGATHER_MEMO_NAME()
		}
	}

	return nil, NewForthicError("Got EOS in START_MEMO").
		WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
}

func (t *Tokenizer) gatherDefinitionName() error {
	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.advancePosition(1)

		if t.isWhitespace(ch) {
			break
		}
		if t.isQuote(ch) {
			return NewForthicError("Definition names can't have quotes in them").
				WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
		}
		if strings.ContainsRune("[]{}", ch) {
			return NewForthicError("Definition names can't have '" + string(ch) + "' in them").
				WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
		}
		t.tokenString.WriteRune(ch)
	}
	return nil
}

func (t *Tokenizer) transitionFromGATHER_DEFINITION_NAME() (*Token, error) {
	t.noteStartToken()
	if err := t.gatherDefinitionName(); err != nil {
		return nil, err
	}
	return NewToken(TOKEN_START_DEF, t.tokenString.String(), t.getTokenLocation()), nil
}

func (t *Tokenizer) transitionFromGATHER_MEMO_NAME() (*Token, error) {
	t.noteStartToken()
	if err := t.gatherDefinitionName(); err != nil {
		return nil, err
	}
	return NewToken(TOKEN_START_MEMO, t.tokenString.String(), t.getTokenLocation()), nil
}

func (t *Tokenizer) transitionFromGATHER_MODULE() (*Token, error) {
	t.noteStartToken()
	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.advancePosition(1)

		if t.isWhitespace(ch) {
			break
		} else if ch == '}' {
			t.advancePosition(-1)
			break
		} else {
			t.tokenString.WriteRune(ch)
		}
	}
	return NewToken(TOKEN_START_MODULE, t.tokenString.String(), t.getTokenLocation()), nil
}

func (t *Tokenizer) transitionFromGATHER_TRIPLE_QUOTE_STRING(delim rune) (*Token, error) {
	t.noteStartToken()
	stringDelimiter := delim
	t.stringDelta = &stringDelta{start: t.inputPos, end: t.inputPos}

	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])

		if ch == stringDelimiter && t.isTripleQuote(t.inputPos, ch) {
			// Check if this triple quote is followed by at least one more quote (greedy mode trigger)
			if t.inputPos+3 < len(t.inputString) && rune(t.inputString[t.inputPos+3]) == stringDelimiter {
				// Greedy mode: include this quote as content and continue looking for the end
				t.advancePosition(1) // Advance by 1 to catch overlapping sequences
				t.tokenString.WriteRune(stringDelimiter)
				t.stringDelta.end = t.inputPos
				continue
			}

			// Normal behavior: close at first triple quote
			t.advancePosition(3)
			token := NewToken(TOKEN_STRING, t.tokenString.String(), t.getTokenLocation())
			t.stringDelta = nil
			return token, nil
		} else {
			t.advancePosition(1)
			t.tokenString.WriteRune(ch)
			t.stringDelta.end = t.inputPos
		}
	}

	if t.streaming {
		return nil, nil
	}
	return nil, NewForthicError("Unterminated string").
		WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
}

func (t *Tokenizer) transitionFromGATHER_STRING(delim rune) (*Token, error) {
	t.noteStartToken()
	stringDelimiter := delim
	t.stringDelta = &stringDelta{start: t.inputPos, end: t.inputPos}

	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.advancePosition(1)

		if ch == stringDelimiter {
			token := NewToken(TOKEN_STRING, t.tokenString.String(), t.getTokenLocation())
			t.stringDelta = nil
			return token, nil
		} else {
			t.tokenString.WriteRune(ch)
			t.stringDelta.end = t.inputPos
		}
	}

	if t.streaming {
		return nil, nil
	}
	return nil, NewForthicError("Unterminated string").
		WithLocation(&CodeLocation{Line: t.tokenLine, Column: t.tokenColumn})
}

func (t *Tokenizer) transitionFromGATHER_WORD() (*Token, error) {
	t.noteStartToken()
	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.advancePosition(1)

		if t.isWhitespace(ch) {
			break
		}
		if strings.ContainsRune(";{}#", ch) {
			t.advancePosition(-1)
			break
		}

		// Handle RFC 9557 datetime with IANA timezone: 2025-05-20T08:00:00[America/Los_Angeles]
		// When we see '[', check if token looks like a datetime (contains 'T')
		// If so, include the bracketed timezone as part of the token
		if ch == '[' {
			if strings.Contains(t.tokenString.String(), "T") {
				// This looks like a datetime, gather until ']'
				t.tokenString.WriteRune(ch)
				for t.inputPos < len(t.inputString) {
					tzChar := rune(t.inputString[t.inputPos])
					t.advancePosition(1)
					t.tokenString.WriteRune(tzChar)
					if tzChar == ']' {
						break
					}
				}
				break
			} else {
				// Not a datetime, treat '[' as delimiter
				t.advancePosition(-1)
				break
			}
		}
		if ch == ']' {
			t.advancePosition(-1)
			break
		}
		t.tokenString.WriteRune(ch)
	}
	return NewToken(TOKEN_WORD, t.tokenString.String(), t.getTokenLocation()), nil
}

func (t *Tokenizer) transitionFromGATHER_DOT_SYMBOL() (*Token, error) {
	t.noteStartToken()
	fullTokenString := strings.Builder{}

	for t.inputPos < len(t.inputString) {
		ch := rune(t.inputString[t.inputPos])
		t.advancePosition(1)

		if t.isWhitespace(ch) {
			break
		}
		if strings.ContainsRune(";[]{}#", ch) {
			t.advancePosition(-1)
			break
		} else {
			fullTokenString.WriteRune(ch)
			t.tokenString.WriteRune(ch)
		}
	}

	fullToken := fullTokenString.String()

	// If dot symbol has no characters after the dot, treat it as a word
	if len(fullToken) < 2 { // "." + at least 1 char = 2 minimum
		return NewToken(TOKEN_WORD, fullToken, t.getTokenLocation()), nil
	}

	// For DOT_SYMBOL, return the string without the dot prefix
	symbolWithoutDot := fullToken[1:]
	return NewToken(TOKEN_DOT_SYMBOL, symbolWithoutDot, t.getTokenLocation()), nil
}
