package argparser

type ArgParser struct {
	Title       string
	Description string
	SubTitle    string
	mainArgSet  *ArgSet
	children    map[string]*ArgSet
}

func NewArgParser(set *ArgSet) *ArgParser {
	return &ArgParser{mainArgSet: set, SubTitle: "SubCommands:"}
}

func (parser *ArgParser) ParseFrom(args []string) {
	// parser.mainArgSet.optArgs["salute"].common.value.Set("SDFSDGSGGS")
	// parser.mainArgSet.allArgs["salutation"].Value.Set("XXX")
	// parser.mainArgSet.allArgs["salary"].Value.Set("456.66")
	// parser.mainArgSet.allArgs["tt"].Value.Set("456", "56")
	// parser.mainArgSet.allArgs["location"].Value.Set("-1,-9")
}

// func (parser *ArgParser) AddSubParser

// // IntVar defines an int flag with specified name, default value, and usage string.
// // The argument p points to an int variable in which to store the value of the flag.
// func (f *ArgSet) IntVar(p *int, name string, value int, usage string) {

// 	f.Var(newIntValue(value, p), name, usage)

// }

// // Var defines a flag with the specified name and usage string. The type and

// // value of the flag are represented by the first argument, of type Value, which

// // typically holds a user-defined implementation of Value. For instance, the

// // caller could create a flag that turns a comma-separated string into a slice

// // of strings by giving the slice the methods of Value; in particular, Set would

// // decompose the comma-separated string into the slice.

// func (f *FlagSet) Var(value Value, name string, usage string) {
// 	// Remember the default value as a string; it won't change.
// 	flag := &Flag{name, usage, value, value.String()}
// 	_, alreadythere := f.formal[name]
// 	if alreadythere {
// 		var msg string
// 		if f.name == "" {
// 			msg = fmt.Sprintf("flag redefined: %s", name)
// 		} else {
// 			msg = fmt.Sprintf("%s flag redefined: %s", f.name, name)
// 		}
// 		fmt.Fprintln(f.Output(), msg)
// 		panic(msg) // Happens only if flags are declared with identical names
// 	}
// 	if f.formal == nil {
// 		f.formal = make(map[string]*Flag)
// 	}
// 	f.formal[name] = flag
// }
