package db

import (
	"context"
	//"log"
	"testing"

	"github.com/MeganViga/Bank/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)
func createRandomAccount(t *testing.T)Account{
	arg := CreateAccountParams{
		Owner: utils.RandomOwner(),
		Balance: int64(utils.RandomBalance()),
		Currency: utils.RandomCurrency(),
	
	}
	account, err := testQueries.CreateAccount(context.Background(),arg)
	require.NoError(t,err)
	require.NotEmpty(t,account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}
func TestCreateAccount(t *testing.T){
	createRandomAccount(t)
}

func TestGetAccountByID(t *testing.T){
	account := createRandomAccount(t)
	account2 , err := testQueries.GetAccountByID(context.Background(),account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Balance, account2.Balance)
	//require.WithinDuration(t, account.CreatedAt, account2.CreatedAt,time.Second)
	require.Equal(t, account.CreatedAt, account2.CreatedAt)
}

func TestUpdateAccount(t *testing.T){
	account:= createRandomAccount(t)
	debit_value := 10
	arg := UpdateAccountParams{
		ID: account.ID,
		Balance: account.Balance - int64(debit_value),
	}
	account2, err := testQueries.UpdateAccount(context.Background(),arg)
	//log.Println(account,account2)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account.ID, account2.ID)
	require.Equal(t, account.Owner, account2.Owner)
	require.Equal(t, account.Currency, account2.Currency)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account.Balance - int64(debit_value), account2.Balance)
	require.Equal(t, account.CreatedAt, account2.CreatedAt)

}

func TestDeleteAccount(t *testing.T){
	account:= createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(),account.ID)
	require.NoError(t, err)
	account2, err := testQueries.GetAccountByID(context.Background(), account.ID)
	require.Error(t, err)
	require.EqualError(t,err,pgx.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T){
	for i:=0;i<10;i++{
		createRandomAccount(t)
	}
	arg := ListAccountsParams{
		Limit: 5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)
	for _, account := range accounts{
		require.NotEmpty(t, account)
	}
}