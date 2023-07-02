package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/jellycat-io/gero/lexer"
	"github.com/jellycat-io/gero/parser"
	"github.com/jellycat-io/gero/util"
)

func main() {
	out := os.Stdout
	input := `
		"hello"
		3.14
		"I"
		"am"
		33
	`
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Program()
	if len(p.Errors()) != 0 {
		util.PrintParserErrors(out, p.Errors())
	}
	json, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		io.WriteString(out, err.Error())
	}
	io.WriteString(out, string(json))
}
