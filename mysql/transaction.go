package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

type TransactionCloser struct {
	tx *sql.Tx
}

// Commit
func (t *TransactionCloser) Commit() error {
	return t.tx.Commit()
}

// Rollback
func (t *TransactionCloser) Rollback() error {
	return t.tx.Rollback()
}

type Transaction struct {
	// Attributes
	tx *sql.Tx
}

// Check if it's in a transaction
func (t *Transaction) InTransaction() bool {
	return true
}

// Begin a new transaction
func (t *Transaction) BeginTransaction(context.Context) (*Transaction, *TransactionCloser, error) {
	return nil, nil, fmt.Errorf("Already in transaction")
}

// Prepare
func (t *Transaction) Prepare(query string) (stmt *Statement, err error) {
	return t.PrepareContext(context.Background(), query)
}

// Prepare
func (t *Transaction) PrepareContext(ctx context.Context, query string) (stmt *Statement, err error) {
	s, err := t.tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Statement{s}, nil
}

// Exec
func (t *Transaction) Exec(query string, args ...interface{}) (num int64, id int64, err error) {
	return t.ExecContext(context.Background(), query, args...)
}

// Exec
func (t *Transaction) ExecContext(ctx context.Context, query string, args ...interface{}) (num int64, id int64, err error) {
	result, err := t.tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, 0, err
	}
	num, err = result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	id, err = result.LastInsertId()
	if err != nil {
		return 0, 0, err
	}
	return num, id, nil
}

// Query
func (t *Transaction) Query(alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error) {
	return t.QueryContext(context.Background(), alloc, query, args...)
}

// Query
func (t *Transaction) QueryContext(ctx context.Context, alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	idx := 0
	for rows.Next() {
		set := alloc(idx)
		if len(set) != len(columns) {
			return 0, fmt.Errorf("Inconsistent number of fields")
		}
		err = rows.Scan(set...)
		if err != nil {
			return 0, err
		}
		idx += 1
		num += 1
	}
	return num, nil
}
