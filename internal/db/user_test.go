package db

import (
	"context"
	"github.com/kavinddd/mangtoon_be/pkg/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Username: util.RandomUsername(),
		Email:    util.RandomEmail(),
		Password: "012356914912391293",
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.Password, user.Password)

	require.NotNil(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user

}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUserById(t *testing.T) {
	account1 := createRandomUser(t)
	account2, err := testQueries.GetUserById(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.CreatedAt, account2.CreatedAt)
	require.Equal(t, account1.Username, account2.Username)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.Password, account2.Password)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
	require.False(t, account2.IsActive)
}

func TestUpdateUserEmail(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateUserEmailParams{
		ID:    user1.ID,
		Email: util.RandomEmail(),
	}
	err := testQueries.UpdateUserEmail(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateUserIsActive(t *testing.T) {
	// new user will have in_active = false by default
	user1 := createRandomUser(t)
	arg := UpdateUserIsActiveParams{
		ID:       user1.ID,
		IsActive: true,
	}
	err := testQueries.UpdateUserIsActive(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateUserPassword(t *testing.T) {
	user1 := createRandomUser(t)
	arg := UpdateUserPasswordParams{
		ID:       user1.ID,
		Password: "newpasswordgenerated",
	}
	err := testQueries.UpdateUserPassword(context.Background(), arg)
	require.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUsers(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}
}
