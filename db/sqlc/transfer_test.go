package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/terenjit/simplebank/util"
)

func createRandomTranfers(t *testing.T, account1 Account, account2 Account) Transfer {
	arg := CreateTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	trans, err := testQueries.CreateTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, trans)

	require.Equal(t, arg.FromAccountID, trans.FromAccountID)
	require.Equal(t, arg.ToAccountID, trans.ToAccountID)
	require.Equal(t, arg.Amount, trans.Amount)

	require.NotZero(t, trans.ID)
	require.NotZero(t, trans.CreatedAt)

	return trans
}

func TestCreateTranfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTranfers(t, account1, account2)
}

func TestGetTranfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	trans := createRandomTranfers(t, account1, account2)
	trans2, err := testQueries.GetTransfers(context.Background(), trans.ID)
	require.NoError(t, err)
	require.NotEmpty(t, trans2)

	require.Equal(t, trans.ID, trans2.ID)
	require.Equal(t, trans.FromAccountID, trans2.FromAccountID)
	require.Equal(t, trans.ToAccountID, trans2.ToAccountID)
	require.Equal(t, trans.Amount, trans2.Amount)
	require.WithinDuration(t, trans.CreatedAt, trans2.CreatedAt, time.Second)
}

func TestListTranfers(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 5; i++ {
		createRandomTranfers(t, account1, account2)
		createRandomTranfers(t, account2, account1)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account1.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account1.ID)
	}
}
