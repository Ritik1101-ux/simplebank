package db

import (
	"context"
	"testing"
	"time"

	"github.com/Ritik1101-ux/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomTransfer(t *testing.T, account1 Account, account2 Account) Transfer {

	arg := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.Empty(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)
	require.NotZero(t, transfer.ID)

	return transfer

}

func TestCreateTransfer(t *testing.T) {
	account1 := CreateAccount(t)
	account2 := CreateAccount(t)

	CreateRandomTransfer(t, account1, account2)
}

func TestGetTransfer(t *testing.T) {
	account1 := CreateAccount(t)
	account2 := CreateAccount(t)

	transfer := CreateRandomTransfer(t, account1, account2)

	transferG, err := testQueries.GetTransfer(context.Background(), transfer.ID)

	require.Empty(t, err)
	require.NotEmpty(t, transferG)

	require.Equal(t, transfer.FromAccountID, transferG.FromAccountID)
	require.Equal(t, transfer.ToAccountID, transferG.ToAccountID)
	require.Equal(t, transfer.Amount, transferG.Amount)

	require.Equal(t, transfer.ID, transferG.ID)
	require.WithinDuration(t, transfer.CreatedAt, transferG.CreatedAt, time.Second)

}

func TestListTransfer(t *testing.T) {

	account1 := CreateAccount(t)
	account2 := CreateAccount(t)

	for i := 0; i < 10; i++ {
		CreateRandomTransfer(t, account1, account2)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         5,
		Offset:        5,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)

	require.Empty(t, err)
	require.Len(t, transfers, 5)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == account1.ID || transfer.ToAccountID == account2.ID)
	}

}
