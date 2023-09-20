package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nphattai/go-simple-bank/util"
	"github.com/stretchr/testify/require"
)

func randomUser(t *testing.T) User {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
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

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	randomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := randomUser(t)

	user2, err := testQueries.GetUser(context.Background(), user1.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.PasswordChangedAt, user2.PasswordChangedAt)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {
	user1 := randomUser(t)

	newFullName := util.RandomOwner()

	arg := UpdateUserParams{
		Username: user1.Username,
		FullName: sql.NullString{String: newFullName, Valid: true},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.Username, user2.Username)
	require.NotEqual(t, user1.FullName, user2.FullName)
	require.Equal(t, newFullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	user1 := randomUser(t)

	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	arg := UpdateUserParams{
		Username:       user1.Username,
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.NotEqual(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, hashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	user1 := randomUser(t)

	newEmail := util.RandomEmail()

	arg := UpdateUserParams{
		Username: user1.Username,
		Email:    sql.NullString{String: newEmail, Valid: true},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, arg.Username, user2.Username)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.NotEqual(t, user1.Email, user2.Email)
	require.Equal(t, newEmail, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}

func TestUpdateUserAllField(t *testing.T) {
	user1 := randomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)

	require.NoError(t, err)
	arg := UpdateUserParams{
		Username:       user1.Username,
		Email:          sql.NullString{String: newEmail, Valid: true},
		FullName:       sql.NullString{String: newFullName, Valid: true},
		HashedPassword: sql.NullString{String: hashedPassword, Valid: true},
	}

	user2, err := testQueries.UpdateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, newFullName, user2.FullName)
	require.Equal(t, hashedPassword, user2.HashedPassword)
	require.Equal(t, newEmail, user2.Email)

	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
}
