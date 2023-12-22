package db

import (
	"context"
	"github.com/Ritik1101-ux/simplebank/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func CreateAccount(t *testing.T) Account {

	user := CreateRandomUser(t)
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

	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = CreateAccount(t)
	}

	arg := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
