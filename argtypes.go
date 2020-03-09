package argparser

import (
	"fmt"
)

const UnlimitedNArgs int = -1

type PosArg struct {
	Value ArgValue
	Name  string
	Help  string
	nArgs int // later convert to string for patterns like '*', '+'
}

func NewPosArg(value ArgValue, name string, help string) *PosArg {
	return &PosArg{
		nArgs: 1,
		Name:  name,
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

func (pos *PosArg) Usage() string {
	return fmt.Sprintf("%-15s%s", pos.Name, pos.Help)
}

type OptArg struct {
	Value ArgValue
	Name  string
	Help  string
	nArgs int // later convert to string for patterns like '*', '+'
	// isSwitch bool
	// mutex   map[string]bool
	// visited bool
	//repeat bool
}

func NewOptArg(value ArgValue, name string, help string) *OptArg {
	nargs := 1
	if value.IsBoolValue() {
		nargs = 0
	}
	return &OptArg{
		nArgs: nargs,
		Name:  name,
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

func (opt *OptArg) Usage() string {
	return fmt.Sprintf("%-15s%s", opt.Name, opt.Help)
}
