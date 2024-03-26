package stemmer

import (
	"bufio"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToken(t *testing.T) {
	t.Run("tree (m=0)", func(t *testing.T) {
		token := "tree"
		expected := &porterStructure{
			Token: token,
			M:     0,
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
			Components: []tokenComponent{
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
	path := "./stemmer_examples.txt"
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open examples file: %s (%s)", path, err)
	}
	defer f.Close()

	splitPatt := regexp.MustCompile(`^(\w+)\s*->\s*(\w+)$`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		match := splitPatt.FindAllStringSubmatch(line, -1)
		if match == nil {
			t.Fatalf("Line did not match expected format: %s", line)
		}
		input := match[0][1]
		expected := match[0][2]
		actual := Stem(input)
		t.Logf("Testing: %s", line)
		assert.Equal(t, expected, actual)
	}
}

func Test(t *testing.T) {
	splitPatt := regexp.MustCompile(`^(\w+)\s*->\s*(\w+)$`)
	match := splitPatt.FindAllStringSubmatch("abc -> bac", -1)
	t.Logf("%s", match)
}
