package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexEmpty(t *testing.T) {
	input := strings.NewReader("")
	output := []Token{
		{Type: TokenEOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestLexDiv(t *testing.T) {
	input := strings.NewReader("<div></div>")
	output := []Token{
		{Type: TokenStartTag, Value: "div"},
		{Type: TokenEndTag, Value: "div"},
		{Type: TokenEOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestNestedElements(t *testing.T) {
	input := strings.NewReader("<div><a><div><a class=\"hi\"></a><b></b></div></a></div>")
	output := []Token{
		{Type: TokenStartTag, Value: "div"},
		{Type: TokenStartTag, Value: "a"},
		{Type: TokenStartTag, Value: "div"},
		{Type: TokenStartTag, Value: "a"},
		{Type: TokenAttributeName, Value: "class"},
		{Type: TokenAttributeValue, Value: "hi"},
		{Type: TokenEndTag, Value: "a"},
		{Type: TokenStartTag, Value: "b"},
		{Type: TokenEndTag, Value: "b"},
		{Type: TokenEndTag, Value: "div"},
		{Type: TokenEndTag, Value: "a"},
		{Type: TokenEndTag, Value: "div"},
		{Type: TokenEOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestLexDivId(t *testing.T) {
	input := strings.NewReader("<div id=\"main\"></div>")
	output := []Token{
		{Type: TokenStartTag, Value: "div"},
		{Type: TokenAttributeName, Value: "id"},
		{Type: TokenAttributeValue, Value: "main"},
		{Type: TokenEndTag, Value: "div"},
		{Type: TokenEOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestLexDivMultiline(t *testing.T) {
	input := strings.NewReader(`<div>
		<a class="big"></a>
	</div>
	`)
	output := []Token{
		{Type: TokenStartTag, Value: "div"},
		{Type: TokenStartTag, Value: "a"},
		{Type: TokenAttributeName, Value: "class"},
		{Type: TokenAttributeValue, Value: "big"},
		{Type: TokenEndTag, Value: "a"},
		{Type: TokenEndTag, Value: "div"},
		{Type: TokenEOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}
