package db

import (
	"context"
	"fmt"
	"log"

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
var txKey = struct{}{}
func (store *Store)TransferTx(ctx context.Context, arg TransferTxParams)(TransferTxResult, error){
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		log.Println(txName,"Create Transfer record")
		result.TransferRecord , err = q.CreateTransfer(ctx,CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}
		log.Println(txName,"Create Entry1")
		result.FromAccountEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil{
			return err
		}
		log.Println(txName,"Create Entry2")
		result.ToAccountEntry, err = q.CreateEntry(ctx,CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil{
			return err
		}

		//TODO: Update account balance
	//Get Account --> Then update it
		// log.Println(txName,"Get Account1 ")
		// account1, err := q.GetAccountByIDForUpdate(context.Background(), arg.FromAccountID)
		// if err != nil{
		// 	return err
		// }
		// 
		if arg.FromAccountID < arg.ToAccountID{
			//log.Println(txName,"Update Account1 balance")
		// result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		// 	ID: arg.FromAccountID,
		// 	Amount: -arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }

		// //log.Println(txName,"Update Account2 balance")
		// result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
		// 	ID: arg.ToAccountID,
		// 	Amount: arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		result.FromAccount , result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID,-arg.Amount,arg.ToAccountID, arg.Amount)
		if err != nil{
			return err
		}

		}else{
			// result.ToAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			// 	ID: arg.ToAccountID,
			// 	Amount: arg.Amount,
			// })
			// if err != nil{
			// 	return err
			// }


			// result.FromAccount, err = q.AddAccountBalance(context.Background(), AddAccountBalanceParams{
			// 	ID: arg.FromAccountID,
			// 	Amount: -arg.Amount,
			// })
			result.ToAccount , result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID,arg.Amount,arg.FromAccountID, -arg.Amount)
			if err != nil{
				return err
			}

		}
		
		// log.Println(txName,"Get Account2")
		// account2, err := q.GetAccountByIDForUpdate(context.Background(), arg.ToAccountID)
		// if err != nil{
		// 	return err
		// }
		// log.Println(txName,"Update Account2 balance")
		// result.ToAccount, err = q.UpdateAccount(context.Background(), UpdateAccountParams{
		// 	ID: arg.ToAccountID,
		// 	Balance: account2.Balance + arg.Amount,
		// })
		// if err != nil{
		// 	return err
		// }
		
		
		
		return nil
	})
	

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
)(account1 Account, account2 Account, err error){
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil{
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	return
}