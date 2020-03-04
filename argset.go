package argparser

type optArg struct {
	usage string
	value Value
}

type posArg struct {
	usage string
	value Value
}

type ArgSet struct {
	Description string
	posArgs     map[string]*posArg
	optArgs     map[string]*optArg
}

func NewArgSet() *ArgSet {
	return &ArgSet{posArgs: make(map[string]*posArg), optArgs: make(map[string]*optArg)}
}

func (set *ArgSet) AddOptional(val Value, name string, usage string) {
	set.optArgs[name] = &optArg{value: val, usage: usage}
}

func (set *ArgSet) AddPositional(val Value, name string, usage string) {
	set.posArgs[name] = &posArg{value: val, usage: usage}
}
