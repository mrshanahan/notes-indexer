package markdown

import (
	"fmt"
	"regexp"
	"strings"
)

// NL

// TEXT(...)

// LEADING_SPACE(n)

// UNORDERED_LIST_INDIC
// ORDERED_LIST_INDIC
// HEADER_INDIC(N)
// EXPLICIT_CODEBLOCK_INDIC

// INLINE_BOLD_INDIC
// INLINE_ITALICS_INDIC
// INLINE_UNDERLINE_INDIC
// INLINE_STRIKETHROUGH_INDIC
// INLINE_CODE_INDIC

// INLINE_LINK_DESC_START
// INLINE_LINK_DESC_END
// INLINE_LINK_URL_START
// INLINE_LINK_URL_END

type MDTokenType int

const (
	NL MDTokenType = iota

	TEXT // Requires content

	LEADING_SPACE // Requires (n)

	UNORDERED_LIST_INDIC
	ORDERED_LIST_INDIC
	HEADER_INDIC // Requires (n)
	EXPLICIT_CODEBLOCK_INDIC

	INLINE_FORMAT_START
	INLINE_FORMAT_MID
	INLINE_FORMAT_END
	// INLINE_ASTERISK_START
	// INLINE_ASTERISK_MID
	// INLINE_ASTERISK_END
	// INLINE_UNDERSCORE_START
	// INLINE_UNDERSCORE_MID
	// INLINE_UNDERSCORE_END
	// INLINE_TILDE_START
	// INLINE_TILDE_MID
	// INLINE_TILDE_END
	// INLINE_BACKTICK_START
	// INLINE_BACKTICK_MID
	// INLINE_BACKTICK_END

	INLINE_LINK_DESC_START
	INLINE_LINK_DESC_END
	INLINE_LINK_URL_START
	INLINE_LINK_URL_END
)

var mdTokenTypeName map[MDTokenType]string = map[MDTokenType]string{
	NL:                       "NL",
	TEXT:                     "TEXT",
	LEADING_SPACE:            "LEADING_SPACE",
	UNORDERED_LIST_INDIC:     "UNORDERED_LIST_INDIC",
	ORDERED_LIST_INDIC:       "ORDERED_LIST_INDIC",
	HEADER_INDIC:             "HEADER_INDIC",
	EXPLICIT_CODEBLOCK_INDIC: "EXPLICIT_CODEBLOCK_INDIC",
	INLINE_FORMAT_START:      "INLINE_FORMAT_START",
	INLINE_FORMAT_MID:        "INLINE_FORMAT_MID",
	INLINE_FORMAT_END:        "INLINE_FORMAT_END",
	INLINE_LINK_DESC_START:   "INLINE_LINK_DESC_START",
	INLINE_LINK_DESC_END:     "INLINE_LINK_DESC_END",
	INLINE_LINK_URL_START:    "INLINE_LINK_URL_START",
	INLINE_LINK_URL_END:      "INLINE_LINK_URL_END",
}

func (t MDTokenType) String() string {
	return mdTokenTypeName[t]
}

type MDToken interface {
	GetType() MDTokenType
}

type MDSimpleToken struct {
	Type MDTokenType
}

// TODO: Should these be on the pointer instead of the plain type? Should it depend on the token type?

func (t MDSimpleToken) GetType() MDTokenType { return t.Type }

type MDInlineFormatToken struct {
	Type    MDTokenType
	Content string
}

func (t MDInlineFormatToken) GetType() MDTokenType { return t.Type }

type MDTextToken struct {
	Content string
}

func (t MDTextToken) GetType() MDTokenType { return TEXT }

type MDHeaderIndicToken struct {
	Count   int
	Content string
}

func (t MDHeaderIndicToken) GetType() MDTokenType { return HEADER_INDIC }

type MDOrderedListIndicToken struct {
	Content string
}

func (t MDOrderedListIndicToken) GetType() MDTokenType { return ORDERED_LIST_INDIC }

type MDUnorderedListIndicToken struct {
	Content string
}

func (t MDUnorderedListIndicToken) GetType() MDTokenType { return UNORDERED_LIST_INDIC }

type MDLeadingSpaceToken struct {
	Count int
}

func (t MDLeadingSpaceToken) GetType() MDTokenType { return LEADING_SPACE }

const (
	INLINE_FORMAT_CHARS string = "*_~`"
)

