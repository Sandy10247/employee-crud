package employeehandler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"server/http/helper"
	"server/http/middleware"
	"server/http/response"
	"server/sql/database"

	db "server/init"
)

func CreateEmp(w http.ResponseWriter, r *http.Request) {
	var reqBody EmpBody

	// Decode the request body into the struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		response.RespondeWithError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	// Extract UserInfo from context
	userInfo, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		response.RespondeWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	// Extract Salary to pgtype.Numeric
	salaryNumeric, err := helper.FloatToNumeric(reqBody.Salary, 2)
	if err != nil {
		response.RespondeWithError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	// Create Obj for DB Insertion
	createEmp := database.CreateEmployeeParams{
		UserID:   userInfo.ID,
		JobTitle: reqBody.JobTitle,
		Country:  reqBody.Country,
		Salary:   salaryNumeric,
	}

	// Create Employee Query Call
	empCreated, err := db.Queries.CreateEmployee(r.Context(), createEmp)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldnot create Employee %v", err))
		return
	}

	response.RespondeWithJSON(w, http.StatusCreated, dbEmployeeToEmpJson(empCreated))
}

func UpdateEmp(w http.ResponseWriter, r *http.Request) {
	var reqBody EmpBody

	// Decode the request body into the struct
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		response.RespondeWithError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	// Extract UserInfo req content
	userInfo, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		response.RespondeWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	// Extract Salary to pgtype.Numeric
	salaryNumeric, err := helper.FloatToNumeric(reqBody.Salary, 2)
	if err != nil {
		response.RespondeWithError(w, http.StatusUnprocessableEntity, "invalid json")
		return
	}

	// Create UpdateEmployeeByUserIdParams
	updateEmp := database.UpdateEmployeeByUserIdParams{
		UserID:   userInfo.ID,
		JobTitle: reqBody.JobTitle,
		Country:  reqBody.Country,
		Salary:   salaryNumeric,
	}

	// Update Employee
	empCreated, err := db.Queries.UpdateEmployeeByUserId(r.Context(), updateEmp)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldnot create Employee %v", err))
		return
	}

	response.RespondeWithJSON(w, http.StatusOK, dbEmployeeToEmpJson(empCreated))
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract UserInfo from context
	userInfo, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		response.RespondeWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	emp, err := db.Queries.GetEmployeByuserById(r.Context(), userInfo.ID)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldnot fetch Employee %v", err))
		return
	}

	response.RespondeWithJSON(w, http.StatusOK, dbEmployeeToEmpJson(emp))
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	// Extract UserInfo req content
	userInfo, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		response.RespondeWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	// Delete Employee
	emp, err := db.Queries.DeleteEmployeeByUserId(r.Context(), userInfo.ID)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldnot Delete Employee %v", err))
		return
	}

	response.RespondeWithJSON(w, http.StatusOK, dbEmployeeToEmpJson(emp))
}

func NetSalary(w http.ResponseWriter, r *http.Request) {
	// Extract UserInfo from context
	userInfo, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		response.RespondeWithError(w, http.StatusBadRequest, "user not found")
		return
	}

	// Fetch Employee Details from DB
	emp, err := db.Queries.GetEmployeByuserById(r.Context(), userInfo.ID)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldnot fetch Employee %v", err))
		return
	}

	grossSalary, err := emp.Salary.Float64Value()
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, "Parsing Float Error")
		return
	}

	// Calculate all Salary types
	countryTaxRate := helper.GetTaxRatePerCountry(emp.Country)
	taxAnount := helper.CalculatePercentage(countryTaxRate, grossSalary.Float64)
	takeHomeSalary := helper.CalculateNetSalary(grossSalary.Float64, taxAnount)

	// construct the payload
	resp := map[string]float64{
		"real_salary":      grossSalary.Float64,
		"government_cut":   taxAnount,
		"take_home_salary": takeHomeSalary,
	}

	// Send the Responses
	response.RespondeWithJSON(w, http.StatusOK, resp)
}

// Admin Route
func GetSalaryMetricsByCountry(w http.ResponseWriter, r *http.Request) {
	// extract country from Query
	country := r.URL.Query().Get("country")

	// Delete Employee
	salaryMetrics, err := db.Queries.GetSalaryMetricsByCountry(r.Context(), country)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldnot Delete Employee %v", err))
		return
	}

	response.RespondeWithJSON(w, http.StatusOK, salaryMetrics)
}

func GetAvgSalaryPerJobTitle(w http.ResponseWriter, r *http.Request) {
	// extract country from Query
	JobTitle := r.URL.Query().Get("job_title")

	// Delete Employee
	AvgSalaryResp, err := db.Queries.GetAvgSalaryPerJobTitle(r.Context(), JobTitle)
	if err != nil {
		response.RespondeWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldnot Delete Employee %v", err))
		return
	}

	response.RespondeWithJSON(w, http.StatusOK, AvgSalaryResp)
}
