package error

import (
	"fmt"
	"os"

	"github.com/TwiN/go-color"
)

type ErrorType string

const (
	SYNTAX_ERROR = "Syntax error"
	PARSER_ERROR = "Parser error"
)

func PrintError(t ErrorType, msg string) {
	fmt.Println(color.InRed(fmt.Sprintf(`%s: %s`, t, msg)))
	os.Exit(1)
}