var (
	// A little tricky - these three groups are start/end/mid respectively.
	// For example, in the following example:
	//
	//     This is *formatting that w~o~ill be italicized.*
	//
	// The following items each type:
	// - start: This is >*formatting that ...
	// - end:   ... be italicized.*<
	// - mid:   ... that w>~i~<ll be ...
	//
	// mid has to be a different category because it could play the role of
	// either a start or end indicator, whereas the other two can only play
	// start/end roles:
	//
	//     The first *mid will work as an en*d, and th*e second works* as a start.
	//
	// If an end indicator is followed by a start indicator they will not form
	// an inline formatting block:
	//
	//     This will* not be formatted in any *way.
	//
	// NB: We have to be careful with the mid group, as it matches non-zero-width
	// assertions by necessity - we don't want to say a '*' is mid only because
	// it is followed by another '*'.
	inlineFormattingPatt *regexp.Regexp = regexp.MustCompile(
		fmt.Sprintf("(\\s|^)([%[1]s]+)\\S|\\S([%[1]s]+)(\\s|$)|\\S([%[1]s]+)\\S", INLINE_FORMAT_CHARS))

	inlineFormatStartGroup int = 2

	inlineFormatEndGroup int = 3

	inlineFormatMidGroup int = 5

	unorderedListIndicPatt *regexp.Regexp = regexp.MustCompile(`^-\s`)

	orderedListIndicPatt *regexp.Regexp = regexp.MustCompile(`^\d+\.\s`)

	headerIndicPatt *regexp.Regexp = regexp.MustCompile(`^(#+)\s`)

	leadingWhitespacePatt *regexp.Regexp = regexp.MustCompile(`^\s+`)

	endOfLinePatt *regexp.Regexp = regexp.MustCompile("(?m)^.*$")
)

func Lex(text string) []MDToken {
	// Standard preprocessing - CRLF-to-LF for simplicity, but the tabs-to-four-spaces
	// appears to be standard behavior for Markdown parsers. We'll follow it here
	// as well b/c it makes things sooo much easier.
	//
	// TODO: Should we also replace '\r' with '\n'? That's standard for other formats but I don't
	//       know if it also is for Markdown.
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\t", "    ")

	bytes := []byte(text)
	tokens := []MDToken{} // TODO: better initial capacity

	cur := 0
	for cur < len(bytes) {
		if bytes[cur] == '\n' {
			token := MDSimpleToken{NL}
			tokens = append(tokens, token)
			cur++
		} else if m := leadingWhitespacePatt.FindIndex(bytes[cur:]); m != nil {
			token := MDLeadingSpaceToken{m[1] - m[0]}
			tokens = append(tokens, token)
			cur += m[1]
		} else if m := unorderedListIndicPatt.FindIndex(bytes[cur:]); m != nil {
			token := MDUnorderedListIndicToken{string(bytes[cur+m[0] : cur+m[1]])}
			tokens = append(tokens, token)
			cur += m[1]
		} else if m := orderedListIndicPatt.FindIndex(bytes[cur:]); m != nil {
			token := MDOrderedListIndicToken{string(bytes[cur+m[0] : cur+m[1]])}
			tokens = append(tokens, token)
			cur += m[1]
		} else if m := headerIndicPatt.FindSubmatchIndex(bytes[cur:]); m != nil {
			token := MDHeaderIndicToken{m[3] - m[2], string(bytes[cur+m[0] : cur+m[1]])}
			tokens = append(tokens, token)
			cur += m[1]
		} else {
			// otherwise, paragraph block until end of line
			m := endOfLinePatt.FindIndex(bytes[cur:])
			lineEnd := cur + m[1]

			for cur < lineEnd {
				m := inlineFormattingPatt.FindSubmatchIndex(bytes[cur:lineEnd])
				if m == nil {
					text := string(bytes[cur:lineEnd])
					tokens = append(tokens, MDTextToken{text})
					cur = lineEnd
				} else {
					var fmtToken MDInlineFormatToken
					var textToken MDTextToken
					if m[inlineFormatStartGroup*2] >= 0 { // start
						startidx, endidx := inlineFormatStartGroup*2, inlineFormatStartGroup*2+1
						chr := string(bytes[cur+m[startidx] : cur+m[endidx]])
						fmtToken = MDInlineFormatToken{INLINE_FORMAT_START, chr}
						textToken = MDTextToken{string(bytes[cur : cur+m[startidx]])}
						cur += m[endidx]
					} else if m[inlineFormatEndGroup*2] >= 0 { // end
						startidx, endidx := inlineFormatEndGroup*2, inlineFormatEndGroup*2+1
						chr := string(bytes[cur+m[startidx] : cur+m[endidx]])
						fmtToken = MDInlineFormatToken{INLINE_FORMAT_END, chr}
						textToken = MDTextToken{string(bytes[cur : cur+m[startidx]])}
						cur += m[endidx]
					} else { // mid
						startidx, endidx := inlineFormatMidGroup*2, inlineFormatMidGroup*2+1
						chr := string(bytes[cur+m[startidx] : cur+m[endidx]])
						fmtToken = MDInlineFormatToken{INLINE_FORMAT_MID, chr}
						textToken = MDTextToken{string(bytes[cur : cur+m[startidx]])}
						cur += m[endidx]
					}

					tokens = append(tokens, textToken)
					tokens = append(tokens, fmtToken)
				}
			}
		}
	}

	return tokens
}

var LINE_SPLIT_PATT *regexp.Regexp = regexp.MustCompile("\r?\n")

func splitLines(text string) []string {
	indices := LINE_SPLIT_PATT.FindAllStringIndex(text, -1)
	lines := []string{}
	cur := 0
	for _, xy := range indices {
		x, y := xy[0], xy[1]
		lines = append(lines, text[cur:x])
		cur = y
	}
	lines = append(lines, text[cur:])

	return lines
}
