-- name: CreateDocument :one
INSERT INTO document (
  username,
  is_private,
  has_summary,
  file_name
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAllDocumentsByUser :many
SELECT * FROM document
WHERE username = $1;

-- name: GetAllNonPrivateDocuments :many
SELECT * FROM document
WHERE private = false;

-- name: UpdateSummary :one
UPDATE document SET has_summary = $2 WHERE id = $1 RETURNING *;

-- name: UpdatePrivateDocument :one
UPDATE document SET has_summary = $2 WHERE id = $1 RETURNING *;

-- name: GetDocument :one
SELECT * FROM document
WHERE id = $1 LIMIT 1;