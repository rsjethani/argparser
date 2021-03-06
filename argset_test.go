package argparser

import (
	"testing"
)

func TestNewArgSet(t *testing.T) {
	argset := NewArgSet()
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

func TestArgSetAdd(t *testing.T) {
	argset := NewArgSet()

	argset.Add("dummy", nil)
	if len(argset.posArgs) != 0 || argset.optArgs["--dummy"] != nil {
		t.Errorf(`testing: argset.AddArgument("dummy", nil); expected: no positional/optional argument named 'dummy should get added; got: 'dummy' got added`)
	}

	argset.Add("pos1", NewPosArg(nil, ""))
	if len(argset.posArgs) == 0 || argset.posArgs[0].name != "pos1" {
		t.Errorf(`testing: argset.AddArgument("pos1", NewPosArg(nil, "")); expected: argset.posArgs[0].name == "pos1"; got: len(argset.posArgs) == 0`)
	}

	argset.Add("opt1", NewOptArg(nil, ""))
	if argset.optArgs["--opt1"] == nil {
		t.Errorf(`testing: argset.AddArgument("opt1", NewOptArg(nil, "")); expected: argset.optArgs["opt1"] != nil; got: argset.optArgs["opt1"] == nil`)
	}
}

func TestNewArgSetFromInvalidInputs(t *testing.T) {
	data := []interface{}{
		// Test nil as input
		nil,
		// Test non-pointer as input
		*new(bool),
		// Test pointer to a non-struct as input
		new(bool),
		// Test unexported tagged field as input
		&struct {
			field1 int `argparser:""`
		}{},
		// Test unsupported field type as input
		&struct {
			Field1 int8 `argparser:""`
		}{},
		// Test invalid tag/value as input
		&struct {
			Field1 int `argparser:"type=xxx"`
		}{},
	}
	for _, input := range data {
		if argset, err := NewArgSetFrom(input); argset != nil || err == nil {
			t.Errorf("testing: NewArgSet(%#v); expected: (nil, error); got: (%v, %#v)", input, argset, err)
		}
	}
}

func TestNewArgSetFromValidInputs(t *testing.T) {
	// Test skipping of untagged fields
	args1 := struct {
		Field1 int // no 'argparser' tag hence should be skipped
	}{}
	argset, err := NewArgSetFrom(&args1)
	if err != nil {
		t.Errorf("testing: NewArgSet(%#v); expected: non-nil *ArgSet and nil error; got: %v", args1, err)
	}
	if len(argset.posArgs) != 0 || len(argset.optArgs) != 1 {
		t.Errorf("testing: NewArgSet(%#v); expected: no arguments except --help in argset; got: %#v", &args1, argset)
	}

	// Test parsing of tagged fields and no error with valid key/values
	args2 := struct {
		Field1 int `argparser:""`         // optional argument
		Field2 int `argparser:"type=pos"` // positional argument
	}{}
	argset, err = NewArgSetFrom(&args2)
	if err != nil {
		t.Errorf("testing: NewArgSet(%#v); expected: non-nil *ArgSet and nil error; got: %v", args2, err)
	}
	if len(argset.posArgs) == 0 || len(argset.optArgs) == 0 {
		t.Errorf("testing: NewArgSet(%#v); expected: 1 optional and 1 positional arguments in argset; got: %#v", &args2, argset)
	}
}

func TestUsage(t *testing.T) {
	args1 := struct {
		Pos1                   int     `argparser:"type=pos,help=pos1 help"`
		Pos2                   bool    `argparser:"type=pos,help=pos2 help"`
		Posssssssssssssssssss3 string  `argparser:"type=pos,help=pos3 help"`
		Pos4                   float64 `argparser:"type=pos,help=pos4 help"`
		Pos5                   []int   `argparser:"type=pos,help=pos5 help,nargs=2"`
		Opt1                   int     `argparser:"help=opt1 help"`
		Opt2                   bool    `argparser:"help=opt2 help"`
		Optfffffffffffffff3    string  `argparser:"help=opt3 help"`
		Opt4                   float64 `argparser:"help=opt4 help"`
		Opt5                   []int   `argparser:"help=opt5 help,nargs=3"`
		Sw1                    bool    `argparser:"type=switch,help=sw1 help"`
	}{}

	argSet, err := NewArgSetFrom(&args1)
	if err != nil {
		t.Error(err)
	}
	argSet.Description = "affaadadaddadsdasdasdddddddddddddddddddddddddddasda adsdasddd\nasdddadadadada"

	argSet.usage()
}
