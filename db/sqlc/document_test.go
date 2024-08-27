package db

import (
	"context"
	"testing"

	"github.com/emmaahmads/summafy/util"
	"github.com/stretchr/testify/require"
)

func TestUpdatePrivateDocument(t *testing.T) {
	usr, _ := createRandomUser(t)
	doc, err := testQueries.CreateDocument(context.Background(), CreateDocumentParams{
		Username:   usr.Username,
		IsPrivate:  true,
		HasSummary: true,
		FileName:   util.RandomString(6),
	})
	require.NoError(t, err)
	_, err = testQueries.UpdatePrivateDocument(context.Background(), UpdatePrivateDocumentParams{
		ID:        doc.ID,
		IsPrivate: false,
	})
	require.NoError(t, err)
	doc2, err := testQueries.GetDocument(context.Background(), doc.ID)
	require.NoError(t, err)
	require.Equal(t, false, doc2.IsPrivate)
	err = testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}
