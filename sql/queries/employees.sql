
-- name: GetEmployeByuserById :one
SELECT * FROM employees WHERE user_id = $1;

-- name: CreateEmployee :one
INSERT INTO employees
(
    user_id,
    job_title,
    country,
    salary,
    created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING * ;

-- name: UpdateEmployeeByUserId :one
UPDATE employees
SET 
    job_title  = $2,
    country    = $3,
    salary     = $4
WHERE user_id = $1
RETURNING *;                                 

-- name: DeleteEmployeeByUserId :one
DELETE FROM employees WHERE user_id = $1 RETURNING *;


-- name: GetSalaryMetricsByCountry :one
SELECT 
    ROUND(MIN(salary), 2)   AS min_salary,
    ROUND(MAX(salary), 2)   AS max_salary,
    ROUND(AVG(salary), 2)   AS avg_salary,
    COUNT(*)                AS employee_count
FROM employees
WHERE country = $1;

-- name: GetAvgSalaryPerJobTitle :one
SELECT 
    ROUND(AVG(salary), 2)   AS average_salary,
    COUNT(*)                AS employee_count
FROM employees
WHERE job_title = $1;
