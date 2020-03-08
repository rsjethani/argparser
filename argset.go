package argparser

import (
	"fmt"
	"reflect"
)

type ArgSet struct {
	Title       string
	Description string
	posArgs     map[string]*PosArg
	optArgs     map[string]*OptArg
}

func NewArgSet(src interface{}) (*ArgSet, error) {
	set := DefaultArgSet()

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
	fmt.Printf("\n%+v\n", NewPosArg("sdf", addr.(ArgValue), "usage"))
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
