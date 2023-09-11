package tokenizer

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestDefaultTokenization(t *testing.T) {
    tokenizer := New()

    actual := tokenizer.Tokenize("This is a test sample (it contains many characters) some of which aren't\n very straightforward \"to get\" right. Do you think so? I do.")

    expected := []string{"This", "is", "a", "test", "sample", "it", "contains", "many", "characters", "some", "of", "which", "aren't", "very", "straightforward", "to", "get", "right", "Do", "you", "think", "so", "I", "do"}
    assert.Equal(t, expected, actual)
}

func TestCustomTokenization(t *testing.T) {
    tokenizer := NewWithSeparators(",\r\n")

    actual := tokenizer.Tokenize("Testing. That,\r\nindeed,\ris what we\ndo here.")

    expected := []string{"Testing. That", "indeed", "is what we", "do here."}
    assert.Equal(t, expected, actual)
}

func TestEndingWithSeparator(t *testing.T) {
    tokenizer := NewWithSeparators(",\r\n")

    actual := tokenizer.Tokenize("This is\na test\r")

    expected := []string{"This is", "a test"}
    assert.Equal(t, expected, actual)
}

func TestEndingWithNonSeparator(t *testing.T) {
    tokenizer := NewWithSeparators(",\r\n")

    actual := tokenizer.Tokenize("This is\nalso, a test.")

    expected := []string{"This is", "also", " a test."}
    assert.Equal(t, expected, actual)
}

func TestOnlySeparators(t *testing.T) {
    tokenizer := NewWithSeparators(",\r\n")

    actual := tokenizer.Tokenize("\r\n,")

    expected := []string{}
    assert.Equal(t, expected, actual)
}

func TestNoSeparators(t *testing.T) {
    tokenizer := NewWithSeparators(",\r\n")

    actual := tokenizer.Tokenize("This is a fabulous test.")

    expected := []string{"This is a fabulous test."}
    assert.Equal(t, expected, actual)
}
