package argparser

import (
	"fmt"
	"reflect"
)

type optArg struct {
	value interface{}
	usage string
	short string
}

type posArg struct {
	value interface{}
	usage string
}

type ArgSet struct {
	Description string
	posArgs     map[string]*posArg
	optArgs     map[string]*optArg
}

func NewArgSet(args interface{}) *ArgSet {
	argset := &ArgSet{
		posArgs: make(map[string]*posArg),
		optArgs: make(map[string]*optArg),
	}
	t := reflect.TypeOf(args).Elem()
	// v := reflect.ValueOf(args).Elem()
	// Get the type and kind of our user variable
	fmt.Println("Type:", t.Name())
	fmt.Println("Kind:", t.Kind())

	// Iterate over all available fields and read the tag value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the field tag value
		if _, present := field.Tag.Lookup("opt"); present {
			//TODO: check and raise error if value other than "yes" is given
			if nm, present := field.Tag.Lookup("name"); present {
				argset.optArgs[nm] = &optArg{
					usage: field.Tag.Get("usage"),
					short: field.Tag.Get("short"),
				}
			}
		} else {
			if nm, present := field.Tag.Lookup("name"); present {
				argset.posArgs[nm] = &posArg{usage: field.Tag.Get("usage")}
			}
		}

		fmt.Printf("%v %v %v\n", field.Name, field.Type.Name(), field.Type.Kind())
	}
	// fmt.Printf("\n%+v\n", argset.optArgs["emp-id"])
	return argset
}

// func (set *ArgSet) AddOptional(val Value, name string, usage string) {
// 	set.optArgs[name] = &optArg{value: val, usage: usage}
// }

// func (set *ArgSet) AddPositional(val Value, name string, usage string) {
// 	set.posArgs[name] = &posArg{value: val, usage: usage}
// }
