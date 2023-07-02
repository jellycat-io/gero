package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/user"

	"github.com/TwiN/go-color"
	"github.com/jellycat-io/gero/lexer"
	"github.com/jellycat-io/gero/parser"
	"github.com/jellycat-io/gero/util"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(in)

	fmt.Print(color.InBold(color.InBlue(fmt.Sprintf("Gero REPL 0.1.0 - Welcome %s\n", user.Username))))

	for {
		fmt.Fprint(out, color.InBold(PROMPT))
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.Program()

		if len(p.Errors()) != 0 {
			util.PrintParserErrors(out, p.Errors())
		}

		json, err := json.MarshalIndent(program, "", "  ")
		if err != nil {
			io.WriteString(out, err.Error())
		}

		io.WriteString(out, string(json)+"\n")
	}
}
