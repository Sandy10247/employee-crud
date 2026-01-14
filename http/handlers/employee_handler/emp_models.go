package employeehandler

import (
	"log"

	"server/sql/database"
)

type EmpBody struct {
	JobTitle string  `json:"job_title"`
	Country  string  `json:"country"`
	Salary   float64 `json:"salary"`
}

type Employee struct {
	ID       int32   `json:"id"`
	JobTitle string  `json:"job_title"`
	Country  string  `json:"country"`
	Salary   float64 `json:"salary"`
}

func dbEmployeeToEmpJson(dbEmp *database.Employee) Employee {
	salary, err := dbEmp.Salary.Float64Value()
	if err != nil {
		log.Printf("Error :- %v\n", err)
	}

	return Employee{
		ID:       dbEmp.ID,
		JobTitle: dbEmp.JobTitle,
		Country:  dbEmp.Country,
		Salary:   salary.Float64,
	}
}
