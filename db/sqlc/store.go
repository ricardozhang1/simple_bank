package db

import (
	"context"
	"database/sql"
	"fmt"
)
// Store provide all functions to execute db queries and translations
// Store 对象提供了所有数据库的操作的查询和事务方法
type Store interface {
	Querier
	TransferTx(context.Context, TransferTxParams) (TransferTxResult, error)
}

// SQLStore provide all functions to execute db queries and translations
// SQLStore 对象提供了所有数据库的操作的查询和事务方法
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore create a Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db: db,
		Queries: New(db),
	}
}

// execTx executes a function with database translation
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			// 返回两个错误的处理方法
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTxParams contains the input parameters of transfer translation.
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer translation.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`     // entries 表
	FromAccount Account  `json:"from_account"` // accounts 表
	ToAccount   Account  `json:"to_account"`   // accounts 表
	FromEntry   Entry    `json:"from_entry"`   // entries 表
	ToEntry     Entry    `json:"to_entry"`     // entries 表
}

var txKey = struct {}{}

// TransferTx perform a money transfer from one account the other.
// It create a transfer record, and account entries, and update accounts' balance with a single translation.
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		txName := ctx.Value(txKey)
		fmt.Println(txName, "create transfer")

		result.Transfer, err = q.CreateTransfers(ctx, CreateTransfersParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "create entry 2")
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "get account 1")
		// get account -> update its balance  need lock
		//account1, err := q.GetAccountForUpdate(ctx, arg.FromAccountID)
		//if err != nil {
		//	return err
		//}

		fmt.Println(txName, "update account 1")

		if arg.FromAccountID < arg.ToAccountID {
			//result.FromAccount, err = q.AddaAccountBalance(ctx, AddaAccountBalanceParams{
			//	ID:     arg.FromAccountID,
			//	Amount: -arg.Amount,
			//})
			//if err != nil {
			//	return err
			//}
			//
			////fmt.Println(txName, "get account 2")
			////account2, err := q.GetAccountForUpdate(ctx, arg.ToAccountID)
			////if err != nil {
			////	return err
			////}
			//
			//fmt.Println(txName, "update account 2")
			//result.ToAccount, err = q.AddaAccountBalance(ctx, AddaAccountBalanceParams{
			//	ID:     arg.ToAccountID,
			//	Amount: arg.Amount,
			//})
			//if err != nil {
			//	return err
			//}
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}
		return nil
	})
	return result, err
}

func addMoney(
	ctx context.Context, q *Queries, accountID1 int64, amount1 int64, accountID2 int64, amount2 int64,
	) (account1, account2 Account, err error) {
	account1, err = q.AddaAccountBalance(ctx, AddaAccountBalanceParams{
		ID: accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddaAccountBalance(ctx, AddaAccountBalanceParams{
		ID: accountID2,
		Amount: amount2,
	})
	return
}


