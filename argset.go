package argparser

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	tagKey         string = "argparser"
	tagValueSep    string = ","
	mapKeyValueSep string = "="
)

func makeMapFromTags(tagValue string) (map[string]string, error) {
	tagMap := make(map[string]string)
	for _, value := range strings.Split(tagValue, tagValueSep) {
		parts := strings.Split(value, mapKeyValueSep)
		//TODO: verify key is a proper non-empty string without special symbols etc.
		tagMap[parts[0]] = parts[1]
	}
	return tagMap, nil
}

type ArgSet struct {
	Title       string
	Description string
	posArgs     map[string]*PosArg
	optArgs     map[string]*OptArg
}

func (argset *ArgSet) addArgument(fieldType reflect.StructField, fieldVal reflect.Value, argAttrs map[string]string) error {
	var argName string
	if val, ok := argAttrs["name"]; ok {
		argName = val
	} else {
		//TODO: convert field name with '-' if multiple words
		argName = fieldType.Name
	}

	var argUsage string
	if val, ok := argAttrs["usage"]; ok {
		argUsage = val
	}

	argVal, err := NewArgValue(fieldVal.Addr().Interface())
	if err != nil {
		return err
	}

	// check whether user wants positional or optional argument and process accordinly
	if _, wantsPos := argAttrs["positional"]; wantsPos {
		// TODO: verify value of 'positional is yes/true only'

		argset.AddPositional(NewPosArg(argName, argVal, argUsage))

	} else { // user wants optional argument
		isSwitch := false
		if val, ok := argAttrs["switch"]; ok {
			if val == "true" {
				isSwitch = true
			}
		}
		argset.AddOptional(NewOptArg(argName, argVal, isSwitch, argUsage))

	}
	return nil
}

func NewArgSet(src interface{}) (*ArgSet, error) {
	// get Type data of src, verify that it is of pointer type
	srcTyp := reflect.TypeOf(src)
	if srcTyp.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("src must be a pointer")
	}

	// get Type data of the actual struct pointed by the pointer,
	// verify that it is a struct
	srcTyp = srcTyp.Elem()
	if srcTyp.Kind() != reflect.Struct {
		return nil, fmt.Errorf("src must be a pointer to a struct")
	}

	srcVal := reflect.ValueOf(src).Elem()

	newArgSet := DefaultArgSet()
	// iterate over all fields of the struct, parse the value of 'argparse' tag
	// and create arguments accordingly. Skip any field not tagged with 'argparse'
	for i := 0; i < srcTyp.NumField(); i++ {
		fieldType := srcTyp.Field(i)
		fieldVal := srcVal.Field(i)
		tagValue, hasTag := srcTyp.Field(i).Tag.Lookup(tagKey)
		if !hasTag {
			continue
		}

		// create map of user provided tag values
		argAttrs, err := makeMapFromTags(tagValue)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing tags for field '%s': %s", fieldType.Name, err)
		}

		err = newArgSet.addArgument(fieldType, fieldVal, argAttrs)
		if err != nil {
			return nil, fmt.Errorf("Error while adding argument: %s", err)
		}
	}

	// val := reflect.ValueOf(src).Elem()
	// x := val.Field(4)
	// // ptrtype := reflect.PtrTo(x.Type())
	// addr := x.Addr().Interface()

	// // tt := addr.Convert(ptrtype)
	// fmt.Printf("\n%+v\n", NewPosArg("sdf", addr.(ArgValue), "usage"))
	// _ = NewInt(tt)
	// NewArgInfo("sdffsf", NewInt(*int(x)), "sdf")

	// for i := 1; i <= 1; i++ {
	// 	// for i := 0; i < typ.NumField(); i++ {

	// 	field := typ.Field(i)
	// 	nm := field.Tag.Get("name")
	// 	// us := field.Tag.Get("usage")
	// 	val := reflect.ValueOf(field)
	// 	fmt.Printf("\n%v---%+v\n", nm, val)
	// 	fmt.Println(val.Type().)
	// arg := NewArgInfo(nm, NewInt(val.Int()), us)
	// }
	return newArgSet, nil
}

func DefaultArgSet() *ArgSet {
	return &ArgSet{
		posArgs: make(map[string]*PosArg),
		optArgs: make(map[string]*OptArg),
	}
}

func (argset *ArgSet) AddOptional(arg *OptArg) {
	argset.optArgs[arg.common.name] = arg
}

func (argset *ArgSet) AddPositional(arg *PosArg) {
	argset.posArgs[arg.common.name] = arg
}
