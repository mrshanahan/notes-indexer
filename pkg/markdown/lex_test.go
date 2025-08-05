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
				MDSimpleToken{NL},
				MDTextToken{"a test"},
				MDSimpleToken{NL},
				MDTextToken{"of many things"},
				MDSimpleToken{NL},
				MDSimpleToken{NL},
				MDTextToken{"and various trials"},
				MDSimpleToken{NL},
				MDSimpleToken{NL},
				MDSimpleToken{NL},
			},
		},
		{
			"list-basic",
			"- Initial list item\n    - Nested list item\n        2. Ordered item under that",
			[]MDToken{
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"Initial list item"},
				MDSimpleToken{NL},
				MDLeadingSpaceToken{4},
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"Nested list item"},
				MDSimpleToken{NL},
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
				MDSimpleToken{NL},
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
				MDInlineFormatToken{INLINE_FORMAT_START, "**"},
				MDTextToken{"cool"},
				MDInlineFormatToken{INLINE_FORMAT_END, "**"},
				MDTextToken{" paragraph, which has "},
				MDInlineFormatToken{INLINE_FORMAT_START, "*"},
				MDTextToken{"many"},
				MDInlineFormatToken{INLINE_FORMAT_END, "*"},
				MDTextToken{" things, such as "},
				MDInlineFormatToken{INLINE_FORMAT_START, "__"},
				MDTextToken{"underlines"},
				MDInlineFormatToken{INLINE_FORMAT_END, "__"},
				MDTextToken{" and "},
				MDInlineFormatToken{INLINE_FORMAT_START, "~~"},
				MDTextToken{"strikethroughs"},
				MDInlineFormatToken{INLINE_FORMAT_MID, "~~"},
				MDTextToken{". For example, "},
				MDInlineFormatToken{INLINE_FORMAT_START, "**"},
				MDTextToken{"this bold is "},
				MDInlineFormatToken{INLINE_FORMAT_START, "*"},
				MDTextToken{"also italicized"},
				MDInlineFormatToken{INLINE_FORMAT_END, "***"},
			},
		},
		{
			"inline-format-with-list",
			"- This list has *several inline elements*.\n    - It has ~~several~~ a few things.",
			[]MDToken{
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"This list has "},
				MDInlineFormatToken{INLINE_FORMAT_START, "*"},
				MDTextToken{"several inline elements"},
				MDInlineFormatToken{INLINE_FORMAT_MID, "*"},
				MDTextToken{"."},
				MDSimpleToken{NL},
				MDLeadingSpaceToken{4},
				MDUnorderedListIndicToken{"- "},
				MDTextToken{"It has "},
				MDInlineFormatToken{INLINE_FORMAT_START, "~~"},
				MDTextToken{"several"},
				MDInlineFormatToken{INLINE_FORMAT_END, "~~"},
				MDTextToken{" a few things."},
			},
		},
		{
			"inline-format-heterogeneous",
			"We have *several nested types `of elements`* at the ~~**same time**~~",
			[]MDToken{
				MDTextToken{"We have "},
				MDInlineFormatToken{INLINE_FORMAT_START, "*"},
				MDTextToken{"several nested types "},
				MDInlineFormatToken{INLINE_FORMAT_START, "`"},
				MDTextToken{"of elements"},
				MDInlineFormatToken{INLINE_FORMAT_END, "`*"},
				MDTextToken{" at the "},
				MDInlineFormatToken{INLINE_FORMAT_START, "~~**"},
				MDTextToken{"same time"},
				MDInlineFormatToken{INLINE_FORMAT_END, "**~~"},
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
				MDInlineFormatToken{INLINE_FORMAT_START, "**"},
				MDTextToken{"has nested "},
				MDInlineFormatToken{INLINE_FORMAT_START, "*"},
				MDTextToken{"formatting."},
				MDInlineFormatToken{INLINE_FORMAT_END, "***"},
			},
		},
		{
			"header",
			"# This is a basic header\nThis one is not\n##   and this one is a header again\n#but this one is not\n# - And this is not a list!\n## But we do keep *processing* inline elements",
			[]MDToken{
				MDHeaderIndicToken{1, "# "}, MDTextToken{"This is a basic header"},
				MDSimpleToken{NL},
				MDTextToken{"This one is not"},
				MDSimpleToken{NL},
				MDHeaderIndicToken{2, "## "}, MDLeadingSpaceToken{2}, MDTextToken{"and this one is a header again"},
				MDSimpleToken{NL},
				MDTextToken{"#but this one is not"},
				MDSimpleToken{NL},
				MDHeaderIndicToken{1, "# "}, MDUnorderedListIndicToken{"- "}, MDTextToken{"And this is not a list!"},
				MDSimpleToken{NL},
				MDHeaderIndicToken{2, "## "}, MDTextToken{"But we do keep "}, MDInlineFormatToken{INLINE_FORMAT_START, "*"}, MDTextToken{"processing"}, MDInlineFormatToken{INLINE_FORMAT_END, "*"}, MDTextToken{" inline elements"},
			},
		},
		{
			"mixed-formatting-ordering",
			"Standard *formatting* line\n    - Plain list item with _formatting_\n # Header with leading space\n  ## Header with **formatting**\n- # List item with header and **formatting**\n 1. Same with **ordered** list\n# - Header ignores rest of list but *not formatting*",
			[]MDToken{
				MDTextToken{"Standard "}, MDInlineFormatToken{INLINE_FORMAT_START, "*"}, MDTextToken{"formatting"}, MDInlineFormatToken{INLINE_FORMAT_END, "*"}, MDTextToken{" line"},
				MDSimpleToken{NL},
				MDLeadingSpaceToken{4}, MDUnorderedListIndicToken{"- "}, MDTextToken{"Plain list item with "}, MDInlineFormatToken{INLINE_FORMAT_START, "_"}, MDTextToken{"formatting"}, MDInlineFormatToken{INLINE_FORMAT_END, "_"},
				MDSimpleToken{NL},
				MDLeadingSpaceToken{1}, MDHeaderIndicToken{1, "# "}, MDTextToken{"Header with leading space"},
				MDSimpleToken{NL},
				MDLeadingSpaceToken{2}, MDHeaderIndicToken{2, "## "}, MDTextToken{"Header with "}, MDInlineFormatToken{INLINE_FORMAT_START, "**"}, MDTextToken{"formatting"}, MDInlineFormatToken{INLINE_FORMAT_END, "**"},
				MDSimpleToken{NL},
				MDUnorderedListIndicToken{"- "}, MDHeaderIndicToken{1, "# "}, MDTextToken{"List item with header and "}, MDInlineFormatToken{INLINE_FORMAT_START, "**"}, MDTextToken{"formatting"}, MDInlineFormatToken{INLINE_FORMAT_END, "**"},
				MDSimpleToken{NL},
				MDLeadingSpaceToken{1}, MDOrderedListIndicToken{"1. "}, MDTextToken{"Same with "}, MDInlineFormatToken{INLINE_FORMAT_START, "**"}, MDTextToken{"ordered"}, MDInlineFormatToken{INLINE_FORMAT_END, "**"}, MDTextToken{" list"},
				MDSimpleToken{NL},
				MDHeaderIndicToken{1, "# "}, MDUnorderedListIndicToken{"- "}, MDTextToken{"Header ignores rest of list but "}, MDInlineFormatToken{INLINE_FORMAT_START, "*"}, MDTextToken{"not formatting"}, MDInlineFormatToken{INLINE_FORMAT_END, "*"},
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
