package argparser

import (
	"errors"
	"fmt"
	"strconv"
)

type ArgValue interface {
	String() string
	Set(...string) error
	Get() interface{}
}

// errParse is returned by Set if a flag's value fails to parse, such as with an invalid integer for Int.
// It then gets wrapped through failf to provide more information.
var errParse = errors.New("parse error")

// errRange is returned by Set if a flag's value is out of range.
// It then gets wrapped through failf to provide more information.
var errRange = errors.New("value out of range")

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}
	if ne.Err == strconv.ErrSyntax {
		return errParse
	}
	if ne.Err == strconv.ErrRange {
		return errRange
	}
	return err
}

// Int type represents an int value
type Int int

func NewInt(p *int) *Int {
	return (*Int)(p)
}

func (i *Int) Set(values ...string) error {
	v, err := strconv.ParseInt(values[0], 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = Int(v)
	return err
}

func (i *Int) Get() interface{} { return int(*i) }

func (i *Int) String() string { return strconv.Itoa(int(*i)) }

// IntList type representing a list of integer values
type IntList []int

func NewIntList(p *[]int) *IntList {
	return (*IntList)(p)
}

func (il *IntList) Set(values ...string) error {
	*il = make([]int, len(values))
	for i, val := range values {
		n, err := strconv.ParseInt(val, 0, strconv.IntSize)
		if err != nil {
			return numError(err)
		}
		(*il)[i] = int(n)
	}
	return nil
}

func (il *IntList) Get() interface{} { return (*[]int)(il) }

func (il *IntList) String() string { return fmt.Sprint(*il) }

// String type represents a string value and satisfies ArgValue interface
type String string

func NewString(p *string) *String {
	return (*String)(p)
}

func (s *String) Set(val ...string) error {
	*s = String(val[0])
	return nil
}

func (s *String) Get() interface{} { return string(*s) }

func (s *String) String() string { return string(*s) }

// StringList type represents a list string value and satisfies ArgValue interface
type StringList []string

func NewStringList(p *[]string) *StringList {
	return (*StringList)(p)
}

func (sl *StringList) Set(values ...string) error {
	*sl = make([]string, len(values))
	for i, val := range values {
		(*sl)[i] = val
	}
	return nil
}

func (sl *StringList) Get() interface{} { return (*[]string)(sl) }

func (sl *StringList) String() string { return fmt.Sprint(*sl) }

// Float64 represents a float64 value and also satisfies ArgValue interface
type Float64 float64

func NewFloat64(p *float64) *Float64 {
	return (*Float64)(p)
}

func (f *Float64) Set(s ...string) error {
	v, err := strconv.ParseFloat(s[0], 64)
	if err != nil {
		err = numError(err)
	}
	*f = Float64(v)
	return err
}

func (f *Float64) Get() interface{} { return float64(*f) }

func (f *Float64) String() string { return strconv.FormatFloat(float64(*f), 'g', -1, 64) }

// Float64List type representing a list of float64 values and satisfies ArgValue interface
type Float64List []float64

func NewFloat64List(p *[]float64) *Float64List {
	return (*Float64List)(p)
}

func (fl *Float64List) Set(values ...string) error {
	*fl = make([]float64, len(values))
	for i, val := range values {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return numError(err)
		}
		(*fl)[i] = float64(f)
	}
	return nil
}

func (fl *Float64List) Get() interface{} { return (*[]float64)(fl) }

func (fl *Float64List) String() string { return fmt.Sprint(*fl) }

//---
func NewValue(v interface{}) (ArgValue, error) {
	switch addr := v.(type) {
	case *string:
		return NewString(addr), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", addr)
	}
}
