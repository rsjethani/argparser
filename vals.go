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

type IntList []int

func NewIntList(p *[]int) *IntList {
	return (*IntList)(p)
}

func (il *IntList) Set(values ...string) error {
	*il = make([]int, len(values))
	for i, val := range values {
		n, err := strconv.ParseInt(val, 0, strconv.IntSize)
		(*il)[i] = int(n)
		if err != nil {
			return numError(err)
		}
	}
	return nil
}

func (il *IntList) Get() interface{} { return (*[]int)(il) }

func (il *IntList) String() string { return fmt.Sprint(*il) }

// --- int value
type Int int

func NewInt(p *int) *Int {
	return (*Int)(p)
}

func (i *Int) Set(s ...string) error {
	v, err := strconv.ParseInt(s[0], 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = Int(v)
	return err
}

func (i *Int) Get() interface{} { return int(*i) }

func (i *Int) String() string { return strconv.Itoa(int(*i)) }

// -- string Value
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

// -- float64 Value
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
