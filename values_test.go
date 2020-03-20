package argparser_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rsjethani/argparser"
)

// type point implements ArgValue interface
type point struct {
	x, y int
}

func (p *point) Set(values ...string) error {
	vals := strings.Split(values[0], ",")
	v, err := strconv.ParseInt(vals[0], 0, strconv.IntSize)
	if err != nil {
		return err
	}
	p.x = int(v)

	v, err = strconv.ParseInt(vals[1], 0, strconv.IntSize)
	p.y = int(v)

	return nil
}

// func (p *point) Get() interface{} {
// 	return p
// }

func (p *point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (p *point) IsBoolValue() bool { return false }

func TestSupportedTypeValueCreation(t *testing.T) {
	// Test value creation for types implementing ArgValue interface
	supported := []interface{}{
		new(point),
		new(int),
		new([]int),
		new(bool),
		new([]bool),
		new(uint),
		new([]uint),
		new(int64),
		new([]int64),
		new(string),
		new([]string),
		new(uint64),
		new([]uint64),
		new(float64),
		new([]float64),
		new(time.Duration),
		new([]time.Duration),
	}
	for _, val := range supported {
		_, err := argparser.NewValue(val)
		if err != nil {
			t.Errorf("Expected: NewValue(%T) should succeed, Got: %s", val, err)
		}
	}
}

func TestUnsupportedTypeValueCreation(t *testing.T) {
	type unsupported struct{}
	var x unsupported
	val, err := argparser.NewValue(&x)
	if err == nil {
		t.Errorf("Expected: unsupported type error , Got: value of %T type", val)
	}
}

func TestBoolType(t *testing.T) {
	var testVar bool
	arg := argparser.NewBool(&testVar)

	data := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	// Test valid values
	for _, val := range data {
		if err := arg.Set(val.input); err != nil {
			t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, val.input)
		}
		if val.expected != testVar {
			t.Errorf("Expected: %v, Got: %v", val.expected, testVar)
		}
		if val.input != arg.String() {
			t.Errorf("Expected: %v, Got: %v", val.input, arg.String())
		}
	}

	// Test invalid values
	for _, input := range []string{"tRUe", "hello", "1.1"} {
		if err := arg.Set(input); err == nil {
			t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
		}
	}
}

func TestBoolListType(t *testing.T) {
	var testVar []bool
	arg := argparser.NewBoolList(&testVar)
	data := []struct {
		input    []string
		expected []bool
	}{
		{[]string{"true", "false"}, []bool{true, false}},
	}

	// Test valid values
	for _, val := range data {
		arg.Set(val.input...)
		// check that all values from expected are set without error
		if err := arg.Set(val.input...); err != nil {
			t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, val.input)
		}
		// check whether each value in expected is same as set in testVar
		for i, _ := range val.expected {
			if val.expected[i] != testVar[i] {
				t.Errorf("Expected: %v, Got: %v", val.expected[i], testVar[i])
			}
		}
		// check whether string representation on input is same as that of arg
		if fmt.Sprint(val.input) != arg.String() {
			t.Errorf("Expected: %v, Got: %v", val.input, arg.String())
		}
	}

	// Test invalid values
	input := []string{"tRUe", "hello", "1.1"}
	if err := arg.Set(input...); err == nil {
		t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
	}

}

func TestIntType(t *testing.T) {
	var testVar int
	arg := argparser.NewInt(&testVar)

	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	data := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"10", 10},
		{"-10", -10},
		{fmt.Sprint(maxInt), maxInt},
		{fmt.Sprint(minInt), minInt},
	}

	// Test valid values
	for _, val := range data {
		if err := arg.Set(val.input); err != nil {
			t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, val.input)
		}
		if val.expected != testVar {
			t.Errorf("Expected: %v, Got: %v", val.expected, testVar)
		}
		if val.input != arg.String() {
			t.Errorf("Expected: %v, Got: %v", val.input, arg.String())
		}
	}

	// Test invalid values
	for _, input := range []string{"hello", "1.1", "true", "666666666666666666666666"} {
		if err := arg.Set(input); err == nil {
			t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
		}
	}
}
