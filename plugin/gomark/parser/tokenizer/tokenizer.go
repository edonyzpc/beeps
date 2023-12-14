package tokenizer

type TokenType = string

const (
	Underline          TokenType = "_"
	Asterisk           TokenType = "*"
	Hash               TokenType = "#"
	Backtick           TokenType = "`"
	LeftSquareBracket  TokenType = "["
	RightSquareBracket TokenType = "]"
	LeftParenthesis    TokenType = "("
	RightParenthesis   TokenType = ")"
	ExclamationMark    TokenType = "!"
	Tilde              TokenType = "~"
	Dash               TokenType = "-"
	GreaterThan        TokenType = ">"
	Newline            TokenType = "\n"
	Space              TokenType = " "
)

const (
	Text TokenType = ""
)

type Token struct {
	Type  TokenType
	Value string
}

func NewToken(tp, text string) *Token {
	return &Token{
		Type:  tp,
		Value: text,
	}
}

func Tokenize(text string) []*Token {
	tokens := []*Token{}
	for _, c := range text {
		switch c {
		case '_':
			tokens = append(tokens, NewToken(Underline, "_"))
		case '*':
			tokens = append(tokens, NewToken(Asterisk, "*"))
		case '#':
			tokens = append(tokens, NewToken(Hash, "#"))
		case '`':
			tokens = append(tokens, NewToken(Backtick, "`"))
		case '[':
			tokens = append(tokens, NewToken(LeftSquareBracket, "["))
		case ']':
			tokens = append(tokens, NewToken(RightSquareBracket, "]"))
		case '(':
			tokens = append(tokens, NewToken(LeftParenthesis, "("))
		case ')':
			tokens = append(tokens, NewToken(RightParenthesis, ")"))
		case '!':
			tokens = append(tokens, NewToken(ExclamationMark, "!"))
		case '~':
			tokens = append(tokens, NewToken(Tilde, "~"))
		case '-':
			tokens = append(tokens, NewToken(Dash, "-"))
		case '>':
			tokens = append(tokens, NewToken(GreaterThan, ">"))
		case '\n':
			tokens = append(tokens, NewToken(Newline, "\n"))
		case ' ':
			tokens = append(tokens, NewToken(Space, " "))
		default:
			var lastToken *Token
			if len(tokens) > 0 {
				lastToken = tokens[len(tokens)-1]
			}
			if lastToken == nil || lastToken.Type != Text {
				tokens = append(tokens, NewToken(Text, string(c)))
			} else {
				lastToken.Value += string(c)
			}
		}
	}
	return tokens
}

func (t *Token) String() string {
	return t.Value
}

func Stringify(tokens []*Token) string {
	text := ""
	for _, token := range tokens {
		text += token.String()
	}
	return text
}
