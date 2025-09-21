package markdown

import (
	"reflect"
	"testing"
)

func TestParseGrammarRule(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput mdGrammarRule
		expectedError  error
	}{
		{
			"unordered_list_indic",
			&mdTokenRule{TOKEN_UNORDERED_LIST_INDIC},
			nil,
		},
	}

	for _, test := range tests {
		t.Run("parseGrammarRule", func(s *testing.T) {
			actualOutput, actualError := parseGrammarRule(test.input)
			if !reflect.DeepEqual(test.expectedOutput, actualOutput) {
				s.Errorf("outputs not equal - expected=%v, actual=%v", test.expectedOutput, actualOutput)
			}
			if !reflect.DeepEqual(test.expectedError, actualError) {
				s.Errorf("errors not equal - expected=%v, actual=%v", test.expectedError, actualError)
			}
		})
	}
}
