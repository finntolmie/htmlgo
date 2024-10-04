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

func TestLexDivText(t *testing.T) {
	input := strings.NewReader("<div>text</div>")
	output := []Token{
		{Type: TokenStartTag, Value: "div"},
		{Type: TokenText, Value: "text"},
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

func TestLexDivMultipleAttributes(t *testing.T) {
	input := strings.NewReader("<div id=\"main\" class=\"huge\"></div>")
	output := []Token{
		{Type: TokenStartTag, Value: "div"},
		{Type: TokenAttributeName, Value: "id"},
		{Type: TokenAttributeValue, Value: "main"},
		{Type: TokenAttributeName, Value: "class"},
		{Type: TokenAttributeValue, Value: "huge"},
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

func TestLexRun(t *testing.T) {
	input := strings.NewReader(`<html>
    <body>
        <h1>Title</h1>
        <div id="main" class="test">
            <p>Hello <em>world</em>!</p>
        </div>
    </body>
</html>`)
	output := []Token{
		{Type: TokenStartTag, Value: "html"},
		{Type: TokenStartTag, Value: "body"},

		{Type: TokenStartTag, Value: "h1"},
		{Type: TokenText, Value: "Title"},
		{Type: TokenEndTag, Value: "h1"},

		{Type: TokenStartTag, Value: "div"},

		{Type: TokenAttributeName, Value: "id"},
		{Type: TokenAttributeValue, Value: "main"},

		{Type: TokenAttributeName, Value: "class"},
		{Type: TokenAttributeValue, Value: "test"},

		{Type: TokenStartTag, Value: "p"},
		{Type: TokenText, Value: "Hello "},
		{Type: TokenStartTag, Value: "em"},
		{Type: TokenText, Value: "world"},
		{Type: TokenEndTag, Value: "em"},
		{Type: TokenText, Value: "!"},
		{Type: TokenEndTag, Value: "p"},

		{Type: TokenEndTag, Value: "div"},

		{Type: TokenEndTag, Value: "body"},
		{Type: TokenEndTag, Value: "html"},

		{Type: TokenEOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}
