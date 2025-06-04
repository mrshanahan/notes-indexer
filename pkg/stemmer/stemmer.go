package stemmer

import (
	"fmt"
	"regexp"
	s "strings"
)

const (
	TYPE_CONSONANT = true
	TYPE_VOWEL     = false
)

type tokenComponent struct {
	Type      bool
	Component string
}

type porterStructure struct {
	Token      string
	Components []*tokenComponent
	M          int
	Consonants []bool
}

func isConsonant(i int, token string, consonants []bool) bool {
	c := token[i]
	if c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u' {
		return false
	}
	if c == 'y' && i > 0 && consonants[i-1] {
		return false
	}
	return true
}

func parseToken(token string) *porterStructure {
	consonants := make([]bool, len(token))
	components := []*tokenComponent{}
	curComponentType := TYPE_CONSONANT
	startOfCurComponent := 0
	for i, _ := range token {
		consonant := isConsonant(i, token, consonants)
		consonants[i] = consonant
		if i == 0 {
			if consonant {
				curComponentType = TYPE_CONSONANT
			} else {
				curComponentType = TYPE_VOWEL
			}
		} else if curComponentType == TYPE_CONSONANT && !consonant {
			component := &tokenComponent{
				Type:      curComponentType,
				Component: token[startOfCurComponent:i],
			}
			components = append(components, component)
			startOfCurComponent = i
			curComponentType = TYPE_VOWEL
		} else if curComponentType == TYPE_VOWEL && consonant {
			component := &tokenComponent{
				Type:      curComponentType,
				Component: token[startOfCurComponent:i],
			}
			components = append(components, component)
			startOfCurComponent = i
			curComponentType = TYPE_CONSONANT
		}
	}
	component := &tokenComponent{
		Type:      curComponentType,
		Component: token[startOfCurComponent:len(token)],
	}
	components = append(components, component)
	return &porterStructure{
		Token:      token,
		Components: components,
		M:          calculateM(components),
		Consonants: consonants,
	}
}

func calculateM(components []*tokenComponent) int {
	init := 0
	if components[init].Type == TYPE_CONSONANT {
		init = 1
	}
	m := (len(components) - init) / 2
	return m
}

func any(p func(*tokenComponent) bool, xs []*tokenComponent) bool {
	for _, x := range xs {
		if p(x) {
			return true
		}
	}
	return false
}

func isVowelComponent(t *tokenComponent) bool {
	return t.Type == TYPE_VOWEL
}

func (p *porterStructure) matchesMWithSuffix(pred func(m int) bool, suffix string) bool {
	if !p.hasSuffix(suffix) {
		return false
	}
	prefixP := p.replaceSuffix(suffix, "")
	return pred(prefixP.M)
}

func (p *porterStructure) matchesRuleWithSuffix(patt, suffix string) bool {
	if !p.hasSuffix(suffix) {
		return false
	}
	prefixP := p.replaceSuffix(suffix, "")
	if patt == "*v*" { // the stems contains a vowel.
		return any(isVowelComponent, prefixP.Components)
	} else if patt == "*d" { // the stem ends with a double consonant.
		lastComponent := prefixP.Components[len(prefixP.Components)-1]
		componentToken := lastComponent.Component
		return lastComponent.Type == TYPE_CONSONANT &&
			len(componentToken) > 1 &&
			componentToken[len(componentToken)-1] == componentToken[len(componentToken)-2]
	} else if patt == "*o" { // the stem ends cvc, where the second c is not W, X, or Y (e.g. -WIL, -HOP).
		numPrefixLetters := len(prefixP.Token)
		endsCvc := numPrefixLetters >= 3 &&
			prefixP.Consonants[numPrefixLetters-3] == TYPE_CONSONANT &&
			prefixP.Consonants[numPrefixLetters-2] == TYPE_VOWEL &&
			prefixP.Consonants[numPrefixLetters-1] == TYPE_CONSONANT
		if !endsCvc {
			return false
		}
		lastPrefixLetter := prefixP.Token[numPrefixLetters-1]
		notWXY := lastPrefixLetter != 'w' &&
			lastPrefixLetter != 'x' &&
			lastPrefixLetter != 'y'
		return notWXY
	} else if len(patt) == 2 && patt[0] == '*' { // any pattern of the form *S, *L, etc.
		letterMatch := s.ToLower(patt)[1]
		return prefixP.Token[len(prefixP.Token)-1] == letterMatch
	} else {
		panic(fmt.Sprintf("Invalid Porter stemmer pattern: %s. Check with library owner.", patt))
	}
}

