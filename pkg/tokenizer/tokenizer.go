package tokenizer

import (
	"strings"
)

var (
	DefaultSeparators = "\t\n\r ,.:?\"!;()"
)

type Tokenizer interface {
	Tokenize(text string) []string
}

type tokenizer struct {
	separators set[rune]
}

func New() Tokenizer {
	return &tokenizer{
		separators: convertSeparator(DefaultSeparators),
	}
}

func NewWithSeparators(seps string) Tokenizer {
	return &tokenizer{
		separators: convertSeparator(seps),
	}
}

func (t *tokenizer) Tokenize(text string) []string {
	first, tokens, seps := -1, []string{}, t.separators
	for i, r := range text {
		issep := seps[r]
		if !issep && first < 0 {
			first = i
		} else if issep && first >= 0 {
			tok := strings.ToLower(text[first:i])
			tokens = append(tokens, tok)
			first = -1
		}
	}
	if first >= 0 {
		tok := strings.ToLower(text[first:])
		tokens = append(tokens, tok)
	}
	return tokens
}

type set[T comparable] map[T]bool

func newSet[T comparable]() set[T] {
	return set[T]{}
}

func convertSeparator(seps string) set[rune] {
	sepset := make(set[rune])
	for _, v := range seps {
		sepset[v] = true
	}
	return sepset
}
