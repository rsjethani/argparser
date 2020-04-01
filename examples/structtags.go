package main

import (
	"fmt"

	"github.com/rsjethani/argparser"
)

func main() {
	config := struct {
		Salute   string  `argparser:"help=Salutation for the employee"`
		Salary   float64 `argparser:"help=Employee salary,type=pos"`
		FullName string  `argparser:"help=Full name of the employee,type=pos"`
		EmpID    []int   `argparser:"name=emp-id,help=Employee ID for new employee,nargs=3"`
		IsIntern bool    `argparser:"name=is-intern,help=Is the new employee an intern,type=switch"`
	}{
		EmpID:  []int{-1},
		Salute: "Mr.",
	}

	mainSet, err := argparser.NewArgSetFrom(&config)
	if err != nil {
		fmt.Println(err)
		return
	}
	mainSet.Description = "CLI for managing employee database"

	mainSet.ArgList = []string{"--help"}
	mainSet.Parse()

	fmt.Printf("\nBEFORE parsing: %+v\n", config)
	mainSet.ArgList = []string{"3.4", "asd", "--salute", "XXX", "--is-intern", "--emp-id", "88888", "345", "33"}
	err = mainSet.Parse()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("\nAFTER parsing: %+v\n", config)
}
