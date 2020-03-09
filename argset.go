package argparser

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	tagKey         string = "argparser"
	tagValueSep    string = ","
	mapKeyValueSep string = "="
)

type ArgSet struct {
	Description string
	posArgs     map[string]*PosArg
	optArgs     map[string]*OptArg
	sortedPos   []string
	// largestName int
}

func DefaultArgSet() *ArgSet {
	return &ArgSet{
		posArgs: make(map[string]*PosArg),
		optArgs: make(map[string]*OptArg),
	}
}

func (argSet *ArgSet) AddOptional(name string, arg *OptArg) {
	argSet.optArgs[name] = arg
}

func (argSet *ArgSet) AddPositional(name string, arg *PosArg) {
	argSet.posArgs[name] = arg
	argSet.sortedPos = append(argSet.sortedPos, name)
	sort.Strings(argSet.sortedPos)
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
			return nil, fmt.Errorf("Error while creating argument from field '%s': %s", fieldType.Name, err)

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

func (argset *ArgSet) Usage() string {
	builder := strings.Builder{}
	builder.WriteString("[] [] [] ...")
	builder.WriteString("\n\n")
	builder.WriteString(argset.Description)
	builder.WriteString("\n\n")
	builder.WriteString("Positional Arguments:")
	for _, name := range argset.sortedPos {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%-15s%s", name, argset.posArgs[name].Help))
	}
	builder.WriteString("\n\n")
	builder.WriteString("Optional Arguments:")
	for name := range argset.optArgs {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%-15s%s", name, argset.optArgs[name].Help))
	}
	return builder.String()
}

func (argset *ArgSet) addArgument(name string, argVal ArgValue, argAttrs map[string]string) error {
	argName := argAttrs["name"]

	argHelp := argAttrs["help"]

	// check whether user wants positional or optional argument and process accordinly
	if _, wantsPos := argAttrs["pos"]; wantsPos {
		// TODO: verify value of 'positional is yes/true only'
		argset.AddPositional(argName, NewPosArg(argVal, argHelp))

	} else { // user wants optional argument
		argset.AddOptional("--"+argName, NewOptArg(argVal, argHelp))
	}
	return nil
}
