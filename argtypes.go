package argparser

import (
	"fmt"
)

const UnlimitedNArgs int = -1

type PosArg struct {
	Value ArgValue
	Help  string
	nArgs int // later convert to string for patterns like '*', '+'
}

func NewPosArg(value ArgValue, help string) *PosArg {
	return &PosArg{
		nArgs: 1,
		Value: value,
		Help:  help,
	}
}

func (pos *PosArg) SetNArgs(n int) error {
	if n == 0 {
		return fmt.Errorf("nargs cannot be zero for positional argument")
	}
	pos.nArgs = n
	return nil
}

type OptArg struct {
	Value ArgValue
	Help  string
	nArgs int // later convert to string for patterns like '*', '+'
	// isSwitch bool
	// mutex   map[string]bool
	// visited bool
	//repeat bool
}

func NewOptArg(value ArgValue, help string) *OptArg {
	nargs := 1
	if value.IsBoolValue() {
		nargs = 0
	}
	return &OptArg{
		nArgs: nargs,
		Value: value,
		Help:  help,
	}

}

func (opt *OptArg) SetNArgs(n int) error {
	if opt.Value.IsBoolValue() {
		return fmt.Errorf("cannot change nargs for optional bool argument, it is always 0")
	}
	if n == 0 {
		return fmt.Errorf("nargs cannot be zero for non Bool optional values")
	}
	opt.nArgs = n
	return nil
}
