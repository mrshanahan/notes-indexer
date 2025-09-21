package markdown

import "fmt"

type MDSyntaxNodeType int

const (
	SYNTAX_NONE MDSyntaxNodeType = iota
	SYNTAX_PARAGRAPH
)

type MDParagraphFormatNodeType int

const (
	FORMAT_NODE_NONE MDParagraphFormatNodeType = iota
	FORMAT_NODE_BOLD
	FORMAT_NODE_ITALICS
	FORMAT_NODE_UNDERLINE
	FORMAT_NODE_CODE
)

var mdFormatNodeTypeName map[MDParagraphFormatNodeType]string = map[MDParagraphFormatNodeType]string{
	FORMAT_NODE_NONE:      "FMT_NONE",
	FORMAT_NODE_BOLD:      "FMT_BOLD",
	FORMAT_NODE_ITALICS:   "FMT_ITALICS",
	FORMAT_NODE_UNDERLINE: "FMT_UNDERLINE",
	FORMAT_NODE_CODE:      "FMT_CODE",
}

func (t MDParagraphFormatNodeType) String() string { return mdFormatNodeTypeName[t] }

type MDSyntaxNode interface {
	GetType() MDSyntaxNodeType
}

type MDParagraph struct {
	Content []MDParagraphFormatNode
}

func (n MDParagraph) String() string { return fmt.Sprintf("P(%v)", n.Content) }

func (n *MDParagraph) GetType() MDSyntaxNodeType { return SYNTAX_PARAGRAPH }

type MDParagraphFormatNode interface {
	GetFormatNodeType() MDParagraphFormatNodeType
}

type MDInlineFormatNode struct {
	Type    MDParagraphFormatNodeType
	Content MDParagraphFormatNode
}

func (n MDInlineFormatNode) String() string { return fmt.Sprintf("%v(%v)", n.Type, n.Content) }

func (n MDInlineFormatNode) GetFormatNodeType() MDParagraphFormatNodeType { return n.Type }

type MDTextFormatNode struct {
	Content string
}

func (n MDTextFormatNode) String() string {
	return fmt.Sprintf("TEXT(%s)", string(n.Content))
}

func (n MDTextFormatNode) GetFormatNodeType() MDParagraphFormatNodeType { return FORMAT_NODE_NONE }

type MDSyntaxTree struct {
	Children []MDSyntaxNode
}

func (t MDSyntaxTree) String() string { return fmt.Sprintf("TREE(%v)", t.Children) }

type mdParseState struct {
	curNode           MDSyntaxNode
	curProcessingType MDSyntaxNodeType
	prevWasNL         bool
}

func initState() mdParseState {
	return mdParseState{
		curNode:           nil,
		curProcessingType: SYNTAX_NONE,
		prevWasNL:         false,
	}
}

func tokenTypeMismatchError(t MDTokenType) error {
	return fmt.Errorf("invalid syntax tree - token of type %v could not be converted as such", t)
}

func invalidTokenError(t MDTokenType) error {
	return fmt.Errorf("invalid token type: %v", t)
}

func unknownTokenTypeError(t MDTokenType) error {
	return fmt.Errorf("unknown token type: %v", t)
}

// func Parse(tokens []MDToken) (MDSyntaxTree, error) {
// 	ns := []MDSyntaxNode{}
// 	state := initState()
// 	for _, t := range tokens {
// 		if state.curProcessingType == SYNTAX_NONE {
// 			typ := t.GetType()
// 			switch typ {
// 			case TOKEN_NL:
// 				continue
// 			case TOKEN_TEXT:
// 				t, ok := t.(MDTextToken)
// 				if !ok {
// 					return MDSyntaxTree{ns}, tokenTypeMismatchError(typ)
// 				}
// 				var p MDSyntaxNode = MDParagraph{[]MDParagraphFormatNode{MDTextFormatNode(t)}}
// 				state.curNode = &p
// 			case TOKEN_LEADING_SPACE:
// 				// TODO: Deal with leading space
// 				continue
// 			case TOKEN_INLINE_FORMAT_START:
// 				// TODO: Deal with inline formatting
// 				continue
// 			case TOKEN_INLINE_FORMAT_MID:
// 			case TOKEN_INLINE_FORMAT_END:
// 				return MDSyntaxTree{ns}, invalidTokenError(typ)
// 			case TOKEN_INLINE_LINK_DESC_START:
// 			case TOKEN_INLINE_LINK_DESC_END:
// 			case TOKEN_INLINE_LINK_URL_START:
// 			case TOKEN_INLINE_LINK_URL_END:
// 			case TOKEN_UNORDERED_LIST_INDIC:
// 			case TOKEN_ORDERED_LIST_INDIC:
// 			case TOKEN_HEADER_INDIC:
// 			case TOKEN_EXPLICIT_CODEBLOCK_INDIC:
// 			case TOKEN_SPECIAL_CHAR_ESCAPE:
// 				// TODO: Deal with all this shit
// 				continue
// 			default:
// 				return MDSyntaxTree{ns}, unknownTokenTypeError(typ)
// 			}
// 		} else if state.curProcessingType == SYNTAX_PARAGRAPH {
// 			typ := t.GetType()
// 			switch typ {
// 			case TOKEN_NL:
// 				state.prevWasNL = true
// 				continue
// 			case TOKEN_TEXT:
// 				t, ok := t.(MDTextToken)
// 				if !ok {
// 					return MDSyntaxTree{ns}, tokenTypeMismatchError(typ)
// 				}
// 				curNode := state.curNode
// 				curPNode := (*curNode).(MDParagraph)
// 				newContent = curPNode.Content
// 				var p MDSyntaxNode = MDParagraph{[]MDParagraphFormatNode{MDTextFormatNode(t)}}
// 				state.curNode = &p
// 			case TOKEN_LEADING_SPACE:
// 				// TODO: Deal with leading space
// 				continue
// 			case TOKEN_INLINE_FORMAT_START:
// 				// TODO: Deal with inline formatting
// 				continue
// 			case TOKEN_INLINE_FORMAT_MID:
// 			case TOKEN_INLINE_FORMAT_END:
// 				return MDSyntaxTree{ns}, invalidTokenError(typ)
// 			case TOKEN_INLINE_LINK_DESC_START:
// 			case TOKEN_INLINE_LINK_DESC_END:
// 			case TOKEN_INLINE_LINK_URL_START:
// 			case TOKEN_INLINE_LINK_URL_END:
// 			case TOKEN_UNORDERED_LIST_INDIC:
// 			case TOKEN_ORDERED_LIST_INDIC:
// 			case TOKEN_HEADER_INDIC:
// 			case TOKEN_EXPLICIT_CODEBLOCK_INDIC:
// 			case TOKEN_SPECIAL_CHAR_ESCAPE:
// 				// TODO: Deal with all this shit
// 				continue
// 			default:
// 				return MDSyntaxTree{ns}, unknownTokenTypeError(typ)
// 		}
// 	}

// 	return MDSyntaxTree{ns}, nil
// }
