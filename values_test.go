package argparser_test

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rsjethani/argparser"
)

const (
	// uintSize      = 32 << (^uint(0) >> 32 & 1)
	minUint uint = 0
	maxUint uint = ^minUint
	maxInt  int  = int(maxUint >> 1)
	minInt  int  = -maxInt - 1
)

// type point implements Value interface
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

func (p *point) IsSwitch() bool { return false }

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

	// Test Set() with no arguments
	arg.Set()
	if testVar != true {
		t.Errorf("Expected: true, Got: %v", testVar)
	}

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
	data := struct {
		input    []string
		expected []bool
	}{
		input:    []string{"true", "false"},
		expected: []bool{true, false},
	}

	// Test valid values
	// check that all values from expected are set without error
	if err := arg.Set(data.input...); err != nil {
		t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, data.input)
	}
	// check whether each value in expected is same as set in testVar
	for i, _ := range data.expected {
		if data.expected[i] != testVar[i] {
			t.Errorf("Expected: %v, Got: %v", data.expected[i], testVar[i])
		}
	}
	// check whether string representation on input is same as that of arg
	if fmt.Sprint(data.input) != arg.String() {
		t.Errorf("Expected: %v, Got: %v", data.input, arg.String())
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

func TestIntListType(t *testing.T) {
	var testVar []int
	arg := argparser.NewIntList(&testVar)
	data := struct {
		input    []string
		expected []int
	}{
		input:    []string{"0", "10", "-10", fmt.Sprint(maxInt), fmt.Sprint(minInt)},
		expected: []int{0, 10, -10, maxInt, minInt},
	}

	// Test valid values
	// check that all values from expected are set without error
	if err := arg.Set(data.input...); err != nil {
		t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, data.input)
	}
	// check whether each value in expected is same as set in testVar
	for i, _ := range data.expected {
		if data.expected[i] != testVar[i] {
			t.Errorf("Expected: %v, Got: %v", data.expected[i], testVar[i])
		}
	}
	// check whether string representation on input is same as that of arg
	if fmt.Sprint(data.input) != arg.String() {
		t.Errorf("Expected: %v, Got: %v", data.input, arg.String())
	}

	// Test invalid values
	input := []string{"hello", "100", "true", "666666666666666666666666"}
	if err := arg.Set(input...); err == nil {
		t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
	}
}

func TestFloat64Type(t *testing.T) {
	var testVar float64
	arg := argparser.NewFloat64(&testVar)

	data := []struct {
		input    string
		expected float64
	}{
		{"0", 0},
		{"100", 100.00},
		{"10.11", 10.11},
		{"-10.11", -10.11},
		{fmt.Sprint(math.MaxFloat64), math.MaxFloat64},
		{fmt.Sprint(math.SmallestNonzeroFloat64), math.SmallestNonzeroFloat64},
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
	for _, input := range []string{"hello", "1.1xx", "true", "100abcd"} {
		if err := arg.Set(input); err == nil {
			t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
		}
	}
}

func TestFloat64ListType(t *testing.T) {
	var testVar []float64
	arg := argparser.NewFloat64List(&testVar)
	data := struct {
		input    []string
		expected []float64
	}{
		input:    []string{"0", "100", "10.11", "-10.11", fmt.Sprint(math.MaxFloat64), fmt.Sprint(math.SmallestNonzeroFloat64)},
		expected: []float64{0, 100.00, 10.11, -10.11, math.MaxFloat64, math.SmallestNonzeroFloat64},
	}

	// Test valid values
	// check that all values from expected are set without error
	if err := arg.Set(data.input...); err != nil {
		t.Errorf("Expected: no error, Got: error '%s' for input \"%s\"", err, data.input)
	}
	// check whether each value in expected is same as set in testVar
	for i, _ := range data.expected {
		if data.expected[i] != testVar[i] {
			t.Errorf("Expected: %v, Got: %v", data.expected[i], testVar[i])
		}
	}
	// check whether string representation on input is same as that of arg
	if fmt.Sprint(data.input) != arg.String() {
		t.Errorf("Expected: %v, Got: %v", data.input, arg.String())
	}

	// Test invalid values
	input := []string{"hello", "1.1", "true", "66666666666"}
	if err := arg.Set(input...); err == nil {
		t.Errorf("Expected: error, Got: no error for input \"%s\"", input)
	}
}
