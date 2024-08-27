package db

import (
	"context"
	"testing"

	"github.com/emmaahmads/summafy/util"
	"github.com/stretchr/testify/require"
)

func TestCreateSummary(t *testing.T) {
	usr, err := createRandomUser(t)
	require.NoError(t, err)
	doc, err := testQueries.CreateDocument(context.Background(), CreateDocumentParams{
		Username:   usr.Username,
		IsPrivate:  true,
		HasSummary: true,
		FileName:   util.RandomString(6),
	})
	require.NoError(t, err)
	arg := CreateSummaryParams{
		DocID:   doc.ID,
		Param1:  true,
		Param2:  true,
		Summary: util.RandomSummary(),
	}
	summary, err := testQueries.CreateSummary(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, summary)
	require.Equal(t, arg.DocID, summary.DocID)
	require.Equal(t, arg.Param1, summary.Param1)
	require.Equal(t, arg.Param2, summary.Param2)
	require.Equal(t, arg.Summary, summary.Summary)
	err = testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
	err = testQueries.DeleteSummary(context.Background(), summary.DocID)
	require.NoError(t, err)
}

func TestChangeSummary(t *testing.T) {
	usr, err := createRandomUser(t)
	require.NoError(t, err)
	doc, err := testQueries.CreateDocument(context.Background(), CreateDocumentParams{
		Username:   usr.Username,
		IsPrivate:  true,
		HasSummary: true,
		FileName:   util.RandomString(6),
	})
	require.NoError(t, err)
	arg := CreateSummaryParams{
		DocID:   doc.ID,
		Param1:  true,
		Param2:  true,
		Summary: util.RandomSummary(),
	}
	summary, err := testQueries.CreateSummary(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, summary)

	arg2 := ChangeSummaryParams{
		DocID:   doc.ID,
		Summary: util.RandomSummary(),
	}
	err = testQueries.ChangeSummary(context.Background(), arg2)
	require.NoError(t, err)
	summary2, err := testQueries.GetSummary(context.Background(), doc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, summary2)
	require.Equal(t, arg2.DocID, summary2.DocID)
	require.Equal(t, arg2.Summary, summary2.Summary)
	err = testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
	err = testQueries.DeleteSummary(context.Background(), summary2.DocID)
	require.NoError(t, err)
}

func TestDeleteSummary(t *testing.T) {
	usr, err := createRandomUser(t)
	require.NoError(t, err)
	doc, err := testQueries.CreateDocument(context.Background(), CreateDocumentParams{
		Username:   usr.Username,
		IsPrivate:  true,
		HasSummary: true,
		FileName:   util.RandomString(6),
	})
	require.NoError(t, err)
	arg := CreateSummaryParams{
		DocID:   doc.ID,
		Param1:  true,
		Param2:  true,
		Summary: util.RandomSummary(),
	}
	summary, err := testQueries.CreateSummary(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, summary)
	err = testQueries.DeleteSummary(context.Background(), summary.DocID)
	require.NoError(t, err)
	require.Equal(t, doc.ID, summary.DocID)
	err = testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}
