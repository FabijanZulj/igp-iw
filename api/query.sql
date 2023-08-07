-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id;

-- name: CreateUser :one
INSERT INTO users (
  email, password, isVerified
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET isVerified = $1
WHERE email = $2
RETURNING *;

-- name: GetVerificationCode :one
SELECT v.userid, v.code FROM verifyData v
INNER JOIN users u on v.userid = u.id
WHERE u.email = $1;


-- name: CreateVerifyData :one
INSERT INTO verifyData (
  userId, code
) VALUES (
  $1, $2
)
RETURNING *;
