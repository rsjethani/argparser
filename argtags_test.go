package argparser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplitKV(t *testing.T) {
	data := make(map[string][]string)
	data[""] = []string{}
	data[fmt.Sprintf("%[1]c%[1]c%[1]c", tagSep)] = []string{}
	data[fmt.Sprintf("a%[1]cb%[1]cc%[1]cd", tagSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca%[1]cb%[1]cc%[1]cd%[1]c", tagSep)] = []string{"a", "b", "c", "d"}
	data[fmt.Sprintf("%[1]ca\\%[1]cb%[1]cc%[1]cd\\%[1]c", tagSep)] = []string{fmt.Sprintf("a%cb", tagSep), "c", fmt.Sprintf("d%c", tagSep)}
	for input, expected := range data {
		if got := splitKV(input, tagSep); !reflect.DeepEqual(expected, got) {
			t.Errorf("testing: splitKV(%q, %c); expected: %q; Got: %q", input, tagSep, expected, got)
		}
	}
}

func TestParseTagsInvalidKeyValues(t *testing.T) {
	invalidKVs := []string{
		"hello",
		"help=",
		"hello=hi",
		"name=arg_name",
		"type=OPT",
		"nargs=1x",
	}

	for _, kv := range invalidKVs {
		if _, err := parseTags(kv); err == nil {
			t.Errorf("testing: parseTags(%#v); expected: error; got: no error", kv)
		}
	}
}

func TestParseTagsValidKeyValues(t *testing.T) {
	data := []struct {
		validKVs string
		expected map[string]string
	}{
		{
			"name=arg-name,type=pos,help=a help message,nargs=10",
			map[string]string{
				"nargs": "10",
				"type":  "pos",
				"name":  "arg-name",
				"help":  "a help message",
			},
		},
		{
			"name=ArgName10,type=opt,help=a,nargs=-10",
			map[string]string{
				"nargs": "-10",
				"type":  "opt",
				"name":  "ArgName10",
				"help":  "a",
			},
		},
		{
			"type=switch,nargs=-10",
			map[string]string{
				"nargs": "-10",
				"type":  "switch",
			},
		},
	}

	for _, input := range data {
		got, err := parseTags(input.validKVs)
		if err != nil {
			t.Errorf("testing: parseTags(%#v); expected: no error; got: %s", input.validKVs, err)
		}

		if !reflect.DeepEqual(input.expected, got) {
			t.Errorf("testing: parseTags(%#v); expected: %+v; got: %+v", input.validKVs, input.expected, got)
		}
	}
}

// func TestNewArgFromTags(t *testing.T) {
// 	testTags := make(map[string]string)

// 	testTags["type"] = "pos"
// 	testTags["nargs"] = "0"
// 	if arg, err := newArgFromTags(nil, testTags); arg != nil || err == nil {
// 		t.Errorf("testing: newArgFromTags(%#v); expected: nil Argument, non-nil error, got: %v, %v ", testTags, arg, err)
// 	}

// 	testTags["type"] = "switch"
// 	testTags["nargs"] = "10"
// 	if arg, err := newArgFromTags(nil, testTags); arg != nil || err == nil {
// 		t.Errorf("testing: newArgFromTags(%#v); expected: nil Argument, non-nil error, got: %v, %v ", testTags, arg, err)
// 	}

// }

// // Test default argument name should be lower case of field name
// args2 := struct {
// 	Field4 int `argparser:"type=opt"`
// }{}
// argset, err = NewArgSet(&args2)
// if err != nil {
// 	t.Errorf("testing: NewArgSet(%#v); expected: (non-nil *ArgSet, nil); got: (%v, %#v)", args2, argset, err)
// }
// if argset.optArgs["--field4"] == nil {
// 	t.Errorf(`testing: testing: NewArgSet(%#v); expected: ArgSet.optArgs["--field4"]!=nil; got: nil`, &args2)
// }
