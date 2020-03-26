package argparser

import (
	"fmt"
	"os"
)

type ArgParser struct {
	mainArgSet *ArgSet
}

func NewArgParser(argSet *ArgSet) *ArgParser {
	parser := &ArgParser{mainArgSet: argSet}
	return parser
}

func (parser *ArgParser) Usage() {
	fmt.Printf("Usage of %s:\n%s", os.Args[0], parser.mainArgSet.Usage())
}

func (parser *ArgParser) ParseFrom(args []string) error {
	return parser.mainArgSet.ParseFrom(args)
}

func (parser *ArgParser) Parse() error {
	return parser.ParseFrom(os.Args)
}
