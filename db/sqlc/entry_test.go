package db

import (
	"context"
	"testing"

	"github.com/Ritik1101-ux/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func CreateEntry(t *testing.T,account Account) Entry {

	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:   utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.Empty(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
            
	return entry

}

func TestCreateEntry(t *testing.T) {
	account:=CreateAccount(t)
	CreateEntry(t,account)
}

func TestGetEntry(t *testing.T){
	account:=CreateAccount(t)

	entry:=CreateEntry(t,account)

	entry1,err:=testQueries.GetEntry(context.Background(),entry.ID)

	require.Empty(t,err)
	require.NotEmpty(t,entry1)
	require.Equal(t,entry.ID,entry1.ID)
	require.Equal(t,entry.Amount,entry1.Amount)
	require.Equal(t,entry.AccountID,entry1.AccountID)
}

func TestGetEntriesList(t *testing.T){
	account:=CreateAccount(t)
	for i:=0; i<10; i++{
        CreateEntry(t,account)
	}

	arg:=ListEntriesParams {
		AccountID:account.ID,
        Limit: 5,
		Offset: 5,
	}

	entries,err:=testQueries.ListEntries(context.Background(),arg)

	require.Empty(t,err)
	require.Len(t,entries,5)

	for _,entry:=range entries{
		require.NotEmpty(t,entry)
		require.Equal(t,entry.AccountID,account.ID)
	}
}
