package argparser

import (
	"fmt"
	"io"
	"os"
)

type ArgParser struct {
	mainArgSet *ArgSet
	usageOut   io.Writer
}

func (parser *ArgParser) SetOutput(w io.Writer) {
	if w == nil {
		parser.usageOut = os.Stderr
		return
	}
	parser.usageOut = w
}

func NewArgParser(argSet *ArgSet) *ArgParser {
	parser := &ArgParser{mainArgSet: argSet, usageOut: os.Stderr}
	return parser
}

func (parser *ArgParser) Usage() {
	fmt.Fprintf(parser.usageOut, "Usage of %s:\n", os.Args[0])
	parser.mainArgSet.usage(parser.usageOut)
	fmt.Fprintln(parser.usageOut, "")
}

func (parser *ArgParser) ParseFrom(args []string) error {
	return parser.mainArgSet.parseFrom(args)
}

func (parser *ArgParser) Parse() error {
	return parser.ParseFrom(os.Args[1:])
}
