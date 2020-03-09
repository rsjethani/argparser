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

type ArgSet struct {
	Description string
	posArgs     []*PosArg
	optArgs     map[string]*OptArg
	// largestName int
}

func DefaultArgSet() *ArgSet {
	return &ArgSet{
		posArgs: make([]*PosArg, 0),
		optArgs: make(map[string]*OptArg),
	}
}

func (argset *ArgSet) AddOptional(arg *OptArg) {
	argset.optArgs[arg.Name] = arg
}

func (argset *ArgSet) AddPositional(arg *PosArg) {
	argset.posArgs = append(argset.posArgs, arg)
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
	for _, pos := range argset.posArgs {
		builder.WriteString("\n")
		builder.WriteString(pos.Usage())
	}
	builder.WriteString("\n\n")
	builder.WriteString("Optional Arguments:")
	for _, opt := range argset.optArgs {
		builder.WriteString("\n")
		builder.WriteString(opt.Usage())
	}
	return builder.String()
}

func (argset *ArgSet) addArgument(name string, argVal ArgValue, argAttrs map[string]string) error {
	argName := argAttrs["name"]

	argHelp := argAttrs["help"]

	// check whether user wants positional or optional argument and process accordinly
	if _, wantsPos := argAttrs["pos"]; wantsPos {
		// TODO: verify value of 'positional is yes/true only'
		argset.AddPositional(NewPosArg(argVal, argName, argHelp))

	} else { // user wants optional argument
		argset.AddOptional(NewOptArg(argVal, "--"+argName, argHelp))
	}
	return nil
}
