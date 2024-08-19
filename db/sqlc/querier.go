// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"
)

type Querier interface {
	ChangePassword(ctx context.Context, arg ChangePasswordParams) error
	ChangeSummary(ctx context.Context, arg ChangeSummaryParams) error
	CreateActivity(ctx context.Context, arg CreateActivityParams) (Activity, error)
	CreateDocument(ctx context.Context, arg CreateDocumentParams) (Document, error)
	CreateSummary(ctx context.Context, arg CreateSummaryParams) (Summary, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteSummary(ctx context.Context, summary []byte) error
	DeleteUser(ctx context.Context, username string) error
	GetAllActivities(ctx context.Context) ([]Activity, error)
	GetAllActivitiesForUser(ctx context.Context, dollar_1 interface{}) ([]Activity, error)
	GetAllDocumentsByUser(ctx context.Context, username string) ([]Document, error)
	GetAllNonPrivateDocuments(ctx context.Context) ([]Document, error)
	GetDocument(ctx context.Context, id int64) (Document, error)
	GetLastActivityForUser(ctx context.Context, dollar_1 interface{}) (Activity, error)
	GetSummary(ctx context.Context, docID int64) (Summary, error)
	GetUser(ctx context.Context, username string) (User, error)
	UpdatePrivateDocument(ctx context.Context, arg UpdatePrivateDocumentParams) (Document, error)
	UpdateSummary(ctx context.Context, arg UpdateSummaryParams) (Document, error)
}

var _ Querier = (*Queries)(nil)
