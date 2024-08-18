package db

import (
	"context"
	"testing"

	"github.com/emmaahmads/summafy/util"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	usr, _ := createRandomUser(t)
	err := testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}

func createRandomUser(t *testing.T) (User, error) {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.NotZero(t, user.CreatedAt)

	return user, nil
}

func TestGetUser(t *testing.T) {
	usr, err := createRandomUser(t)
	require.NoError(t, err)

	user, err := testQueries.GetUser(context.Background(), usr.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, usr.Username, user.Username)
	require.Equal(t, usr.FullName, user.FullName)
	testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}

func TestChangePassword(t *testing.T) {
	usr, err := createRandomUser(t)
	require.NoError(t, err)
	randomPassword := util.RandomString(6)
	err = testQueries.ChangePassword(context.Background(), ChangePasswordParams{usr.Username, randomPassword})
	require.NoError(t, err)
	err = testQueries.DeleteUser(context.Background(), usr.Username)
	require.NoError(t, err)
}