func (p *porterStructure) hasSuffix(suffix string) bool {
	return s.HasSuffix(p.Token, suffix)
}

func (p *porterStructure) truncateFromRight(n int) *porterStructure {
	stemmed := p.Token[0:(len(p.Token) - n)]
	// TODO: Don't do a full reparse every time
	return parseToken(stemmed)
}

func (p *porterStructure) replaceSuffix(before string, after string) *porterStructure {
	// TODO: Statically compile all necessary patterns
	patt := regexp.MustCompile(before + "$")
	stemmed := patt.ReplaceAllString(p.Token, after)
	// TODO: Don't do a full reparse every time
	return parseToken(stemmed)
}

func (p *porterStructure) replaceSuffixIfMatchesM(pred func(m int) bool, before, after string) (*porterStructure, bool) {
	prefixP := p.replaceSuffix(before, "")
	if pred(prefixP.M) {
		return p.replaceSuffix(before, after), true
	}
	return p, false
}

func (p *porterStructure) replaceSuffixIfMatchesRule(patt, before, after string) (*porterStructure, bool) {
	if p.matchesRuleWithSuffix(patt, before) {
		return p.replaceSuffix(before, after), true
	}
	return p, false
}

func Stem(token string) string {
	// TODO: Just statically encode every fucking transition. There's < 100, right?
	//       Can't we just generate those?

	// This is a standard change from the published algorithm: don't touch words of length 1 or 2
	if len(token) <= 2 {
		return token
	}

	// TODO: Don't need to parse until 1b
	props := parseToken(token)

	mGreaterThan0 := func(m int) bool { return m > 0 }
	mEquals1 := func(m int) bool { return m == 1 }
	mGreaterThan1 := func(m int) bool { return m > 1 }

	// Step 1a
	if props.hasSuffix("sses") {
		props = props.replaceSuffix("sses", "ss")
	} else if props.hasSuffix("ies") {
		props = props.replaceSuffix("ies", "i")
	} else if !props.hasSuffix("ss") && props.hasSuffix("s") {
		props = props.replaceSuffix("s", "")
	}

	// Step 1b
	matchedEDorING := false
	if props.hasSuffix("eed") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "eed", "ee")
	} else if props.hasSuffix("ed") {
		props, matchedEDorING = props.replaceSuffixIfMatchesRule("*v*", "ed", "")
	} else if props.hasSuffix("ing") {
		props, matchedEDorING = props.replaceSuffixIfMatchesRule("*v*", "ing", "")
	}

	if matchedEDorING {
		if props.hasSuffix("at") {
			props = props.replaceSuffix("at", "ate")
		} else if props.hasSuffix("bl") {
			props = props.replaceSuffix("bl", "ble")
		} else if props.hasSuffix("iz") {
			props = props.replaceSuffix("iz", "ize")
		} else if props.matchesRuleWithSuffix("*d", "") && !(props.hasSuffix("l") || props.hasSuffix("s") || props.hasSuffix("z")) {
			props = props.truncateFromRight(1)
		} else if props.matchesMWithSuffix(mEquals1, "") && props.matchesRuleWithSuffix("*o", "") {
			props = props.replaceSuffix("", "e")
		}
	}

	// Step 1c
	if props.hasSuffix("y") {
		props, _ = props.replaceSuffixIfMatchesRule("*v*", "y", "i")
	}

	// Step 2
	// TODO: Lots of availability for truncation in here
	// TODO: Note from Porter Stemmer def.txt:
	//     The test for the string S1 can be made fast by doing a program switch on
	//     the penultimate letter of the word being tested. This gives a fairly even
	//     breakdown of the possible values of the string S1. It will be seen in fact
	//     that the S1-strings in step 2 are presented here in the alphabetical order
	//     of their penultimate letter. Similar techniques may be applied in the other
	//     steps.

	if props.hasSuffix("ational") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ational", "ate")
	} else if props.hasSuffix("tional") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "tional", "tion")
	} else if props.hasSuffix("enci") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "enci", "ence")
	} else if props.hasSuffix("anci") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "anci", "ance")
	} else if props.hasSuffix("izer") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "izer", "ize")
	} else if props.hasSuffix("bli") {
		// This is a standard change from the published algorithm: bli -> bl instead of abli -> able
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "bli", "ble")
	} else if props.hasSuffix("alli") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "alli", "al")
	} else if props.hasSuffix("entli") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "entli", "ent")
	} else if props.hasSuffix("eli") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "eli", "e")
	} else if props.hasSuffix("ousli") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ousli", "ous")
	} else if props.hasSuffix("ization") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ization", "ize")
	} else if props.hasSuffix("ation") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ation", "ate")
	} else if props.hasSuffix("ator") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ator", "ate")
	} else if props.hasSuffix("alism") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "alism", "al")
	} else if props.hasSuffix("iveness") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "iveness", "ive")
	} else if props.hasSuffix("fulness") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "fulness", "ful")
	} else if props.hasSuffix("ousness") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ousness", "ous")
	} else if props.hasSuffix("aliti") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "aliti", "al")
	} else if props.hasSuffix("iviti") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "iviti", "ive")
	} else if props.hasSuffix("biliti") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "biliti", "ble")
	} else if props.hasSuffix("logi") {
		// This is a standard change from the published algorithm: extra rule to account for "-ology" words
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "logi", "log")
	}

	// Step 3
	if props.hasSuffix("icate") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "icate", "ic")
	} else if props.hasSuffix("ative") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ative", "")
	} else if props.hasSuffix("alize") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "alize", "al")
	} else if props.hasSuffix("iciti") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "iciti", "ic")
	} else if props.hasSuffix("ical") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ical", "ic")
	} else if props.hasSuffix("ful") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ful", "")
	} else if props.hasSuffix("ness") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan0, "ness", "")
	}

	// Step 4
	if props.hasSuffix("al") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "al", "")
	} else if props.hasSuffix("ance") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ance", "")
	} else if props.hasSuffix("ence") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ence", "")
	} else if props.hasSuffix("er") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "er", "")
	} else if props.hasSuffix("ic") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ic", "")
	} else if props.hasSuffix("able") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "able", "")
	} else if props.hasSuffix("ible") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ible", "")
	} else if props.hasSuffix("ant") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ant", "")
	} else if props.hasSuffix("ement") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ement", "")
	} else if props.hasSuffix("ment") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ment", "")
	} else if props.hasSuffix("ent") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ent", "")
	} else if props.hasSuffix("ion") { // TODO: Clean this one up a bit
		if props.matchesMWithSuffix(mGreaterThan1, "ion") &&
			(props.matchesRuleWithSuffix("*S", "ion") || props.matchesRuleWithSuffix("*T", "ion")) {
			props = props.replaceSuffix("ion", "")
		}
	} else if props.hasSuffix("ou") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ou", "")
	} else if props.hasSuffix("ism") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ism", "")
	} else if props.hasSuffix("ate") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ate", "")
	} else if props.hasSuffix("iti") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "iti", "")
	} else if props.hasSuffix("ous") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ous", "")
	} else if props.hasSuffix("ive") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ive", "")
	} else if props.hasSuffix("ize") {
		props, _ = props.replaceSuffixIfMatchesM(mGreaterThan1, "ize", "")
	}

	// Step 5a
	if props.hasSuffix("e") {
		if props.matchesMWithSuffix(mGreaterThan1, "e") || (props.matchesMWithSuffix(mEquals1, "e") && !props.matchesRuleWithSuffix("*o", "e")) {
			props = props.replaceSuffix("e", "")
		}
	}

	// Step 5b
	if props.matchesMWithSuffix(mGreaterThan1, "") && props.matchesRuleWithSuffix("*d", "") && props.matchesRuleWithSuffix("*L", "") {
		props = props.truncateFromRight(1)
	}

	return props.Token
}
