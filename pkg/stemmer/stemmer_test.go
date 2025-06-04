package stemmer

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	t.Run("tree (m=0)", func(t *testing.T) {
		token := "tree"
		expected := &porterStructure{
			Token: token,
			M:     0,
			Components: []*tokenComponent{
				{TYPE_CONSONANT, "tr"},
				{TYPE_VOWEL, "ee"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("trouble (m=1)", func(t *testing.T) {
		token := "trouble"
		expected := &porterStructure{
			Token: token,
			M:     1,
			Components: []*tokenComponent{
				{TYPE_CONSONANT, "tr"},
				{TYPE_VOWEL, "ou"},
				{TYPE_CONSONANT, "bl"},
				{TYPE_VOWEL, "e"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("oats (m=1)", func(t *testing.T) {
		token := "oats"
		expected := &porterStructure{
			Token: token,
			M:     1,
			Components: []*tokenComponent{
				{TYPE_VOWEL, "oa"},
				{TYPE_CONSONANT, "ts"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("trees (m=1)", func(t *testing.T) {
		token := "trees"
		expected := &porterStructure{
			Token: token,
			M:     1,
			Components: []*tokenComponent{
				{TYPE_CONSONANT, "tr"},
				{TYPE_VOWEL, "ee"},
				{TYPE_CONSONANT, "s"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("ivy (m=1)", func(t *testing.T) {
		token := "ivy"
		expected := &porterStructure{
			Token: token,
			M:     1,
			Components: []*tokenComponent{
				{TYPE_VOWEL, "i"},
				{TYPE_CONSONANT, "v"},
				{TYPE_VOWEL, "y"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("troubles (m=2)", func(t *testing.T) {
		token := "troubles"
		expected := &porterStructure{
			Token: token,
			M:     2,
			Components: []*tokenComponent{
				{TYPE_CONSONANT, "tr"},
				{TYPE_VOWEL, "ou"},
				{TYPE_CONSONANT, "bl"},
				{TYPE_VOWEL, "e"},
				{TYPE_CONSONANT, "s"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("private (m=2)", func(t *testing.T) {
		token := "private"
		expected := &porterStructure{
			Token: token,
			M:     2,
			Components: []*tokenComponent{
				{TYPE_CONSONANT, "pr"},
				{TYPE_VOWEL, "i"},
				{TYPE_CONSONANT, "v"},
				{TYPE_VOWEL, "a"},
				{TYPE_CONSONANT, "t"},
				{TYPE_VOWEL, "e"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("oaten (m=2)", func(t *testing.T) {
		token := "oaten"
		expected := &porterStructure{
			Token: token,
			M:     2,
			Components: []*tokenComponent{
				{TYPE_VOWEL, "oa"},
				{TYPE_CONSONANT, "t"},
				{TYPE_VOWEL, "e"},
				{TYPE_CONSONANT, "n"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
	t.Run("orrery (m=2)", func(t *testing.T) {
		token := "orrery"
		expected := &porterStructure{
			Token: token,
			M:     2,
			Components: []*tokenComponent{
				{TYPE_VOWEL, "o"},
				{TYPE_CONSONANT, "rr"},
				{TYPE_VOWEL, "e"},
				{TYPE_CONSONANT, "r"},
				{TYPE_VOWEL, "y"},
			},
		}
		actual := parseToken(token)
		assert.Equal(t, expected, actual)
	})
}

func TestStem(t *testing.T) {
	inpPath := "./sample_data/voc.txt"
	expPath := "./sample_data/output_porter.txt"

	inpf, err := os.Open(inpPath)
	if err != nil {
		t.Fatalf("Failed to open input vocabulary file: %s (%s)", inpPath, err)
	}
	defer inpf.Close()

	expf, err := os.Open(expPath)
	if err != nil {
		t.Fatalf("Failed to open expected output file: %s (%s)", expPath, err)
	}
	defer expf.Close()

	inpScanner, expScanner := bufio.NewScanner(inpf), bufio.NewScanner(expf)
	for inpScanner.Scan() && expScanner.Scan() {
		input, expected := inpScanner.Text(), expScanner.Text()
		actual := Stem(input)
		assert.Equal(t, expected, actual, "input: %s", input)
	}
}
