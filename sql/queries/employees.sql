
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
