package mysql

import (
	"context"
)

type RowAllocFunc func(index int) (pointers []interface{})

// Implement by cluster/database/transaction
type Client interface {
	// Transaction
	InTransaction() bool
	BeginTransaction(context.Context) (*Transaction, *TransactionCloser, error)
	// Operations
	Prepare(query string) (stmt *Statement, err error)
	PrepareContext(ctx context.Context, query string) (stmt *Statement, err error)
	Exec(query string, args ...interface{}) (num int64, id int64, err error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (num int64, id int64, err error)
	Query(alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error)
	QueryContext(ctx context.Context, alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error)
}
