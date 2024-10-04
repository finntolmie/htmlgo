package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexEmpty(t *testing.T) {
	input := strings.NewReader("")
	output := []Token{
		{Type: EOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestLexDiv(t *testing.T) {
	input := strings.NewReader("<div></div>")
	output := []Token{
		{Type: StartTag, Value: "div"},
		{Type: EndTag, Value: "div"},
		{Type: EOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestLexDivText(t *testing.T) {
	input := strings.NewReader("<div>text</div>")
	output := []Token{
		{Type: StartTag, Value: "div"},
		{Type: Text, Value: "text"},
		{Type: EndTag, Value: "div"},
		{Type: EOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestNestedElements(t *testing.T) {
	input := strings.NewReader("<div><a><div><a class=\"hi\"></a><b></b></div></a></div>")
	output := []Token{
		{Type: StartTag, Value: "div"},
		{Type: StartTag, Value: "a"},
		{Type: StartTag, Value: "div"},
		{Type: StartTag, Value: "a"},
		{Type: AttributeName, Value: "class"},
		{Type: AttributeValue, Value: "hi"},
		{Type: EndTag, Value: "a"},
		{Type: StartTag, Value: "b"},
		{Type: EndTag, Value: "b"},
		{Type: EndTag, Value: "div"},
		{Type: EndTag, Value: "a"},
		{Type: EndTag, Value: "div"},
		{Type: EOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestLexDivId(t *testing.T) {
	input := strings.NewReader("<div id=\"main\"></div>")
	output := []Token{
		{Type: StartTag, Value: "div"},
		{Type: AttributeName, Value: "id"},
		{Type: AttributeValue, Value: "main"},
		{Type: EndTag, Value: "div"},
		{Type: EOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}

func TestLexDivMultipleAttributes(t *testing.T) {
	input := strings.NewReader("<div id=\"main\" class=\"huge\"></div>")
	output := []Token{
		{Type: StartTag, Value: "div"},
		{Type: AttributeName, Value: "id"},
		{Type: AttributeValue, Value: "main"},
		{Type: AttributeName, Value: "class"},
		{Type: AttributeValue, Value: "huge"},
		{Type: EndTag, Value: "div"},
		{Type: EOF},
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
		{Type: StartTag, Value: "div"},
		{Type: StartTag, Value: "a"},
		{Type: AttributeName, Value: "class"},
		{Type: AttributeValue, Value: "big"},
		{Type: EndTag, Value: "a"},
		{Type: EndTag, Value: "div"},
		{Type: EOF},
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
		{Type: StartTag, Value: "html"},
		{Type: StartTag, Value: "body"},

		{Type: StartTag, Value: "h1"},
		{Type: Text, Value: "Title"},
		{Type: EndTag, Value: "h1"},

		{Type: StartTag, Value: "div"},

		{Type: AttributeName, Value: "id"},
		{Type: AttributeValue, Value: "main"},

		{Type: AttributeName, Value: "class"},
		{Type: AttributeValue, Value: "test"},

		{Type: StartTag, Value: "p"},
		{Type: Text, Value: "Hello "},
		{Type: StartTag, Value: "em"},
		{Type: Text, Value: "world"},
		{Type: EndTag, Value: "em"},
		{Type: Text, Value: "!"},
		{Type: EndTag, Value: "p"},

		{Type: EndTag, Value: "div"},

		{Type: EndTag, Value: "body"},
		{Type: EndTag, Value: "html"},

		{Type: EOF},
	}
	l := NewLexer(input)
	l.Run()
	assert.Equal(t, output, l.Tokens)
}
