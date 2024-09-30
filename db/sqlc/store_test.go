package db

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/emmaahmads/summafy/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestNewDocumentTx(t *testing.T) {
	testDB, err := sql.Open("postgres", "postgresql://emma:happybirthday@localhost:5432/summafy?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := NewStore(testDB)

	usr, err := createRandomUser(t)
	require.NoError(t, err)

	arg := NewDocumentParams{
		Username:   usr.Username,
		IsPrivate:  true,
		HasSummary: true,
		FileName:   util.RandomString(6),
		Summary:    util.RandomString(10),
	}
	doc, err := store.NewDocumentTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, doc)
	require.Equal(t, arg.FileName, doc.Document.FileName)
	require.Equal(t, int(doc.Activity.Activity), Uploaded)
	require.Equal(t, doc.Summary.DocID, doc.Document.ID)

	err = store.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}
