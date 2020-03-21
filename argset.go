package argparser

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	stateInit int = iota
	statePosArg
	stateOptArg
	stateNoArgsLeft
	defaultOptArgPrefix string = "--"
	packageTag          string = "argparser"
)

type posArgWithName struct {
	name string
	arg  *Argument
}

type ArgSet struct {
	Description  string
	OptArgPrefix string
	posArgs      []posArgWithName
	optArgs      map[string]*Argument
	// largestName int
	// mutex
	// choices
	//short option and short prefix
	// only modify source vars if no errors ie make it atomic
	// make usage tabular: name,type/format,default,help
}

func DefaultArgSet() *ArgSet {
	return &ArgSet{
		OptArgPrefix: defaultOptArgPrefix,
		optArgs:      make(map[string]*Argument),
	}
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
	// iterate over all fields of the struct, parse the value of 'argparser' tag
	// and create arguments accordingly. Skip any field not tagged with 'argparser'
	for i := 0; i < srcTyp.NumField(); i++ {
		fieldType := srcTyp.Field(i)
		fieldVal := srcVal.Field(i)
		structTags, tagged := srcTyp.Field(i).Tag.Lookup(packageTag)
		if !tagged {
			continue
		}

		tags, err := parseTags(structTags)
		if err != nil {
			return nil, fmt.Errorf("Error while creating argument from field '%s': %s", fieldType.Name, err)
		}

		argVal, err := NewValue(fieldVal.Addr().Interface())
		if err != nil {
			return nil, fmt.Errorf("Error while creating argument from field '%s': %s", fieldType.Name, err)
		}

		arg, err := newArgFromTags(argVal, tags)
		if err != nil {
			return nil, fmt.Errorf("Error while creating argument from field '%s': %s", fieldType.Name, err)
		}

		newArgSet.AddArgument(tags["name"], arg)
	}

	return newArgSet, nil
}

func (argSet *ArgSet) Arg(name string) *Argument {
	if arg, ok := argSet.optArgs[name]; ok {
		return arg
	}
	for _, arg := range argSet.posArgs {
		if arg.name == name {
			return arg.arg
		}
	}
	return nil
}

func (argSet *ArgSet) AddArgument(name string, arg *Argument) {
	if arg.IsPositional() {
		argSet.posArgs = append(argSet.posArgs, posArgWithName{name: name, arg: arg})
		return
	}
	argSet.optArgs[argSet.OptArgPrefix+name] = arg
}

func (argSet *ArgSet) Usage() string {
	// TODO: show list of pos/opt args in sorted order
	builder := strings.Builder{}
	builder.WriteString("[] [] [] ...")
	builder.WriteString("\n\n")
	builder.WriteString(argSet.Description)
	builder.WriteString("\n\n")
	builder.WriteString("Positional Arguments:")
	for _, p := range argSet.posArgs {
		builder.WriteString("\n")
		//fmt.Sprintf("\n%-[2]*[1]s",p.Name,argSet.largestName)
		builder.WriteString(fmt.Sprintf("%-15s%s", p.name, p.arg.Help))
	}
	builder.WriteString("\n\n")
	builder.WriteString("Optional Arguments:")
	for name := range argSet.optArgs {
		builder.WriteString("\n")
		builder.WriteString(fmt.Sprintf("%-15s%s. Default: %s", name, argSet.optArgs[name].Help, argSet.optArgs[name].Value))
	}
	return builder.String()
}

func (argSet *ArgSet) ParseFrom(args []string) error {
	curState := stateInit
	var curArg string
	visited := make(map[string]bool)
	var posIndex, argsIndex int

	getArg := func(i int) string {
		if i < len(args) {
			return args[i]
		}
		return ""
	}

	for {
		switch curState {
		case stateInit:
			arg := getArg(argsIndex)
			if arg == "" {
				curState = stateNoArgsLeft
				break
			}
			curArg = arg

			// if curArg starts with the configured prefix then process it as an optional arg
			if strings.HasPrefix(curArg, argSet.OptArgPrefix) {
				if _, found := argSet.optArgs[curArg]; found {
					if visited[curArg] { // if curArg is defined but already processed then return error
						return fmt.Errorf("option '%s' already given", curArg)
					}
					curState = stateOptArg
					break
				} else { // if curArg is not defined as an opt arg then return error
					return fmt.Errorf("unknown optional argument: %s", curArg)
				}
			}

			// if all positional args have not been processed yet then consider
			// curArg as the value for next positional arg
			if posIndex < len(argSet.posArgs) {
				curState = statePosArg
				break
			}

			// since all defined positional and optional args have been processed, curArg
			// is an undefined positional arg
			return fmt.Errorf("Unknown positional argument: %s", curArg)
		case statePosArg:
			if err := argSet.posArgs[posIndex].arg.Value.Set(curArg); err != nil {
				return fmt.Errorf("error while setting option '%s': %s", argSet.posArgs[posIndex].name, err)
			}
			visited[argSet.posArgs[posIndex].name] = true
			posIndex++
			argsIndex++
			curState = stateInit
		case stateOptArg:
			if argSet.optArgs[curArg].NArgs() == 0 {
				argSet.optArgs[curArg].Value.Set()
				argsIndex++
			} else if argSet.optArgs[curArg].NArgs() < 0 {
				if err := argSet.optArgs[curArg].Value.Set(args[argsIndex+1:]...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				argsIndex = len(args)

			} else {
				inp := []string{}
				for i := 1; i <= argSet.optArgs[curArg].NArgs(); i++ {
					v := getArg(i + argsIndex)
					if v == "" {
						return fmt.Errorf("invalid no. of arguments for option '%s'; required: %d, given: %d", curArg, argSet.optArgs[curArg].NArgs(), i-1)
					}
					inp = append(inp, v)
				}
				if err := argSet.optArgs[curArg].Value.Set(inp...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				argsIndex += argSet.optArgs[curArg].NArgs() + 1
			}
			curState = stateInit
		case stateNoArgsLeft:
			for _, pos := range argSet.posArgs {
				if !visited[pos.name] {
					return fmt.Errorf("Error: value for positional argument '%s' not given", pos.name)
				}
			}
			return nil
		}
	}
}
