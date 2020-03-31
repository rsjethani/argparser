package argparser

import (
	"testing"
)

func TestNewPosArg(t *testing.T) {
	arg := NewPosArg(nil, "help string")
	if !arg.isPositional() || arg.Value != nil || arg.NArgs() != 1 || arg.Help != "help string" {
		t.Errorf("Expected: positional=true, value=nil, nargs=1, help=help string; Got: %+v", arg)
	}
}

func TestNewOptArg(t *testing.T) {
	arg := NewOptArg(nil, "help string")
	if arg.isPositional() || arg.Value != nil || arg.NArgs() != 1 || arg.Help != "help string" {
		t.Errorf("Expected: positional=false, value=nil, nargs=1, help=help string; Got: %+v", arg)
	}
}

func TestNewSwitchArg(t *testing.T) {
	arg := NewSwitchArg(nil, "help string")
	if arg.isPositional() || arg.Value != nil || arg.NArgs() != 0 || arg.Help != "help string" {
		t.Errorf("Expected: positional=false, value=nil, nargs=0, help=help string; Got: %+v", arg)
	}
}

func TestSetNArgs(t *testing.T) {
	posArg := NewPosArg(nil, "")
	if err := posArg.SetNArgs(10); err != nil || posArg.NArgs() != 10 {
		t.Errorf("Expected: for positional argument %[1]T.SetNArgs(10) suceeds with nil error setting %[1]T.NArgs()==10; Got: error", posArg)
	}
	if err := posArg.SetNArgs(0); err == nil {
		t.Errorf("Expected: for positional argument %T.SetNArgs(0) results in error; Got: nil error", posArg)
	}

	optArg := NewOptArg(nil, "")
	if err := optArg.SetNArgs(10); err != nil || optArg.NArgs() != 10 {
		t.Errorf("Expected: for optional argument %[1]T.SetNArgs(10) suceeds with nil error setting %[1]T.NArgs()==10; Got: error", optArg)
	}
	if err := optArg.SetNArgs(0); err != nil || optArg.NArgs() != 0 {
		t.Errorf("Expected: for optional argument %[1]T.SetNArgs(0) suceeds with no error setting %[1]T.NArgs()==0; Got: error", optArg)
	}
}
