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
	Components []tokenComponent
	M          int
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
	components := []tokenComponent{}
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
			component := tokenComponent{
				Type:      curComponentType,
				Component: token[startOfCurComponent:i],
			}
			components = append(components, component)
			startOfCurComponent = i
			curComponentType = TYPE_VOWEL
		} else if curComponentType == TYPE_VOWEL && consonant {
			component := tokenComponent{
				Type:      curComponentType,
				Component: token[startOfCurComponent:i],
			}
			components = append(components, component)
			startOfCurComponent = i
			curComponentType = TYPE_CONSONANT
		}
	}
	component := tokenComponent{
		Type:      curComponentType,
		Component: token[startOfCurComponent:len(token)],
	}
	components = append(components, component)
	return &porterStructure{
		Token:      token,
		Components: components,
		M:          calculateM(components),
	}
}

func calculateM(components []tokenComponent) int {
	init := 0
	if components[init].Type == TYPE_CONSONANT {
		init = 1
	}
	m := (len(components) - init) / 2
	return m
}

func any(p func(tokenComponent) bool, xs []tokenComponent) bool {
	for _, x := range xs {
		if p(x) {
			return true
		}
	}
	return false
}

func isVowelComponent(t tokenComponent) bool {
	return t.Type == TYPE_VOWEL
}

func (t tokenComponent) isSingleVowel() bool {
	return t.Type == TYPE_VOWEL && len(t.Component) == 1
}

func (t tokenComponent) isSingleConsonant() bool {
	return t.Type == TYPE_CONSONANT && len(t.Component) == 1
}

func (p *porterStructure) matchesMWithSuffix(pred func(m int) bool, suffix string) bool {
	if !p.hasSuffix(suffix) {
		return false
	}
	prefixP := p.replaceSuffix(suffix, "")
	return pred(prefixP.M)
}

