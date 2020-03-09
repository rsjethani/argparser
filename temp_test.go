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
		return numError(err)
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

func TestTraditionalApproach(t *testing.T) {
	config := struct {
		Salute string `argparser:""`
		// Salute   string `argparser:"name=salute,usage=Salutation for the employee"`
		FullName string `argparser:"name=full-name,usage=Full name of the employee,pos=yes"`
		// EmpID    int     `name:"emp-id" opt:"yes" usage:"Employee ID for new employee" short:"i"`
		// salary   float64 `name:"salary" usage:"Employee salary"`
		// Loc      point
	}{
		// EmpID:  -1,
		Salute: "Mr.",
	}

	// config.Loc = point{11, 22}

	fmt.Printf("\nBEFORE: %+v\n", config)
	fmt.Printf("\nSalute: %p\n", &config.Salute)
	fmt.Printf("\nFullName: %p\n", &config.FullName)

	mainSet, err := NewArgSet(&config)
	if err != nil {
		t.Fatal(err)
	}
	mainSet.Description = "CLI for managing employee database"

	fmt.Printf("\nmainset: %#v\n", mainSet)

	parser := NewArgParser(mainSet)
	parser.ParseFrom([]string{})

	fmt.Printf("\nAFTER: %+v\n", config)

}

/*
func TestTagApproach(t *testing.T) {
	myArgs := struct{
		Age argparser.IntValue	`opt:"yes" name:"user-age" short:"a"`
		Salary argparser.IntValue	`opt:"yes" name:"user-salary" short:"s"`
	}

	parser := argparser.NewArgParser(&myArgs)
	parser.ParseFrom([]string{"-a", "33", "-s", "50000"})

	fmt.Printf("\n%+v\n",myargs)
}
*/