// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: activities.sql

package db

import (
	"context"
	"time"
)

const createActivity = `-- name: CreateActivity :one
INSERT INTO activities (
  username,
  activity,
  document_id,
  created_at
 ) VALUES (
  $1, $2, $3, $4
) RETURNING id, username, activity, created_at, document_id
`

type CreateActivityParams struct {
	Username   string    `json:"username"`
	Activity   int64     `json:"activity"`
	DocumentID int64     `json:"document_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (q *Queries) CreateActivity(ctx context.Context, arg CreateActivityParams) (Activity, error) {
	row := q.db.QueryRow(ctx, createActivity,
		arg.Username,
		arg.Activity,
		arg.DocumentID,
		arg.CreatedAt,
	)
	var i Activity
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Activity,
		&i.CreatedAt,
		&i.DocumentID,
	)
	return i, err
}

const getAllActivitiesForUser = `-- name: GetAllActivitiesForUser :many
SELECT id, username, activity, created_at, document_id FROM activities
WHERE user = $1
`

func (q *Queries) GetAllActivitiesForUser(ctx context.Context, dollar_1 interface{}) ([]Activity, error) {
	rows, err := q.db.Query(ctx, getAllActivitiesForUser, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Activity{}
	for rows.Next() {
		var i Activity
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Activity,
			&i.CreatedAt,
			&i.DocumentID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastActivityForUser = `-- name: GetLastActivityForUser :one
SELECT id, username, activity, created_at, document_id FROM activities
WHERE user = $1 AND id = (SELECT MAX(id) FROM activities WHERE user = $1) LIMIT 1
`

func (q *Queries) GetLastActivityForUser(ctx context.Context, dollar_1 interface{}) (Activity, error) {
	row := q.db.QueryRow(ctx, getLastActivityForUser, dollar_1)
	var i Activity
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Activity,
		&i.CreatedAt,
		&i.DocumentID,
	)
	return i, err
}
