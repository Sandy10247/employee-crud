package employeehandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"server/http/middleware"
	"server/http/response"
	"server/sql/database"

	db "server/init"

	"github.com/jackc/pgx/v5/pgtype"
)

func floatToNumeric(f float64, precision int) (pgtype.Numeric, error) {
	var numericValue pgtype.Numeric

	// Format the float to a string with the desired precision
	// 'f' format specifier, precision specifies the number of digits after the decimal point
	str := strconv.FormatFloat(f, 'f', precision, 64)

	// Use Scan to parse the string into the pgtype.Numeric struct
	if err := numericValue.Scan(str); err != nil {
		return pgtype.Numeric{}, fmt.Errorf("failed to scan string to pgtype.Numeric: %w", err)
	}

	return numericValue, nil
}

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
	salaryNumeric, err := floatToNumeric(reqBody.Salary, 2)
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
	salaryNumeric, err := floatToNumeric(reqBody.Salary, 2)
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
