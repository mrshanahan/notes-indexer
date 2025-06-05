package tokenizer

import (
	"strings"
)

type defaultTokenizer struct {
	separators set[rune]
}

func NewDefault() Tokenizer {
	return &defaultTokenizer{
		separators: convertSeparator(DefaultSeparators),
	}
}

func NewDefaultWithSeparators(seps string) Tokenizer {
	return &defaultTokenizer{
		separators: convertSeparator(seps),
	}
}

func (t *defaultTokenizer) Tokenize(text string) ([]Token, error) {
	first, tokens, seps := -1, []Token{}, t.separators
	for i, r := range text {
		issep := seps[r]
		if !issep && first < 0 {
			first = i
		} else if issep && first >= 0 {
			tok := strings.ToLower(text[first:i])
			tokens = append(tokens, Token{tok, TOKEN_TYPE_GENERIC})
			first = -1
		}
	}
	if first >= 0 {
		tok := strings.ToLower(text[first:])
		tokens = append(tokens, Token{tok, TOKEN_TYPE_GENERIC})
	}
	return tokens, nil
}

type set[T comparable] map[T]bool

func convertSeparator(seps string) set[rune] {
	sepset := make(set[rune])
	for _, v := range seps {
		sepset[v] = true
	}
	return sepset
}
