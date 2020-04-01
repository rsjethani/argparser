package argparser

import (
	"fmt"
	"io"
	"os"
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
	name         string
	ArgList      []string
	Description  string
	OptArgPrefix string
	posArgs      []posArgWithName
	optArgs      map[string]*Argument
	usageOut     io.Writer
	Usage        func()

	// mutex
	// choices
	//short option and short prefix
	// only modify source vars if no errors ie make it atomic
	// make usage tabular: name,type/format,default,help
}

func (argSet *ArgSet) SetOutput(w io.Writer) {
	if w == nil {
		argSet.usageOut = os.Stderr
		return
	}
	argSet.usageOut = w
}

func (argSet *ArgSet) addHelp() {
	var help bool
	argSet.Add("help", NewSwitchArg(NewBool(&help), "Show this help message"))
}

func NewArgSet() *ArgSet {
	argSet := &ArgSet{
		OptArgPrefix: defaultOptArgPrefix,
		optArgs:      make(map[string]*Argument),
		usageOut:     os.Stderr,
		name:         os.Args[0],
		ArgList:      os.Args[1:],
	}
	argSet.addHelp()
	return argSet
}

func NewArgSetFrom(src interface{}) (*ArgSet, error) {
	if src == nil {
		return nil, fmt.Errorf("src cannot be nil")
	}
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

	newArgSet := NewArgSet()
	// iterate over all fields of the struct, parse the value of 'argparser' tag
	// and create arguments accordingly. Skip any field not tagged with 'argparser'
	for i := 0; i < srcTyp.NumField(); i++ {
		fieldType := srcTyp.Field(i)
		fieldVal := srcVal.Field(i)
		structTags, tagged := srcTyp.Field(i).Tag.Lookup(packageTag)
		if !tagged {
			continue
		}

		if !fieldVal.Addr().CanInterface() {
			return nil, fmt.Errorf("Error while creating argument from field '%s': %s", fieldType.Name, "unexported struct field")
		}
		argVal, err := NewValue(fieldVal.Addr().Interface())
		if err != nil {
			return nil, fmt.Errorf("Error while creating argument from field '%s': %s", fieldType.Name, err)
		}

		arg, name, err := newArgFromTags(argVal, fieldType.Name, structTags)
		if err != nil {
			return nil, fmt.Errorf("Error while creating argument from field '%s': %s", fieldType.Name, err)
		}

		newArgSet.Add(name, arg)
	}

	return newArgSet, nil
}

func (argSet *ArgSet) Add(name string, arg *Argument) {
	if arg == nil {
		return
	}
	if arg.positional {
		argSet.posArgs = append(argSet.posArgs, posArgWithName{name: name, arg: arg})
		return
	}
	argSet.optArgs[argSet.OptArgPrefix+name] = arg
}

// usage calls the Usage method for the ArgSet if one is specified,
// or the appropriate default usage function otherwise.
func (argSet *ArgSet) usage() {
	if argSet.Usage == nil {
		argSet.defaultUsage()
	} else {
		argSet.Usage()
	}
}

func (argSet *ArgSet) defaultUsage() {
	out := argSet.usageOut
	fmt.Fprintf(out, "Usage of %s:\n\n", argSet.name)
	// TODO: show list of opt args in sorted order
	fmt.Fprint(out, argSet.Description)
	fmt.Fprint(out, "\n\nPositional Arguments:")
	for _, p := range argSet.posArgs {
		val := p.arg.value.Get()
		fmt.Fprintf(out, "\n  %[1]s  %[2]T\n\t%[3]s  (Default: %[2]v)", p.name, val, p.arg.help, val)
	}

	fmt.Fprint(out, "\n\nOptional Arguments:")
	for name, arg := range argSet.optArgs {
		sw := ""
		if arg.isSwitch() {
			sw = "  (switch)"
		}
		val := arg.value.Get()
		fmt.Fprintf(out, "\n  %[1]s  %[2]T%[4]s\n\t%[3]s  (Default: %[2]v)", name, val, arg.help, sw)
	}

	fmt.Fprintln(out, "")
}

func (argSet *ArgSet) Parse() error {
	argsToParse := argSet.ArgList
	curState := stateInit
	var curArg string
	visited := make(map[string]bool)
	var posIndex, argsIndex int

	getArg := func(i int) string {
		if i < len(argsToParse) {
			return argsToParse[i]
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
			if err := argSet.posArgs[posIndex].arg.value.Set(curArg); err != nil {
				return fmt.Errorf("error while setting option '%s': %s", argSet.posArgs[posIndex].name, err)
			}
			visited[argSet.posArgs[posIndex].name] = true
			posIndex++
			argsIndex++
			curState = stateInit
		case stateOptArg:
			if argSet.optArgs[curArg].nArgs == 0 {
				if curArg == "--help" {
					argSet.usage()
					return nil
				}
				argSet.optArgs[curArg].value.Set()
				argsIndex++
			} else if argSet.optArgs[curArg].nArgs < 0 {
				if err := argSet.optArgs[curArg].value.Set(argsToParse[argsIndex+1:]...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				argsIndex = len(argsToParse)

			} else {
				inp := []string{}
				for i := 1; i <= argSet.optArgs[curArg].nArgs; i++ {
					v := getArg(i + argsIndex)
					if v == "" {
						return fmt.Errorf("invalid no. of arguments for option '%s'; required: %d, given: %d", curArg, argSet.optArgs[curArg].nArgs, i-1)
					}
					inp = append(inp, v)
				}
				if err := argSet.optArgs[curArg].value.Set(inp...); err != nil {
					return fmt.Errorf("error while setting option '%s': %s", curArg, err)
				}
				argsIndex += argSet.optArgs[curArg].nArgs + 1
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
