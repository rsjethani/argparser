package argparser

import (
	"testing"
)

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

func TestAddArgument(t *testing.T) {
	argset := DefaultArgSet()

	argset.AddArgument("dummy", nil)
	if len(argset.posArgs) != 0 || argset.optArgs["--dummy"] != nil {
		t.Errorf(`testing: argset.AddArgument("dummy", nil); expected: no positional/optional argument named 'dummy should get added; got: 'dummy' got added`)
	}

	argset.AddArgument("pos1", NewPosArg(nil, ""))
	if len(argset.posArgs) == 0 || argset.posArgs[0].name != "pos1" {
		t.Errorf(`testing: argset.AddArgument("pos1", NewPosArg(nil, "")); expected: argset.posArgs[0].name == "pos1"; got: len(argset.posArgs) == 0`)
	}

	argset.AddArgument("opt1", NewOptArg(nil, ""))
	if argset.optArgs["--opt1"] == nil {
		t.Errorf(`testing: argset.AddArgument("opt1", NewOptArg(nil, "")); expected: argset.optArgs["opt1"] != nil; got: argset.optArgs["opt1"] == nil`)
	}
}

func TestNewArgSetInvalidInputs(t *testing.T) {
	data := []interface{}{
		// Test invalid case of nil input
		nil,
		// Test invalid case of non-pointer input
		*new(bool),
		// Test invalid case of pointer to a non-struct input
		new(bool),
		// Test invalid case of unexported tagged field as input
		&struct {
			field1 int `argparser:""`
		}{},
		// Test invalid case of unsupported field type as input
		&struct {
			Field1 int8 `argparser:""`
		}{},
		// Test invalid case of invalid tag/value as input
		&struct {
			Field1 int `argparser:"type=xxx"`
		}{},
		// Test invalid case of invalid tag/value as input
		&struct {
			Field1 int `argparser:"type=pos,nargs=0"`
		}{},
	}
	for _, input := range data {
		if argset, err := NewArgSet(input); argset != nil || err == nil {
			t.Errorf("testing: NewArgSet(%#v); expected: (nil, error); got: (%v, %#v)", input, argset, err)
		}
	}
}

func TestNewArgSetValidInputs(t *testing.T) {
	args1 := struct {
		Field1 int // no 'argparser' tag hence should be skipped
	}{}
	argset, err := NewArgSet(&args1)
	if err != nil {
		t.Errorf("testing: NewArgSet(%#v); expected: (non-nil *ArgSet, nil); got: (%v, %#v)", args1, argset, err)
	}
	if len(argset.posArgs) != 0 || argset.optArgs["--field1"] != nil {
		t.Errorf("testing: NewArgSet(%#v); expected: no optional/positional argument should be added; got: %#v", &args1, argset)
	}

	// Test default argument name should be lower case of field name
	args2 := struct {
		Field4 int `argparser:"type=opt"`
	}{}
	argset, err = NewArgSet(&args2)
	if err != nil {
		t.Errorf("testing: NewArgSet(%#v); expected: (non-nil *ArgSet, nil); got: (%v, %#v)", args2, argset, err)
	}
	if argset.optArgs["--field4"] == nil {
		t.Errorf(`testing: testing: NewArgSet(%#v); expected: ArgSet.optArgs["--field4"]!=nil; got: nil`, &args2)
	}
}
