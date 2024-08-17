-- name: CreateActivity :one
INSERT INTO activities (
  activity
 ) VALUES (
  $1
) RETURNING *;

-- name: GetLastActivityForUser :one
SELECT * FROM activities
WHERE user = $1 AND id = (SELECT MAX(id) FROM activities WHERE user = $1) LIMIT 1;

-- name: GetAllActivitiesForUser :many
SELECT * FROM activities
WHERE user = $1;