package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := CreateAccount(t)
	account2 := CreateAccount(t)

	//run a concurrent transfer transactions

	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)

		go func() { //Go Routine
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	//Check Results

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//Check Transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)

		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)

		require.NoError(t, err)

		//Check Entries

		FromEntry := result.FromEntry

		require.NotEmpty(t, FromEntry)
		require.Equal(t, account1.ID, FromEntry.AccountID)
		require.Equal(t, -amount, FromEntry.Amount)
		require.NotZero(t, FromEntry.ID)
		require.NotZero(t, FromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), FromEntry.ID)
		require.NoError(t, err)

		ToEntry := result.ToEntry

		require.NotEmpty(t, ToEntry)
		require.Equal(t, account2.ID, ToEntry.AccountID)
		require.Equal(t, amount, ToEntry.Amount)
		require.NotZero(t, ToEntry.ID)
		require.NotZero(t, ToEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), ToEntry.ID)
		require.NoError(t, err)

		// Check Account Balance

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">>txn: ", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) //1*amount,2*amount,3*amount ,.... n*amount

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	// Check the final updated balance

	updatedAccount1, err1 := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err1)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err2 := testQueries.GetAccount(context.Background(), account2.ID)

	require.NoError(t, err2)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedAccount2.Balance)

	fmt.Println(">> After:", updatedAccount1.Balance, updatedAccount2.Balance)

}

func TestTransferDeadlockTx(t *testing.T) {
	store := NewStore(testDb)

	account1 := CreateAccount(t)
	account2 := CreateAccount(t)

	//run a concurrent transfer transactions

	fmt.Println(">> Before:", account1.Balance, account2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {

		FromAccountId := account1.ID
		ToAccountId := account2.ID

		if i%2 == 1 {
			FromAccountId = account2.ID
			ToAccountId = account1.ID
		}
		go func() { //Go Routine
		
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: FromAccountId,
				ToAccountID:   ToAccountId,
				Amount:        amount,
			})
			errs <- err

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs

		require.NoError(t, err)

	}

	// Check the final updated balance

	updatedAccount1, err1 := testQueries.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err1)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err2 := testQueries.GetAccount(context.Background(), account2.ID)

	require.NoError(t, err2)
	require.NotEmpty(t, updatedAccount2)

	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)

	fmt.Println(">> After:", updatedAccount1.Balance, updatedAccount2.Balance)

}
