


-- name: GetAdminUser :one
SELECT * FROM adminUsers WHERE user_id = $1;

-- name: CreateAdminUser :one
INSERT INTO adminUsers
(
    user_id
) VALUES (
    $1
) RETURNING * ;

-- name: DeleteAdminUser :one
DELETE FROM adminUsers
WHERE user_id = $1  RETURNING * ;