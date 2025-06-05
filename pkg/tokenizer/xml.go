package tokenizer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	PATT_CHAR string = "[\\x09\\x0A\\x0D]|[\\x20-\\x{D7FF}]|[\\x{E000}-\\x{FFFD}]|[\\x{00010000}-\\x{0010FFFF}]"
	PATT_S    string = "[\\x20\\x09\\x0D\\x0A]+"

	// Translated directly from: https://www.w3.org/TR/xml/#NT-Name
	PATT_NAMESTARTCHAR string = "[:A-Za-z_]|[\\xC0-\\xD6]|[\\xD8-\\xF6]|[\\xF8-\\x{02FF}]|[\\x{0370}-\\x{037D}]|[\\x{037F}-\\x{1FFF}]|[\\x{200C}-\\x{200D}]|[\\x{2070}-\\x{218F}]|[\\x{2C00}-\\x{2FEF}]|[\\x{3001}-\\x{D7FF}]|[\\x{F900}-\\x{FDCF}]|[\\x{FDF0}-\\x{FFFD}]|[\\x{00010000}-\\x{000EFFFF}]"
	PATT_NAMECHAR      string = PATT_NAMESTARTCHAR + "|[\\-.0-9]|\\xB7|[\\x{0300}-\\x{036F}]|[\\x{203F}-\\x{2040}]"
	PATT_NAME          string = "(" + PATT_NAMESTARTCHAR + ")(" + PATT_NAMECHAR + ")*"
	PATT_NMTOKEN       string = "(" + PATT_NAMECHAR + ")+"
	PATT_NMTOKENS      string = "(" + PATT_NMTOKEN + ")" + "(" + "[\x20]" + PATT_NMTOKEN + ")*"

	PATT_PEREFERENCE string = "%(" + PATT_NAME + ");"
	PATT_REFERENCE   string = "(" + PATT_ENTITYREF + ")|(" + PATT_CHARREF + ")"
	PATT_ENTITYVALUE string = "\"" + "([^&%\"]" + "|" + PATT_PEREFERENCE + "|" + PATT_REFERENCE + ")*" + "\""
	PATT_ATTVALUE    string = "(\"" + "([^<&\"]" + "|" + PATT_REFERENCE + ")*" + "\")" + "|" +
		"('" + "([^<&']" + "|" + PATT_REFERENCE + ")*" + "')"
	PATT_SYSTEMLITERAL string = "" // ::= ('"' [^"]* '"') | ("'" [^']* "'")
	PATT_PUBIDLITERAL  string = "" // ::= '"' PubidChar* '"' | "'" (PubidChar - "'")* "'"
	PATT_PUBIDCHAR     string = "" // ::= #x20 | #xD | #xA | [a-zA-Z0-9] | [-'()+,./:=?;!*#@$_%]

	PATT_CHARDATA string = "(!?\\]\\]>)[^<&]*" // ::= [^<&]* - ([^<&]* ']]>' [^<&]*)

	PATT_EQ string = "(" + PATT_S + ")?" + "=" + "(" + PATT_S + ")?"

	PATT_ATTRIBUTE    string = "(" + PATT_NAME + ")(" + PATT_EQ + ")(" + PATT_ATTVALUE + ")"
	PATT_EMPTYELEMTAG string = "<" + ("(" + PATT_NAME + ")") + ("((" + PATT_S + ")(" + PATT_ATTRIBUTE + "))*") + ("(" + PATT_S + ")?") + "/>" // '<' Name (S Attribute)* S? '/>'
	PATT_STAG         string = "<" + ("(" + PATT_NAME + ")") + ("((" + PATT_S + ")(" + PATT_ATTRIBUTE + "))*") + ("(" + PATT_S + ")?") + ">"  // '<' Name (S Attribute)* S? '>'
	PATT_ETAG         string = "</" + ("(" + PATT_NAME + ")") + ("(" + PATT_S + ")?") + ">"                                                   // '</' Name S? '>'

	PATT_CHARREF   string = "&#x?[0-9a-fA-F]+;" // Matches more than valid char refs (e.g. "&#FF;") but we don't really care
	PATT_ENTITYREF string = "&(" + PATT_NAME + ");"

	// NB: This is less restrictive than the given pattern, which attempts to ensure that a
	//     comment does not end in '--->'. It appears that many actually-existing XML/HTML
	//     parsers don't care about this, so we won't either.
	//     Also, it uses the non-greedy match on characters, which appears to be how other
	//     actually-existing parsers work - i.e. "<!-- bar --><foo/><!-- baz -->" would be
	//     one giant comment under the greedy pattern, rather than two comments and a valid
	//     empty element (as it should be) under the non-greedy pattern.
	PATT_COMMENT string = "<!--" + ("(" + PATT_CHAR + ")*?") + "-->" // '<!--' ((Char - '-') | ('-' (Char - '-')))* '-->'

	// PATT_ELEMENT      string = ""

	PATT_XMLENT string = "^((" + PATT_STAG + ")|(" +
		PATT_ETAG + ")|(" + PATT_COMMENT + ")|(" +
		PATT_EMPTYELEMTAG + ")|(" + PATT_CHARREF + ")|(" +
		PATT_ENTITYREF + "))"
)

var (
	// patt_STag         Production = newRegexProduction(PATT_STAG)
	// patt_ETag         Production = newRegexProduction(PATT_ETAG)
	// patt_EmptyElemTag Production = newRegexProduction(PATT_EMPTYELEMTAG)
	// patt_Comment      Production = newRegexProduction(PATT_COMMENT)
	// patt_CharRef      Production = newRegexProduction(PATT_CHARREF)

	// patt_EntityRef *regexp.Regexp = regexp.MustCompile(PATT_ENTITYREF)

	patt_XmlEnt *regexp.Regexp = regexp.MustCompile(PATT_XMLENT)
)

func (t *xmlTokenizer) Tokenize(text string) ([]Token, error) {
	first, cur, tokens, seps := -1, 0, []Token{}, t.separators
	for cur < len(text) {
		if loc := patt_XmlEnt.FindStringIndex(text[cur:]); len(loc) > 0 {
			if first >= 0 {
				tok := strings.ToLower(text[first:cur])
				tokens = append(tokens, Token{tok, TOKEN_TYPE_GENERIC})
				first = -1
			}
			idx, lng := loc[0], loc[1]
			tok := text[cur+idx : cur+idx+lng]
			tokens = append(tokens, Token{tok, TOKEN_TYPE_XML})
			cur += idx + lng
		} else {
			r, rlen := utf8.DecodeRuneInString(text[cur:])
			if r == utf8.RuneError {
				if rlen > 0 {
					return nil, fmt.Errorf("invalid UTF-8 rune at byte %d", cur)
				}
				return nil, fmt.Errorf("unexpected end of input")
			}

			issep := seps[r]
			if !issep && first < 0 {
				first = cur
			} else if issep && first >= 0 {
				tok := strings.ToLower(text[first:cur])
				tokens = append(tokens, Token{tok, TOKEN_TYPE_GENERIC})
				first = -1
			}
			cur += rlen
		}
	}
	if first >= 0 {
		tok := strings.ToLower(text[first:])
		tokens = append(tokens, Token{tok, TOKEN_TYPE_GENERIC})
	}
	return tokens, nil
}

type xmlTokenizer defaultTokenizer

func NewXmlTokenizer() Tokenizer {
	return &xmlTokenizer{
		separators: convertSeparator(DefaultSeparators),
	}
}
