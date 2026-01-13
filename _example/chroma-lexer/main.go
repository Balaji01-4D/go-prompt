package main

import (
    "fmt"
    "github.com/balaji01-4d/go-prompt"
)

func main() {
    // Create the Chroma lexer for Go language with monokai style
    lexer, err := prompt.NewChromaLexer("go", "monokai")
    if err != nil {
        panic(err)
    }

    p := prompt.New(
        func(input string) {
			fmt.Println("You entered:", input)
		},
		prompt.WithLexer(lexer),
	)

	p.Run()
}
