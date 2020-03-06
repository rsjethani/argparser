package argparser

import (
	"fmt"
	"reflect"
)

const UnlimitedNArgs int = -1

type ArgInfo struct {
	Value ArgValue
	Name  string
	Usage string
	NArgs int // later convert to string for patterns like '*', '+'
	// mutex   map[string]bool
	// visited bool
	//repeat bool
}

func NewArgInfo(name string, value ArgValue, usage string) *ArgInfo {
	return &ArgInfo{
		NArgs: 1,
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

type ArgSet struct {
	Title       string
	Description string
	posArgs     []*ArgInfo
	allArgs     map[string]*ArgInfo
}

func NewArgSetFrom(src interface{}) (*ArgSet, error) {
	set := NewArgSet()

	t := reflect.TypeOf(src)
	if t.Kind() != reflect.Ptr {
		return nil, fmt.Errorf("val must be a pointer")
	}

	typ := t.Elem()
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("val must be a pointer to a struct")
	}

	val := reflect.ValueOf(src).Elem()
	x := val.Field(4)
	// ptrtype := reflect.PtrTo(x.Type())
	addr := x.Addr().Interface()

	// tt := addr.Convert(ptrtype)
	fmt.Printf("\n%+v\n", NewArgInfo("sdf", addr.(ArgValue), "usage"))
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
	return set, nil
}

func NewArgSet() *ArgSet {
	return &ArgSet{posArgs: make([]*ArgInfo, 0), allArgs: make(map[string]*ArgInfo)}
}

func (argset *ArgSet) AddOptional(arg *ArgInfo) {
	temp := *arg
	argset.allArgs[arg.Name] = &temp
}

func (argset *ArgSet) AddPositional(arg *ArgInfo) {
	argset.AddOptional(arg)
	argset.posArgs = append(argset.posArgs, argset.allArgs[arg.Name])
}
