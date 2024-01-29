package parser

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/usememos/memos/plugin/gomark/ast"
	"github.com/usememos/memos/plugin/gomark/parser/tokenizer"
	"github.com/usememos/memos/plugin/gomark/restore"
)

func TestHorizontalRuleParser(t *testing.T) {
	tests := []struct {
		text           string
		horizontalRule ast.Node
	}{
		{
			text: "---",
			horizontalRule: &ast.HorizontalRule{
				Symbol: "-",
			},
		},
		{
			text: "---\naaa",
			horizontalRule: &ast.HorizontalRule{
				Symbol: "-",
			},
		},
		{
			text:           "****",
			horizontalRule: nil,
		},
		{
			text: "***",
			horizontalRule: &ast.HorizontalRule{
				Symbol: "*",
			},
		},
		{
			text:           "-*-",
			horizontalRule: nil,
		},
	}

	for _, test := range tests {
		tokens := tokenizer.Tokenize(test.text)
		node, _ := NewHorizontalRuleParser().Match(tokens)
		require.Equal(t, restore.Restore([]ast.Node{test.horizontalRule}), restore.Restore([]ast.Node{node}))
	}
}
