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

func parseStructTag(tagValue string) (map[string]string, error) {
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

		argVal, err := NewArgValue(fieldVal.Addr().Interface())
		if err != nil {
			return nil, err
		}

		// create map of user provided tag values
		argAttrs, err := parseStructTag(tagValue)
		if err != nil {
			return nil, fmt.Errorf("Error while parsing tags for field '%s': %s", fieldType.Name, err)
		}

		err = newArgSet.addArgument(fieldType.Name, argVal, argAttrs)
		if err != nil {
			return nil, fmt.Errorf("Error while adding argument: %s", err)
		}
	}

	return newArgSet, nil
}

func (argset *ArgSet) addArgument(name string, argVal ArgValue, argAttrs map[string]string) error {
	//TODO: convert field name with '-' if multiple words
	argName := name
	if val, ok := argAttrs["name"]; ok {
		argName = val
	}

	var argUsage string
	if val, ok := argAttrs["usage"]; ok {
		argUsage = val
	}

	// check whether user wants positional or optional argument and process accordinly
	if _, wantsPos := argAttrs["pos"]; wantsPos {
		// TODO: verify value of 'positional is yes/true only'
		argset.AddPositional(NewPosArg(argName, argVal, argUsage))

	} else { // user wants optional argument
		argset.AddOptional(NewOptArg(argName, argVal, argUsage))
	}
	return nil
}
