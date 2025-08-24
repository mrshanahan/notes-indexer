package markdown

import (
	"reflect"
	"testing"
)

func TestLex(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		expectedLex []MDToken
	}{
		{
			"basic",
			"This is a test",
			[]MDToken{MDTextToken{"This is a test"}},
		},
		{
			"leading space",
			"    This is a test",
			[]MDToken{
				MDLeadingSpaceToken{4},
				MDTextToken{"This is a test"},
			},
		},
		{
			"leading-tabs",
			"\t  This is a test",
			[]MDToken{
				MDLeadingSpaceToken{6},
				MDTextToken{"This is a test"},
			},
		},
		{
			"newlines",
			"This is\na test\r\nof many things\n\nand various trials\n\r\n\n",
			[]MDToken{
				MDTextToken{"This is"},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"a test"},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"of many things"},
				MDSimpleToken{TOKEN_NL},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"and various trials"},
				MDSimpleToken{TOKEN_NL},
				MDSimpleToken{TOKEN_NL},
				MDSimpleToken{TOKEN_NL},
			},
		},
		{
			"list-basic",
			"- Initial list item\n    - Nested list item\n        2. Ordered item under that",
			[]MDToken{
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"Initial list item"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{4},
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"Nested list item"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{8},
				MDOrderedListIndicToken{"2. "},
				MDTextToken{"Ordered item under that"},
			},
		},
		{
			"list-preserve-spaces",
			"-   We keep the spaces\n 1.   In our lists",
			[]MDToken{
				MDUnorderedListIndicToken{"- "}, MDLeadingSpaceToken{2}, MDTextToken{"We keep the spaces"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{1}, MDOrderedListIndicToken{"1. "}, MDLeadingSpaceToken{2}, MDTextToken{"In our lists"},
			},
		},
		{
			"no-inline-list",
			"- This is a list with a - hyphen in it -",
			[]MDToken{
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"This is a list with a - hyphen in it -"},
			},
		},
		{
			"inline",
			"This is a **cool** paragraph, which has *many* things, such as __underlines__ and ~~strikethroughs~~. For example, **this bold is *also italicized***",
			[]MDToken{
				MDTextToken{"This is a "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "**"},
				MDTextToken{"cool"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "**"},
				MDTextToken{" paragraph, which has "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"},
				MDTextToken{"many"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "*"},
				MDTextToken{" things, such as "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "__"},
				MDTextToken{"underlines"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "__"},
				MDTextToken{" and "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "~~"},
				MDTextToken{"strikethroughs"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_MID, "~~"},
				MDTextToken{". For example, "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "**"},
				MDTextToken{"this bold is "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"},
				MDTextToken{"also italicized"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "***"},
			},
		},
		{
			"inline-format-with-list",
			"- This list has *several inline elements*.\n    - It has ~~several~~ a few things.",
			[]MDToken{
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"This list has "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"},
				MDTextToken{"several inline elements"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_MID, "*"},
				MDTextToken{"."},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{4},
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"It has "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "~~"},
				MDTextToken{"several"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "~~"},
				MDTextToken{" a few things."},
			},
		},
		{
			"inline-format-heterogeneous",
			"We have *several nested types `of elements`* at the ~~**same time**~~",
			[]MDToken{
				MDTextToken{"We have "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"},
				MDTextToken{"several nested types "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "`"},
				MDTextToken{"of elements"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "`*"},
				MDTextToken{" at the "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "~~**"},
				MDTextToken{"same time"},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "**~~"},
			},
		},
		{
			"inline-without-correct-spaces-ignored",
			"These * formatting elements will be ignored. ~",
			[]MDToken{
				MDTextToken{"These * formatting elements will be ignored. ~"},
			},
		},
		{
			"inline-recognizes-multiple-end",
			"This **has nested *formatting.***",
			[]MDToken{
				MDTextToken{"This "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "**"},
				MDTextToken{"has nested "},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"},
				MDTextToken{"formatting."},
				MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "***"},
			},
		},
		{
			"header",
			"# This is a basic header\nThis one is not\n##   and this one is a header again\n#but this one is not\n# - And this is not a list!\n## But we do keep *processing* inline elements",
			[]MDToken{
				MDHeaderIndicToken{1, "# "}, MDTextToken{"This is a basic header"},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"This one is not"},
				MDSimpleToken{TOKEN_NL},
				MDHeaderIndicToken{2, "## "}, MDLeadingSpaceToken{2}, MDTextToken{"and this one is a header again"},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"#but this one is not"},
				MDSimpleToken{TOKEN_NL},
				MDHeaderIndicToken{1, "# "}, MDUnorderedListIndicToken{"- "}, MDTextToken{"And this is not a list!"},
				MDSimpleToken{TOKEN_NL},
				MDHeaderIndicToken{2, "## "}, MDTextToken{"But we do keep "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"}, MDTextToken{"processing"}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "*"}, MDTextToken{" inline elements"},
			},
		},
		{
			"mixed-formatting-ordering",
			"Standard *formatting* line\n    - Plain list item with _formatting_\n # Header with leading space\n  ## Header with **formatting**\n- # List item with header and **formatting**\n 1. Same with **ordered** list\n# - Header ignores rest of list but *not formatting*",
			[]MDToken{
				MDTextToken{"Standard "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"}, MDTextToken{"formatting"}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "*"}, MDTextToken{" line"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{4}, MDUnorderedListIndicToken{"- "}, MDTextToken{"Plain list item with "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "_"}, MDTextToken{"formatting"}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "_"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{1}, MDHeaderIndicToken{1, "# "}, MDTextToken{"Header with leading space"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{2}, MDHeaderIndicToken{2, "## "}, MDTextToken{"Header with "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "**"}, MDTextToken{"formatting"}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "**"},
				MDSimpleToken{TOKEN_NL},
				MDUnorderedListIndicToken{"- "}, MDHeaderIndicToken{1, "# "}, MDTextToken{"List item with header and "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "**"}, MDTextToken{"formatting"}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "**"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{1}, MDOrderedListIndicToken{"1. "}, MDTextToken{"Same with "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "**"}, MDTextToken{"ordered"}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "**"}, MDTextToken{" list"},
				MDSimpleToken{TOKEN_NL},
				MDHeaderIndicToken{1, "# "}, MDUnorderedListIndicToken{"- "}, MDTextToken{"Header ignores rest of list but "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "*"}, MDTextToken{"not formatting"}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_END, "*"},
			},
		},
		{
			"escape-inline",
			"This is an escape\\* line\nAs is ~~The last one\\~\\~ here\\\\.\nBut we don't get any escaped chars for whitespace\\ or end of line.\\\nOr end of text.\\",
			[]MDToken{
				MDTextToken{"This is an escape"}, MDEscapeToken{"*"}, MDTextToken{" line"},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"As is "}, MDInlineFormatToken{TOKEN_INLINE_FORMAT_START, "~~"}, MDTextToken{"The last one"}, MDEscapeToken{"~"}, MDEscapeToken{"~"}, MDTextToken{" here"}, MDEscapeToken{"\\"}, MDTextToken{"."},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"But we don't get any escaped chars for whitespace"}, MDEscapeToken{""}, MDTextToken{" or end of line."}, MDEscapeToken{""},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"Or end of text."}, MDEscapeToken{""},
			},
		},
		{
			"escape-structural-elements",
			"\\- This won't be a list\n    \\- Neither will this\n1\\. Nor this\n\\1. And this neither\nAnd this\\nwon't be a newline\n\\# And this won't be a header\n\\### And neither will this",
			[]MDToken{
				MDEscapeToken{"-"}, MDTextToken{" This won't be a list"},
				MDSimpleToken{TOKEN_NL},
				MDLeadingSpaceToken{4}, MDEscapeToken{"-"}, MDTextToken{" Neither will this"},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"1"}, MDEscapeToken{"."}, MDTextToken{" Nor this"},
				MDSimpleToken{TOKEN_NL},
				MDEscapeToken{"1"}, MDTextToken{". And this neither"},
				MDSimpleToken{TOKEN_NL},
				MDTextToken{"And this"}, MDEscapeToken{"n"}, MDTextToken{"won't be a newline"},
				MDSimpleToken{TOKEN_NL},
				MDEscapeToken{"#"}, MDTextToken{" And this won't be a header"},
				MDSimpleToken{TOKEN_NL},
				MDEscapeToken{"#"}, MDTextToken{"## And neither will this"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(s *testing.T) {
			actualLex := Lex(test.text)
			if !reflect.DeepEqual(test.expectedLex, actualLex) {
				s.Errorf("lexes were not equal - expected=%v, actual=%v", test.expectedLex, actualLex)
			}
		})
	}
}