func (p *porterStructure) hasRuleWithSuffix(patt, suffix string) bool {
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
		numComponents := len(prefixP.Components)
		endsCvc := numComponents >= 3 &&
			prefixP.Components[numComponents-3].isSingleConsonant() &&
			prefixP.Components[numComponents-2].isSingleVowel() &&
			prefixP.Components[numComponents-1].isSingleConsonant()
		if !endsCvc {
			return false
		}
		lastComponentToken := prefixP.Components[numComponents-1].Component
		notWXY := lastComponentToken != "W" &&
			lastComponentToken != "X" &&
			lastComponentToken != "Y"
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

func Stem(token string) string {
	// TODO: Just statically encode every fucking transition. There's < 100, right?
	//       Can't we just generate those?

	// TODO: Don't need to parse until 1b
	props := parseToken(token)

	// Step 1a
	if props.hasSuffix("sses") {
		props = props.truncateFromRight(2)
	} else if props.hasSuffix("ies") {
		props = props.truncateFromRight(2)
	} else if !props.hasSuffix("ss") && props.hasSuffix("s") {
		props = props.truncateFromRight(1)
	}

	// Step 1b
	matchedEDorING := false
	if props.M > 0 && props.hasSuffix("eed") {
		props = props.truncateFromRight(1)
	} else if props.hasRuleWithSuffix("*v*", "ed") {
		props = props.truncateFromRight(2)
		matchedEDorING = true
	} else if props.hasRuleWithSuffix("*v*", "ing") {
		props = props.truncateFromRight(3)
		matchedEDorING = true
	}

	if matchedEDorING {
		if props.hasSuffix("at") {
			props = props.replaceSuffix("at", "ate")
		} else if props.hasSuffix("bl") {
			props = props.replaceSuffix("bl", "ble")
		} else if props.hasSuffix("iz") {
			props = props.replaceSuffix("iz", "ize")
		} else if props.hasRuleWithSuffix("*d", "") && !(props.hasSuffix("l") || props.hasSuffix("s") || props.hasSuffix("z")) {
			props = props.truncateFromRight(1)
		} else if props.M == 1 && props.hasRuleWithSuffix("*o", "") {
			props = props.replaceSuffix("", "e")
		}
	}

	// Step 1c
	if props.hasRuleWithSuffix("*v*", "y") {
		props = props.replaceSuffix("y", "i")
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

	mGreaterThan0 := func(m int) bool { return m > 0 }
	if props.matchesMWithSuffix(mGreaterThan0, "ational") {
		props = props.replaceSuffix("ational", "ate")
	} else if props.matchesMWithSuffix(mGreaterThan0, "tional") {
		props = props.replaceSuffix("tional", "tion")
	} else if props.matchesMWithSuffix(mGreaterThan0, "enci") {
		props = props.replaceSuffix("enci", "ence")
	} else if props.matchesMWithSuffix(mGreaterThan0, "anci") {
		props = props.replaceSuffix("anci", "ance")
	} else if props.matchesMWithSuffix(mGreaterThan0, "izer") {
		props = props.replaceSuffix("izer", "ize")
	} else if props.matchesMWithSuffix(mGreaterThan0, "abli") {
		props = props.replaceSuffix("abli", "able")
	} else if props.matchesMWithSuffix(mGreaterThan0, "alli") {
		props = props.replaceSuffix("alli", "al")
	} else if props.matchesMWithSuffix(mGreaterThan0, "entli") {
		props = props.replaceSuffix("entli", "ent")
	} else if props.matchesMWithSuffix(mGreaterThan0, "eli") {
		props = props.replaceSuffix("eli", "e")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ousli") {
		props = props.replaceSuffix("ousli", "ous")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ization") {
		props = props.replaceSuffix("ization", "ize")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ation") {
		props = props.replaceSuffix("ation", "ate")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ator") {
		props = props.replaceSuffix("ator", "ate")
	} else if props.matchesMWithSuffix(mGreaterThan0, "alism") {
		props = props.replaceSuffix("alism", "al")
	} else if props.matchesMWithSuffix(mGreaterThan0, "iveness") {
		props = props.replaceSuffix("iveness", "ive")
	} else if props.matchesMWithSuffix(mGreaterThan0, "fulness") {
		props = props.replaceSuffix("fulness", "ful")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ousness") {
		props = props.replaceSuffix("ousness", "ous")
	} else if props.matchesMWithSuffix(mGreaterThan0, "aliti") {
		props = props.replaceSuffix("aliti", "al")
	} else if props.matchesMWithSuffix(mGreaterThan0, "iviti") {
		props = props.replaceSuffix("iviti", "ive")
	} else if props.matchesMWithSuffix(mGreaterThan0, "biliti") {
		props = props.replaceSuffix("biliti", "ble")
	}

	// Step 3
	if props.matchesMWithSuffix(mGreaterThan0, "icate") {
		props = props.replaceSuffix("icate", "")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ative") {
		props = props.replaceSuffix("ative", "")
	} else if props.matchesMWithSuffix(mGreaterThan0, "alize") {
		props = props.replaceSuffix("alize", "al")
	} else if props.matchesMWithSuffix(mGreaterThan0, "iciti") {
		props = props.replaceSuffix("iciti", "ic")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ical") {
		props = props.replaceSuffix("ical", "ic")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ful") {
		props = props.replaceSuffix("ful", "")
	} else if props.matchesMWithSuffix(mGreaterThan0, "ness") {
		props = props.replaceSuffix("ness", "")
	}

	// Step 4
	mGreaterThan1 := func(m int) bool { return m > 1 }
	if props.matchesMWithSuffix(mGreaterThan1, "al") {
		props = props.replaceSuffix("al", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ance") {
		props = props.replaceSuffix("ance", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ence") {
		props = props.replaceSuffix("ence", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "er") {
		props = props.replaceSuffix("er", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ic") {
		props = props.replaceSuffix("ic", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "able") {
		props = props.replaceSuffix("able", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ible") {
		props = props.replaceSuffix("ible", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ant") {
		props = props.replaceSuffix("ant", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ement") {
		props = props.replaceSuffix("ement", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ment") {
		props = props.replaceSuffix("ment", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ent") {
		props = props.replaceSuffix("ent", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ion") && // TODO: Do something about this guy
		(props.hasRuleWithSuffix("*S", "ion") || props.hasRuleWithSuffix("*T", "ion")) {
		props = props.replaceSuffix("ion", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ou") {
		props = props.replaceSuffix("ou", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ism") {
		props = props.replaceSuffix("ism", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ate") {
		props = props.replaceSuffix("ate", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "iti") {
		props = props.replaceSuffix("iti", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ous") {
		props = props.replaceSuffix("ous", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ive") {
		props = props.replaceSuffix("ive", "")
	} else if props.matchesMWithSuffix(mGreaterThan1, "ize") {
		props = props.replaceSuffix("ize", "")
	}

	// Step 5a
	mEquals1 := func(m int) bool { return m == 1 }
	if props.matchesMWithSuffix(mGreaterThan1, "e") {
		props = props.replaceSuffix("e", "")
	} else if props.matchesMWithSuffix(mEquals1, "e") && !props.hasRuleWithSuffix("*o", "e") {
		props = props.replaceSuffix("e", "")
	}

	// Step 5b
	if props.M > 1 && props.hasRuleWithSuffix("*d", "l") {
		props = props.truncateFromRight(1)
	}
	return props.Token
}
