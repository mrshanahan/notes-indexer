package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultTokenization(t *testing.T) {
	tokenizer := NewDefault()

	actual, _ := tokenizer.Tokenize("This is a test sample (it contains many characters) some of which aren't\n very straightforward \"to get\" right. Do you think so? I do.")

	expected := []string{"this", "is", "a", "test", "sample", "it", "contains", "many", "characters", "some", "of", "which", "aren't", "very", "straightforward", "to", "get", "right", "do", "you", "think", "so", "i", "do"}
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

	expected := []string{"this", "is", "a", "<b>", "test", "</b>", "sample", "<br/>", "some", "of", "which", "aren't", "very", "&amp;", "straightforward", "to", "get", "right", "do", "you", "think", "so", "<div>", "i", "do"}
	assert.Equal(t, expected, actual)
}

func TestCustomTokenization(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("Testing. That,\r\nindeed,\ris what we\ndo here.")

	expected := []string{"Testing. That", "indeed", "is what we", "do here."}
	assert.Equal(t, expected, actual)
}

func TestEndingWithSeparator(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("This is\na test\r")

	expected := []string{"This is", "a test"}
	assert.Equal(t, expected, actual)
}

func TestEndingWithNonSeparator(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("This is\nalso, a test.")

	expected := []string{"This is", "also", " a test."}
	assert.Equal(t, expected, actual)
}

func TestOnlySeparators(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("\r\n,")

	expected := []string{}
	assert.Equal(t, expected, actual)
}

func TestNoSeparators(t *testing.T) {
	tokenizer := NewDefaultWithSeparators(",\r\n")

	actual, _ := tokenizer.Tokenize("This is a fabulous test.")

	expected := []string{"This is a fabulous test."}
	assert.Equal(t, expected, actual)
}
