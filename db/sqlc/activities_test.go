package db

import (
	"context"
	"testing"

	"github.com/emmaahmads/summafy/util"
	"github.com/stretchr/testify/require"
)

func TestCreateActivity(t *testing.T) {
	usr, err := createRandomUser(t)
	require.NoError(t, err)
	doc, err := testQueries.CreateDocument(context.Background(), CreateDocumentParams{
		Username:   usr.Username,
		IsPrivate:  true,
		HasSummary: true,
		FileName:   util.RandomString(6),
	})
	require.NoError(t, err)
	arg := CreateActivityParams{
		Username:   usr.Username,
		DocumentID: doc.ID,
		Activity:   util.RandomActivity(),
	}
	activity, err := testQueries.CreateActivity(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, activity)
	require.Equal(t, arg.Username, activity.Username)
	require.Equal(t, arg.DocumentID, activity.DocumentID)
	require.Equal(t, arg.Activity, activity.Activity)
	testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}
