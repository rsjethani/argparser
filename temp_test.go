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
	p.y = int(v)

	return nil
}

func (p *point) Get() interface{} {
	return p
}
func (p *point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (p *point) IsSwitch() bool { return false }

func TestTraditionalApproach(t *testing.T) {
	config := struct {
		Salute   string  `argparser:"help=Sal\\,utationfortheemployee"`
		Salary   float64 `argparser:"help=Emp\\,loyee salary,type=pos"`
		FullName string  `argparser:",help=Full name of the employee,type=pos"`
		EmpID    []int   `argparser:"name=emp-id,help=Employee ID for new employee,nargs=3"`
		Loc      point   `argparser:"name=point,help=coordinates"`
		IsIntern bool    `argparser:"name=is-intern,help=Is the new employee an intern,type=switch"`
		// IsIntern bool    `argparser:"name=is-intern|help=Is the new employee an intern|nargs=99999999999999999999"`
	}{
		EmpID:  []int{-1},
		Salute: "Mr.",
	}

	mainSet, err := NewArgSet(&config)
	if err != nil {
		t.Fatal(err)
	}
	mainSet.Description = "CLI for managing employee database"

	parser := NewArgParser(mainSet)
	parser.Usage()

	fmt.Printf("\nBEFORE parsing: %+v\n", config)
	fmt.Println(parser.ParseFrom([]string{"3.4", "asd", "--salute", "XXX", "--point", "5,-7", "--is-intern", "--emp-id", "88888", "345", "33"}))
	fmt.Printf("\nAFTER parsing: %+v\n", config)

}
