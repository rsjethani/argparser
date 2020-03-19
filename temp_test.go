package argparser

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type point struct {
	x, y int
}

func (p *point) Set(values ...string) error {
	vals := strings.Split(values[0], ",")
	v, err := strconv.ParseInt(vals[0], 0, strconv.IntSize)
	if err != nil {
		return formatParseError(values[0], fmt.Sprintf("%T", *p), err)
	}
	p.x = int(v)

	v, err = strconv.ParseInt(vals[1], 0, strconv.IntSize)
	// if err != nil {
	// 	rerturn numError(err)
	// }
	p.y = int(v)

	return nil
}

func (p *point) Get() interface{} {
	return p
}
func (p *point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (p *point) IsBoolValue() bool { return false }

func TestTraditionalApproach(t *testing.T) {
	config := struct {
		Salute   string  `argparser:"help=Salutation for the employee|name=salute"`
		Salary   float64 `argparser:"help=Employee salary|pos=yes|name=salary"`
		FullName string  `argparser:"pos=yes|name=full-name|help=Full name of the employee|pos=yes"`
		EmpID    []int   `argparser:"name=emp-id|help=Employee ID for new employee|nargs=-45"`
		Loc      point   `argparser:"name=point"`
		// IsIntern bool    `argparser:"name=is-intern"`
	}{
		// EmpID:  -1,
		Salute: "Mr.",
	}

	mainSet, err := NewArgSet(&config)
	if err != nil {
		t.Fatal(err)
	}
	mainSet.Description = "CLI for managing employee database\n...\n..."

	fmt.Printf("\nmainset: %+v\n", mainSet)
	// fmt.Println(mainSet.Usage())

	fmt.Printf("\n%+v\n", config)
	fmt.Println(mainSet.ParseFrom([]string{"3.4", "asd", "--salute", "XXX", "--point", "--emp-id", "88888", "345", "-35"}))
	fmt.Printf("\n%+v\n", config)
}
