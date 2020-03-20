package argparser

import (
	"fmt"
)

type Argument struct {
	Value      Value
	Help       string
	Positional bool
	nArgs      int // TODO: convert to string for patterns like '*', '+' etc.
}

func NewPosArg(value Value, help string) *Argument {
	return &Argument{
		nArgs:      1,
		Value:      value,
		Help:       help,
		Positional: true,
	}
}

func NewOptArg(value Value, help string) *Argument {
	nargs := 1
	if value.IsSwitch() {
		nargs = 0
	}
	return &Argument{
		nArgs: nargs,
		Value: value,
		Help:  help,
	}
}

func (arg *Argument) SetNArgs(n int) error {
	if arg.Value.IsSwitch() && !arg.Positional {
		return fmt.Errorf("cannot change nargs for optional bool argument, it is always 0")
	}
	if n == 0 {
		return fmt.Errorf("nargs cannot be zero for non boolean argument")
	}
	arg.nArgs = n
	return nil
}
