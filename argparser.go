package argparser

import (
	"fmt"
	"io"
	"os"
)

type Parser struct {
	mainArgSet *ArgSet
	usageOut   io.Writer
}

func (parser *Parser) SetOutput(w io.Writer) {
	if w == nil {
		parser.usageOut = os.Stderr
		return
	}
	parser.usageOut = w
}

func NewParser(argSet *ArgSet) *Parser {
	parser := &Parser{mainArgSet: argSet, usageOut: os.Stderr}
	return parser
}

func (parser *Parser) Usage() {
	fmt.Fprintf(parser.usageOut, "Usage of %s:\n", os.Args[0])
	parser.mainArgSet.usage(parser.usageOut)
	fmt.Fprintln(parser.usageOut, "")
}

func (parser *Parser) ParseFrom(args []string) error {
	return parser.mainArgSet.parseFrom(args)
}

func (parser *Parser) Parse() error {
	return parser.ParseFrom(os.Args[1:])
}
