package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)
type PoolInterface interface{
	DBTX
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

//Store provide all functions to run db queries
type Store struct{
	*Queries	
	db PoolInterface
}


func NewStore(db PoolInterface)*Store{
	return &Store{
		Queries: New(db),
		db:db,
	
	}
}

func (store *Store)execTx(ctx context.Context, fn func(*Queries)error)error{

	//Starting the transaction
	
	tx, err := store.db.BeginTx(ctx,pgx.TxOptions{})
	if err != nil{
		return err
	}
	//Converting transactions to Queries struct, 
	//here tx(which  is pgx.TX interface) is matching 
	//with DBTX interface since function from DBTX interface are there in pgx.TX interface)
	//that's why passing tx to New func is not throwing any error
	q := New(tx)
	err = fn(q) //--> this is the func where actual transaction operations will be done
	if err != nil{
		if rbErr := tx.Rollback(ctx); rbErr != nil{
			return fmt.Errorf("tx err: %v rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
type TransferTxParams struct{
	FromAccountID int64
	ToAccountID int64
	Amount int64
}
type TransferTxResult struct{
	FromAccount Account
	ToAccount Account
	FromAccountEntry Entry
	ToAccountEntry Entry
	TransferRecord Transfer
}

func (store *Store)TransferTx(ctx context.Context, arg TransferTxParams)(TransferTxResult, error){
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.TransferRecord , err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		result.FromAccountEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}

		result.ToAccountEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		
		return nil
	})
	//TODO: Update account balance

	return result, err
}