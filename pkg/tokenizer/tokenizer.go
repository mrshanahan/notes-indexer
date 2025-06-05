package tokenizer

var (
	DefaultSeparators = "\t\n\r ,.:?\"!;()"
)

const (
	TOKEN_TYPE_GENERIC = 0
	TOKEN_TYPE_XML     = 1
)

type Token struct {
	Value string
	Type  int
}

type Tokenizer interface {
	Tokenize(text string) ([]Token, error)
}
