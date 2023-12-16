package parser

import (
	"errors"

	"github.com/usememos/memos/plugin/gomark/ast"
	"github.com/usememos/memos/plugin/gomark/parser/tokenizer"
)

type OrderedListParser struct{}

func NewOrderedListParser() *OrderedListParser {
	return &OrderedListParser{}
}

func (*OrderedListParser) Match(tokens []*tokenizer.Token) (int, bool) {
	if len(tokens) < 4 {
		return 0, false
	}
	if tokens[0].Type != tokenizer.Number || tokens[1].Type != tokenizer.Dot || tokens[2].Type != tokenizer.Space {
		return 0, false
	}

	contentTokens := []*tokenizer.Token{}
	for _, token := range tokens[3:] {
		contentTokens = append(contentTokens, token)
		if token.Type == tokenizer.Newline {
			break
		}
	}

	if len(contentTokens) == 0 {
		return 0, false
	}

	return len(contentTokens) + 3, true
}

func (p *OrderedListParser) Parse(tokens []*tokenizer.Token) (ast.Node, error) {
	size, ok := p.Match(tokens)
	if size == 0 || !ok {
		return nil, errors.New("not matched")
	}

	contentTokens := tokens[3:size]
	if contentTokens[len(contentTokens)-1].Type == tokenizer.Newline {
		contentTokens = contentTokens[:len(contentTokens)-1]
	}
	children, err := ParseInline(contentTokens)
	if err != nil {
		return nil, err
	}
	return &ast.OrderedList{
		Number:   tokens[0].Value,
		Children: children,
	}, nil
}
