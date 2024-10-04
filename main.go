package main

import (
	"fmt"
	"strings"

	"github.com/finntolmie/htmlgo/lexer"
)

func main() {
	input := strings.NewReader("<div id=\"main\">Hello</div>")
	lexer := lexer.NewLexer(input)
	lexer.Run()

	for _, token := range lexer.Tokens {
		fmt.Printf("%+v\n", token)
	}
}
