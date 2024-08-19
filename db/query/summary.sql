-- name: CreateSummary :one
INSERT INTO summary (
  doc_id,
  param1,
  param2,
  summary
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetSummary :one
SELECT * FROM summary  
WHERE doc_id IS NOT NULL AND doc_id = $1 LIMIT 1;

-- name: ChangeSummary :exec
UPDATE summary SET summary = $2 WHERE doc_id = $1 RETURNING *;

-- name: DeleteSummary :exec
DELETE FROM summary WHERE summary = $1;