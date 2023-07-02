package main

import (
	"encoding/json"
	"fmt"

	"github.com/jellycat-io/gero/lexer"
	"github.com/jellycat-io/gero/parser"
)

func main() {
	input := `42 "Hello"`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Program()
	json, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(json))
}
