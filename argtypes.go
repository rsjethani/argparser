package argparser

import (
	"fmt"
)

type Argument struct {
	value      Value
	help       string
	positional bool
	nArgs      int // TODO: convert to string for patterns like '*', '+' etc.
}

func NewPosArg(value Value, help string) *Argument {
	return &Argument{
		nArgs:      1,
		value:      value,
		help:       help,
		positional: true,
	}
}

func NewOptArg(value Value, help string) *Argument {
	return &Argument{
		nArgs:      1,
		value:      value,
		help:       help,
		positional: false,
	}
}

func NewSwitchArg(value Value, help string) *Argument {
	return &Argument{
		nArgs:      0,
		value:      value,
		help:       help,
		positional: false,
	}
}

func (arg *Argument) isSwitch() bool {
	return !arg.positional && arg.nArgs == 0
}

func (arg *Argument) SetNArgs(n int) error {
	if n == 0 && arg.positional {
		return fmt.Errorf("nargs cannot be 0 for positional argument")
	}
	arg.nArgs = n
	return nil
}
