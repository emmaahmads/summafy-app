// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	ChangePassword(ctx context.Context, arg ChangePasswordParams) (User, error)
	CreateActivity(ctx context.Context, activity int64) (Activity, error)
	CreateDocument(ctx context.Context, arg CreateDocumentParams) (Document, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetAllActivitiesForUser(ctx context.Context, dollar_1 interface{}) ([]Activity, error)
	GetAllDocumentsByUser(ctx context.Context, username string) ([]Document, error)
	GetAllNonPrivateDocuments(ctx context.Context) ([]Document, error)
	GetDocument(ctx context.Context, id int64) (Document, error)
	GetLastActivityForUser(ctx context.Context, dollar_1 interface{}) (Activity, error)
	GetUser(ctx context.Context, username string) (User, error)
	UpdatePrivateDocument(ctx context.Context, arg UpdatePrivateDocumentParams) (Document, error)
	UpdateSummary(ctx context.Context, arg UpdateSummaryParams) (Document, error)
}

var _ Querier = (*Queries)(nil)
