package main

import (
	"fmt"

	"github.com/rsjethani/argparser"
)

func main() {
	var pos1 int
	var sw1 bool
	var opt1 string

	set := argparser.NewArgSet()
	set.AddArgument("pos1", argparser.NewPosArg(argparser.NewInt(&pos1), "pos1 help"))
	set.AddArgument("opt1", argparser.NewOptArg(argparser.NewString(&opt1), "pos1 help"))
	set.AddArgument("sw1", argparser.NewSwitchArg(argparser.NewBool(&sw1), "sw1 help"))

	parser := argparser.NewArgParser(set)
	parser.Usage()
	fmt.Println("before parse: ", pos1, opt1, sw1)
	fmt.Println(parser.Parse())
	fmt.Println("after parse: ", pos1, opt1, sw1)
}
