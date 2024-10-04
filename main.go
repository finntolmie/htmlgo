package main

import (
	"fmt"
	"strings"

	"github.com/finntolmie/htmlgo/lexer"
)

func main() {
	input := strings.NewReader("<div id=\"main\">Hello</div>")
	lexer := lexer.NewLexer(input)
	tokens := lexer.Lex()
	for token := range tokens {
		fmt.Println(token)
	}
}
