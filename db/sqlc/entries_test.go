package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/jackc/pgx/v5"
)


func createRandomEntry(t *testing.T)Entry{
	account := createRandomAccount(t)
	entry_value := int64(10)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount: int64(entry_value),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, account.ID, entry.AccountID)
	require.Equal(t, entry_value, entry.Amount)
	require.NotZero(t,entry.CreatedAt)
	return entry
}
func TestCreateEntry(t *testing.T){
	createRandomEntry(t)
}

func TestGetEntryByID(t *testing.T){
	entry := createRandomEntry(t)
	entry2, err := testQueries.GetEntryByID(context.Background(), entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)
	require.Equal(t, entry.AccountID,entry2.AccountID)
	require.Equal(t, entry.Amount, entry2.Amount)
	require.Equal(t, entry.CreatedAt, entry2.CreatedAt)
	require.Equal(t, entry.ID, entry2.ID)
}

func TestDeleteEntry(t *testing.T){
	entry:= createRandomEntry(t)
	err := testQueries.DeleteEntry(context.Background(),entry.ID)
	require.NoError(t, err)
	entry2, err := testQueries.GetEntryByID(context.Background(), entry.ID)
	require.Error(t, err)
	require.EqualError(t,err,pgx.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListEntries(t *testing.T){
	for i:=0;i<10;i++{
		createRandomEntry(t)
	}
	arg := ListEntriesParams{
		Limit: 5,
		Offset: 5,
	}
	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 5)
	for _, entry := range entries{
		require.NotEmpty(t, entry)
	}
}