package argparser

import (
	"fmt"
)

const UnlimitedNArgs int = -1

type commonArgData struct {
	value ArgValue
	name  string
	usage string
	nArgs int // later convert to string for patterns like '*', '+'
}

func (c *commonArgData) setNArgs(n int) error {
	if c.value.IsBoolValue() {
		return fmt.Errorf("cannot change nargs for bool value, nargs for them is always 0")
	}
	if n == 0 {
		return fmt.Errorf("nargs cannot be zero for non Bool values")
	}
	c.nArgs = n
	return nil
}

type PosArg struct {
	common commonArgData
}

func NewPosArg(name string, value ArgValue, usage string) *PosArg {
	return &PosArg{
		common: commonArgData{
			nArgs: 1,
			name:  name,
			value: value,
			usage: usage,
		},
	}
}

func (pos *PosArg) SetNArgs(n int) error {
	return pos.common.setNArgs(n)
}

type OptArg struct {
	common commonArgData
	// isSwitch bool
	// mutex   map[string]bool
	// visited bool
	//repeat bool
}

func NewOptArg(name string, val ArgValue, usage string) *OptArg {
	nargs := 1
	if val.IsBoolValue() {
		nargs = 0
	}

	return &OptArg{
		common: commonArgData{
			nArgs: nargs,
			name:  name,
			usage: usage,
			value: val,
		},
	}
}

// for switch options this change is simply ignored
func (opt *OptArg) SetNArgs(n int) error {
	return opt.common.setNArgs(n)
}
