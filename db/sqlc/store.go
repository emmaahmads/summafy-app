package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to execute db queries and transaction
type Store struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// NewDocumentTx creates a new document
type NewDocumentParams struct {
	Activity   int64  `json:"activity"`
	Username   string `json:"username"`
	IsPrivate  bool   `json:"is_private"`
	HasSummary bool   `json:"has_summary"`
	FileName   string `json:"file_name"`
	Param1     bool   `json:"param1"`
	Param2     bool   `json:"param2"`
	Summary    string `json:"summary,omitempty"`
}

type NewDocumentResults struct {
	Document Document `json:"document"`
	Summary  Summary  `json:"summary"`
	Activity Activity `json:"activity"`
}

func (store *Store) NewDocumentTx(ctx context.Context, arg NewDocumentParams) (NewDocumentResults, error) {
	var result NewDocumentResults
	var doc_id int64
	var err error
	err = store.execTx(ctx, func(q *Queries) error {
		result.Document, err = q.CreateDocument(ctx, CreateDocumentParams{
			Username:   arg.Username,
			IsPrivate:  arg.IsPrivate,
			HasSummary: arg.HasSummary,
			FileName:   arg.FileName,
		})
		if err != nil {
			return err
		}
		doc_id = result.Document.ID
		result.Summary, err = q.CreateSummary(ctx, CreateSummaryParams{
			DocID:   doc_id,
			Param1:  arg.Param1,
			Param2:  arg.Param2,
			Summary: []byte(arg.Summary),
		})
		if err != nil {
			return err
		}
		result.Activity, err = q.CreateActivity(ctx, CreateActivityParams{
			Username:   arg.Username,
			DocumentID: doc_id,
			Activity:   arg.Activity,
		})
		if err != nil {
			return err
		}

		return nil
	})
	return result, err
}
