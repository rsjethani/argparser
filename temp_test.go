package argparser_test

import (
	"fmt"
	"testing"

	"github.com/rsjethani/argparser"
)

func TestTraditionalApproach(t *testing.T) {
	config := struct {
		salute   argparser.StringVal `name:"salute" opt:"yes" usage:"Salutaion for new employee"`
		empID    argparser.Int       `name:"emp-id" opt:"yes" usage:"Employee ID for new employee" short:"i"`
		fullName argparser.StringVal `name:"full-name" usage:"Full name of the employee"`
		salary   argparser.Float64   `name:"salary" usage:"Employee salary"`
	}{empID: -1, salute: "Mr."}

	fmt.Printf("\n%+v\n", config)

	mainSet := argparser.NewArgSet(&config)
	mainSet.Description = "CLI for managing employee database"
	// mainSet.AddOptional(&config.salute, "salutation", "Salutation to use for the new employee")
	// mainSet.AddOptional(&config.empID, "employee-id", "Employee ID of the new employee")
	// mainSet.AddPositional(&config.fullName, "full-name", "Name of the new employee")
	// mainSet.AddPositional(&config.salary, "salary", "Salary of the new employee")
	// fmt.Printf("\n%#v\n", mainSet)

	// parser := argparser.NewArgParser(mainSet)
	// parser.ParseFrom([]string{})

	// fmt.Printf("\n%+v\n", config)

	// mainSet.AddPositional("age", argparser.IntValue(18))

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
