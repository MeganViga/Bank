package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before transaction:", account1.Balance, account2.Balance)
	// run n concurrent transactions
	n := 3
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	//log.Println(testStore)
	for i :=0;i < n;i++{
		txName := fmt.Sprintf("tx: %d", i+1)
		ctx := context.WithValue(context.Background(),txKey,txName)
		go func() {
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID: account2.ID,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}

	//check errors and result from outside
	existed := make(map[int]bool)
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
		//check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)
		log.Println(">> transaction:", fromAccount.Balance, toAccount.Balance)
		//check account balance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1 % amount == 0)

		k := int(diff1/amount)
		require.True(t, k >= 1 && k <= n )
		require.NotContains(t, existed,k)
		existed[k]= true
	
	}

	//check final updated balances
	updatedAccount1, err := testQueries.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	updatedAccount2, err := testQueries.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	log.Println(">> after transactions:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance - int64(n) * amount,updatedAccount1.Balance)
	require.Equal(t, account2.Balance + int64(n) * amount,updatedAccount2.Balance)

	
	
}

func TestTransferTxDeadlock(t *testing.T){
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println(">> before transaction:", account1.Balance, account2.Balance)
	// run n concurrent transactions
	n := 10
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	//log.Println(testStore)
	for i :=0;i < n;i++{
		fromAccount := account1.ID
		toAccount := account2.ID
		if i %2 == 0{
			fromAccount = account2.ID
			toAccount = account1.ID
		}
		txName := fmt.Sprintf("tx: %d", i+1)
		ctx := context.WithValue(context.Background(),txKey,txName)
		go func() {
			result, err := testStore.TransferTx(ctx, TransferTxParams{
				FromAccountID: fromAccount,
				ToAccountID: toAccount,
				Amount: amount,
			})
			errs <- err
			results <- result
		}()
	}

	//check errors and result from outside
	// existed := make(map[int]bool)
	for i:=0; i< n;i++{
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		// //check transfer
		// transfer := result.TransferRecord
		// require.NotEmpty(t, transfer)
		// require.Equal(t,account1.ID, transfer.FromAccountID)
		// require.Equal(t,account2.ID, transfer.ToAccountID)
		// require.Equal(t,amount, transfer.Amount)
		// require.NotZero(t, transfer.ID)
		// require.NotZero(t, transfer.CreatedAt)
		// _, err = testStore.GetTransferByID(context.Background(), transfer.ID)
		// require.NoError(t, err)

		// check entries
		// fromEntry := result.FromAccountEntry
		// require.NotEmpty(t, fromEntry)
		// require.Equal(t,account1.ID, fromEntry.AccountID)
		// require.Equal(t,-amount, fromEntry.Amount)
		// require.NotZero(t, fromEntry.ID)
		// require.NotZero(t, fromEntry.CreatedAt)
		// _, err = testStore.GetEntryByID(context.Background(), fromEntry.ID)
		// require.NoError(t, err)

		// toEntry := result.ToAccountEntry
		// require.NotEmpty(t, toEntry)
		// require.Equal(t,account2.ID, toEntry.AccountID)
		// require.Equal(t,amount, toEntry.Amount)
		// require.NotZero(t, toEntry.ID)
		// require.NotZero(t, toEntry.CreatedAt)
		// _, err = testStore.GetEntryByID(context.Background(), toEntry.ID)
		// require.NoError(t, err)

		//TODO:check account's balance
		//check accounts
		// fromAccount := result.FromAccount
		// require.NotEmpty(t, fromAccount)
		// require.Equal(t, account1.ID, fromAccount.ID)

		// toAccount := result.ToAccount
		// require.NotEmpty(t, toAccount)
		// require.Equal(t, account2.ID, toAccount.ID)
		// log.Println(">> transaction:", fromAccount.Balance, toAccount.Balance)
		// //check account balance
		// diff1 := account1.Balance - fromAccount.Balance
		// diff2 := toAccount.Balance - account2.Balance
		// require.Equal(t, diff1, diff2)
		// require.True(t, diff1 > 0)
		// require.True(t, diff1 % amount == 0)

		// k := int(diff1/amount)
		// require.True(t, k >= 1 && k <= n )
		// require.NotContains(t, existed,k)
		// existed[k]= true
	
	}

	//check final updated balances
	updatedAccount1, err := testQueries.GetAccountByID(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)
	updatedAccount2, err := testQueries.GetAccountByID(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)
	log.Println(">> after transactions:", updatedAccount1.Balance, updatedAccount2.Balance)
	require.Equal(t, account1.Balance,updatedAccount1.Balance)
	require.Equal(t, account2.Balance,updatedAccount2.Balance)


	
}