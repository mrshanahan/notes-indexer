package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTokenization(t *testing.T) {
	tokenizer := NewDefault()

	actual, _ := tokenizer.Tokenize("This is a test sample (it contains many characters) some of which aren't\n very straightforward \"to get\" right. Do you think so? I do.")

	expected := fmap([]string{"this", "is", "a", "test", "sample", "it", "contains", "many", "characters", "some", "of", "which", "aren't", "very", "straightforward", "to", "get", "right", "do", "you", "think", "so", "i", "do"}, genericToken)
	assert.Equal(t, expected, actual)
}

func TestXmlElementDetection(t *testing.T) {
	tests := []struct {
		input        string
		isXmlElement bool
	}{
		{"<b>", true},
		{"</b>", true},
		{"<b/>", true},
		{"<b />", true},
		{"<!-- b/ -->", true},
		{"&amp;", true},
		{"&#09;", true},
		{"&#xFF09;", true},
		{"<b />abc", true},
		{"<abc>", true},
		{"<abc b=\"f\">", true},
		{"<abc b=\"f&amp;d\">", true},
		{"<abc b=\"f&amp;d\" efe=\"ac123\">", true},
		{"abc", false},
		{"<b/ >", false},
		{"< b/>", false},
		{"<!-- b/>", false},
		{"abc<b />", false},
		{"& elsewhere;", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(s *testing.T) {
			actual := patt_XmlEnt.MatchString(test.input)
			assert.Equal(s, test.isXmlElement, actual)
		})
	}
}

func TestHtmlTokenization(t *testing.T) {
	tokenizer := NewXmlTokenizer()

	actual, _ := tokenizer.Tokenize("This is a <b>test</b> sample<br/>some of which aren't\n very&amp;straightforward \"to get\" right. Do you think so?<div>I do.")

	expected := []Token{
		genericToken("this"),
		genericToken("is"),
		genericToken("a"),
		xmlToken("<b>"),
		genericToken("test"),
		xmlToken("</b>"),
		genericToken("sample"),
		xmlToken("<br/>"),
		genericToken("some"),
		genericToken("of"),
		genericToken("which"),
		genericToken("aren't"),
		genericToken("very"),
		xmlToken("&amp;"),
		genericToken("straightforward"),
		genericToken("to"),
		genericToken("get"),
		genericToken("right"),
		genericToken("do"),
		genericToken("you"),
		genericToken("think"),
		genericToken("so"),
		xmlToken("<div>"),
		genericToken("i"),
		genericToken("do"),
	}
	assert.Equal(t, expected, actual)
}

func TestCustomTokenization(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("Testing. That,\r\nindeed,\ris what we\ndo here.")

	expected := fmap([]string{"testing. that", "indeed", "is what we", "do here."}, genericToken)
	assert.Equal(t, expected, actual)
}

func TestEndingWithSeparator(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("This is\na test\r")

	expected := fmap([]string{"this is", "a test"}, genericToken)
	assert.Equal(t, expected, actual)
}

func TestEndingWithNonSeparator(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("This is\nalso, a test.")

	expected := fmap([]string{"this is", "also", " a test."}, genericToken)
	assert.Equal(t, expected, actual)
}

func TestOnlySeparators(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("\r\n,")

	expected := []Token{}
	assert.Equal(t, expected, actual)
}

func TestNoSeparators(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("This is a fabulous test.")

	expected := fmap([]string{"this is a fabulous test."}, genericToken)
	assert.Equal(t, expected, actual)
}

func fmap[T any, S any](ts []T, f func(T) S) []S {
	ss := []S{}
	for _, t := range ts {
		ss = append(ss, f(t))
	}
	return ss
}

func genericToken(t string) Token {
	return Token{t, TOKEN_TYPE_GENERIC}
}

func xmlToken(t string) Token {
	return Token{t, TOKEN_TYPE_XML}
}
