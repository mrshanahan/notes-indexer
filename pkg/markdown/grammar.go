package markdown

import (
	"fmt"
	"strings"

	"mrshanahan.com/notes-indexer/internal/util"
)

type mdGrammarRuleType int

const (
	RULE_TOKEN mdGrammarRuleType = iota
	RULE_OPTION
	RULE_UNION
	RULE_PRODUCT
	RULE_ZERO_PLUS
	RULE_ONE_PLUS
)

var mdGrammarRuleTypeNames map[mdGrammarRuleType]string = map[mdGrammarRuleType]string{
	RULE_TOKEN:     "TOKEN",
	RULE_OPTION:    "OPTIONS",
	RULE_UNION:     "UNION",
	RULE_PRODUCT:   "PRODUCT",
	RULE_ZERO_PLUS: "ZERO_PLUS",
	RULE_ONE_PLUS:  "ONE_PLUS",
}

func (t mdGrammarRuleType) String() string {
	return mdGrammarRuleTypeNames[t]
}

type mdGrammarRule interface {
	GetType() mdGrammarRuleType
}

type mdTokenRule struct {
	tokenType MDTokenType
}

func (r *mdTokenRule) GetType() mdGrammarRuleType { return RULE_TOKEN }

func (r *mdTokenRule) String() string { return fmt.Sprintf("TOKEN(%v)", r.tokenType) }

type mdOptionRule struct {
	term mdGrammarRule
}

func (r *mdOptionRule) GetType() mdGrammarRuleType { return RULE_OPTION }

func (r *mdOptionRule) String() string { return fmt.Sprintf("(%v)?", r.term) }

type mdMultiRule struct {
	typ   mdGrammarRuleType
	terms []mdGrammarRule
}

func (r *mdMultiRule) GetType() mdGrammarRuleType { return r.typ }

func (r *mdMultiRule) String() string {
	strTerms := util.Map(r.terms, func(t mdGrammarRule) string { return fmt.Sprintf("%v", t) })
	if r.typ == RULE_PRODUCT {
		return strings.Join(strTerms, " ")
	}
	return strings.Join(strTerms, "|")
}

type mdRepeatRule struct {
	typ  mdGrammarRuleType
	term mdGrammarRule
}

func (r *mdRepeatRule) GetType() mdGrammarRuleType { return r.typ }

func (r *mdRepeatRule) String() string {
	if r.typ == RULE_ZERO_PLUS {
		return fmt.Sprintf("(%v)*", r.term)
	}
	return fmt.Sprintf("(%v)+", r.term)
}

func parseGrammarRule(rule string) (mdGrammarRule, error) {
	return nil, nil
}
