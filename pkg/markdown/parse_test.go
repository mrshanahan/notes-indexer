package markdown

// func TestParse(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		tokens        []MDToken
// 		expectedParse MDSyntaxTree
// 		expectedError error
// 	}{
// 		{
// 			"plaintext",
// 			[]MDToken{MDTextToken{"This is a test"}},
// 			MDSyntaxTree{[]MDSyntaxNode{paragraph(text("This is a test"))}},
// 			nil,
// 		},
// 		{
// 			"multiple-lines-one-paragraph",
// 			[]MDToken{MDTextToken{"This is a test"}, MDSimpleToken{TOKEN_NL}, MDTextToken{"and so is this"}},
// 			MDSyntaxTree{[]MDSyntaxNode{paragraph(text("This is a test and so is this"))}},
// 			nil,
// 		},
// 		{
// 			"multiple-paragraphs",
// 			[]MDToken{MDTextToken{"This is a test"}, MDSimpleToken{TOKEN_NL}, MDSimpleToken{TOKEN_NL}, MDTextToken{"But this is a new paragraph"}},
// 			MDSyntaxTree{[]MDSyntaxNode{paragraph(text("This is a test"), text("But this is a new paragraph"))}},
// 			nil,
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(s *testing.T) {
// 			actualParse, actualError := Parse(test.tokens)
// 			if !reflect.DeepEqual(test.expectedParse, actualParse) {
// 				s.Errorf("parses were not equal - expected=%v, actual=%v", test.expectedParse, actualParse)
// 			}
// 			if !reflect.DeepEqual(test.expectedError, actualError) {
// 				s.Errorf("parses did not have same error - expected=%v, actual=%v", test.expectedError, actualError)
// 			}
// 		})
// 	}
// }

func text(c string) MDParagraphFormatNode {
	return MDTextFormatNode{c}
}

func paragraph(ns ...MDParagraphFormatNode) MDParagraph {
	return MDParagraph{ns}
}
