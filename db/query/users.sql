-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: ChangePassword :exec
UPDATE users SET hashed_password = $2 WHERE username = $1 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE username = $1;