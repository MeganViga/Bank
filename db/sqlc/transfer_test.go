package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/jackc/pgx/v5"
)


func createRandomTransfer(t *testing.T)Transfer{
	account := createRandomAccount(t)
	account2 :=  createRandomAccount(t)
	transfer_value := int64(10)
	arg := CreateTransferParams{
		FromAccountID: account.ID,
		ToAccountID: account2.ID,
		Amount: transfer_value,
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.NotZero(t, transfer.ID)
	require.Equal(t, account.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, transfer_value, transfer.Amount)
	require.NotZero(t,transfer.CreatedAt)
	return transfer
}
func TestCreateTransfer(t *testing.T){
	createRandomTransfer(t)
}

func TestGetTransferByID(t *testing.T){
	transfer := createRandomTransfer(t)
	transfer2, err := testQueries.GetTransferByID(context.Background(), transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)
	require.Equal(t, transfer.FromAccountID,transfer2.FromAccountID)
	require.Equal(t, transfer.ToAccountID,transfer2.ToAccountID)
	require.Equal(t, transfer.Amount, transfer2.Amount)
	require.Equal(t, transfer.CreatedAt, transfer2.CreatedAt)
	require.Equal(t, transfer.ID, transfer2.ID)
}

func TestDeleteTransfer(t *testing.T){
	transfer:= createRandomTransfer(t)
	err := testQueries.DeleteTransfer(context.Background(),transfer.ID)
	require.NoError(t, err)
	transfer2, err := testQueries.GetTransferByID(context.Background(), transfer.ID)
	require.Error(t, err)
	require.EqualError(t,err,pgx.ErrNoRows.Error())
	require.Empty(t, transfer2)
}

func TestListTransfer(t *testing.T){
	for i:=0;i<10;i++{
		createRandomTransfer(t)
	}
	arg := ListTransfersParams{
		Limit: 5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)
	for _, transfer := range transfers{
		require.NotEmpty(t, transfer)
	}
}