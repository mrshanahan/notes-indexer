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
	NONE MDTokenType = iota

	NL

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

	SPECIAL_CHAR_ESCAPE
)

var mdTokenTypeName map[MDTokenType]string = map[MDTokenType]string{
	NONE:                     "NONE",
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
	SPECIAL_CHAR_ESCAPE:      "SPECIAL_CHAR_ESCAPE",
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

func (t MDSimpleToken) String() string {
	return t.Type.String()
}

type MDInlineFormatToken struct {
	Type    MDTokenType
	Content string
}

func (t MDInlineFormatToken) GetType() MDTokenType { return t.Type }

func (t MDInlineFormatToken) String() string {
	return fmt.Sprintf("%v(%s)", t.Type, t.Content)
}

type MDTextToken struct {
	Content string
}

func (t MDTextToken) GetType() MDTokenType { return TEXT }

func (t MDTextToken) String() string { return fmt.Sprintf("%v(%s)", t.GetType(), t.Content) }

type MDHeaderIndicToken struct {
	Count   int
	Content string
}

func (t MDHeaderIndicToken) GetType() MDTokenType { return HEADER_INDIC }

func (t MDHeaderIndicToken) String() string {
	return fmt.Sprintf("HEADER(%s)", strings.Repeat("#", t.Count))
}

type MDOrderedListIndicToken struct {
	Content string
}

func (t MDOrderedListIndicToken) GetType() MDTokenType { return ORDERED_LIST_INDIC }

func (t MDOrderedListIndicToken) String() string { return fmt.Sprintf("LIST(%s)", t.Content) }

type MDUnorderedListIndicToken struct {
	Content string
}

func (t MDUnorderedListIndicToken) GetType() MDTokenType { return UNORDERED_LIST_INDIC }

func (t MDUnorderedListIndicToken) String() string { return fmt.Sprintf("LIST(%s)", t.Content) }

type MDLeadingSpaceToken struct {
	Count int
}

func (t MDLeadingSpaceToken) GetType() MDTokenType { return LEADING_SPACE }

func (t MDLeadingSpaceToken) String() string {
	return fmt.Sprintf("%v(%s)", t.GetType(), strings.Repeat(" ", t.Count))
}

type MDEscapeToken struct {
	Content string
}

func (t MDEscapeToken) GetType() MDTokenType { return SPECIAL_CHAR_ESCAPE }

func (t MDEscapeToken) String() string { return fmt.Sprintf("ESCAPED(%s)", t.Content) }

const (
	INLINE_FORMAT_CHARS string = "*_~`"
	SPECIAL_CHARS       string = `\\\[\]()`
)

func GetSpecialCharTokenType(b byte) MDTokenType {
	switch b {
	case '\\':
		return SPECIAL_CHAR_ESCAPE
	case '[':
		return INLINE_LINK_DESC_START
	case ']':
		return INLINE_LINK_DESC_END
	case '(':
		return INLINE_LINK_URL_START
	case ')':
		return INLINE_LINK_URL_END
	default:
		return NONE
	}
}

// TODO:
// - Links ([...](...))
// - a/b/c, A/B/C, i/ii/iii, etc. for lists
// - Escape at end of line forces new line (instead of concatenating e.g. two paragraphs)
// - Integration of HTML elements
//     - Do we need to do this explicitly, or does this just work?

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
	inlineFormattingPattStr string = fmt.Sprintf("(\\s|^)([%[1]s]+)\\S|\\S([%[1]s]+)(\\s|$)|\\S([%[1]s]+)\\S", INLINE_FORMAT_CHARS)

	inlineCharPatt *regexp.Regexp = regexp.MustCompile(fmt.Sprintf("(\\\\([^\\s])?)|(%s)", inlineFormattingPattStr))

	specialCharGroup int = 1

	specialCharEscapedGroup int = 2

	inlineFormattingGroup int = 3

	inlineFormatStartGroup int = 5

	inlineFormatEndGroup int = 6

	inlineFormatMidGroup int = 8

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
				m := inlineCharPatt.FindSubmatchIndex(bytes[cur:lineEnd])
				if m == nil {
					text := string(bytes[cur:lineEnd])
					tokens = append(tokens, MDTextToken{text})
					cur = lineEnd
				} else if m[specialCharGroup*2] >= 0 {
					startidx, endidx := specialCharGroup*2, specialCharGroup*2+1
					var escaped string
					if m[specialCharEscapedGroup*2] >= 0 {
						escaped = string(bytes[cur+m[specialCharEscapedGroup*2] : cur+m[specialCharEscapedGroup*2+1]])
					}
					escapeToken := MDEscapeToken{escaped}

					// Only create text token if we have non-empty text to add
					if m[startidx] > 0 {
						textToken := MDTextToken{string(bytes[cur : cur+m[startidx]])}
						tokens = append(tokens, textToken)
					}
					tokens = append(tokens, escapeToken)
					cur += m[endidx]
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
