package db

import (
	"context"
	"testing"

	"github.com/emmaahmads/summafy/util"
	"github.com/stretchr/testify/require"
)

func TestNewDocumentTx(t *testing.T) {

	usr, err := createRandomUser(t)
	require.NoError(t, err)

	arg := NewDocumentParams{
		Username:   usr.Username,
		IsPrivate:  true,
		HasSummary: true,
		FileName:   util.RandomString(6),
	}
	doc, err := testQueries.NewDocumentTx(context.Background(), arg)
	require.Error(t, err)
	require.NotEmpty(t, doc)
	require.Equal(t, arg.FileName, doc.Document.FileName)
	require.Equal(t, doc.Activity.Activity, Uploaded)
	require.Equal(t, doc.Summary.DocID, doc.Document.ID)

	//	err = store.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}
