package db

import (
	"context"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	// run n concurrent transactions
	n :=5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	log.Println(testStore)
	for i :=0;i < n;i++{
		go func() {
			result, err := testStore.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}

	//check errors and result from outside
	for i:=0; i< n;i++{
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//check transfer
		transfer := result.TransferRecord
		require.NotEmpty(t, transfer)
		require.Equal(t,account1.ID, transfer.FromAccountID)
		require.Equal(t,account2.ID, transfer.ToAccountID)
		require.Equal(t,amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)
		_, err = testStore.GetTransferByID(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromAccountEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t,account1.ID, fromEntry.AccountID)
		require.Equal(t,-amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = testStore.GetEntryByID(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToAccountEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t,account2.ID, toEntry.AccountID)
		require.Equal(t,amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		_, err = testStore.GetEntryByID(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO:check account's balance
		
	
	}
	
	
}