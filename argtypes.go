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
	if n == 0 {
		return fmt.Errorf("nargs cannot be zero")
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
	common   commonArgData
	isSwitch bool
	// mutex   map[string]bool
	// visited bool
	//repeat bool
}

func NewOptArg(name string, val ArgValue, sw bool, usage string) *OptArg {
	return &OptArg{
		common: commonArgData{
			nArgs: 1,
			name:  name,
			usage: usage,
			value: val,
		},
		isSwitch: sw,
	}
}

// for switch options this change is simply ignored
func (opt *OptArg) SetNArgs(n int) error {
	return opt.common.setNArgs(n)
}
