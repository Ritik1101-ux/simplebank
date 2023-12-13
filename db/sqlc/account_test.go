package db

import (
	"context"
	"testing"
	"github.com/Ritik1101-ux/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func CreateAccount(t *testing.T) Account {

	user:=CreateRandomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	return account
}
func TestCreateAccount(t *testing.T) {
	CreateAccount(t)
}

func TestGetAccount(t *testing.T) {
	arg := CreateAccount(t)

	account2, err := testQueries.GetAccount(context.Background(), arg.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, arg.Owner, account2.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, arg.Currency, account2.Currency)

}

func TestUpdateAccount(t *testing.T) {
	account1 := CreateAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account1.ID, account2.ID)
}

func TestDeleteAccount(t *testing.T) {
	account1 := CreateAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.ID)

	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.Empty(t, account2)

}

func TestListAccounts(t *testing.T) {

	for i := 0; i < 10; i++ {
		CreateAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
