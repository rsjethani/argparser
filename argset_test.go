package argparser

import "testing"

func TestDefaultArgSet(t *testing.T) {
	argset := DefaultArgSet()
	if argset.OptArgPrefix != defaultOptArgPrefix {
		t.Errorf("testing: DefaultArgSet(); expected: argset.OptArgPrefix==%#v; got: argset.OptArgPrefix==%#v", defaultOptArgPrefix, argset.OptArgPrefix)
	}

	if argset.optArgs == nil {
		t.Errorf("testing: DefaultArgSet(); expected: argset.optArgs!=nil; got: argset.optArgs==nil")
	}
	if len(argset.posArgs) != 0 {
		t.Errorf("testing: DefaultArgSet(); expected: len(argset.posArgs)==0; got: len(argset.posArgs)==%d", len(argset.posArgs))
	}
}

func TestNewArgSetInvalidInputs(t *testing.T) {
	// Test invalid case of nil input
	if argset, err := NewArgSet(nil); argset != nil || err == nil {
		t.Errorf("testing: NewArgSet(nil); expected: (nil, error); got: (%v, %#v)", argset, err)
	}

	// Test invalid case of non-pointer input
	x := 100
	if argset, err := NewArgSet(x); argset != nil || err == nil {
		t.Errorf("testing: NewArgSet(%#v); expected: (nil, error); got: (%v, %#v)", x, argset, err)
	}

	// Test invalid case of pointer to a non-struct input
	if argset, err := NewArgSet(&x); argset != nil || err == nil {
		t.Errorf("testing: NewArgSet(%#v); expected: (nil, error); got: (%v, %#v)", &x, argset, err)
	}

	// Test invalid case of unexported tagged field as input
	arg := struct {
		field1 int `argparser:""`
	}{}
	if argset, err := NewArgSet(&arg); argset != nil || err == nil {
		t.Errorf("testing: NewArgSet(%+v); expected: (nil, error); got: (%v, %#v)", &arg, argset, err)
	}
}

func TestNewArgSetValidInputs(t *testing.T) {

}

func TestAddArgument(t *testing.T) {
	argset := DefaultArgSet()
	argset.AddArgument("dummy", nil)
	argset.AddArgument("pos1", NewPosArg(nil, ""))
	argset.AddArgument("opt1", NewOptArg(nil, ""))

	if len(argset.posArgs) == 0 || argset.posArgs[0].name != "pos1" {
		t.Errorf(`testing: argset.AddArgument("pos1", NewPosArg(nil, "")); expected: argset.posArgs[0].name == "pos1"; got: len(argset.posArgs) == 0`)
	}

	if argset.optArgs["--opt1"] == nil {
		t.Errorf(`testing: argset.AddArgument("opt1", NewOptArg(nil, "")); expected: argset.optArgs["opt1"] != nil; got: argset.optArgs["opt1"] == nil`)
	}
}
