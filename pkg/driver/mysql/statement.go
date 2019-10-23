package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

type Statement struct {
	stmt *sql.Stmt
}

// Exec
func (s *Statement) Exec(args ...interface{}) (num int64, id int64, err error) {
	return s.ExecContext(context.Background(), args...)
}

// ExecContext
func (s *Statement) ExecContext(ctx context.Context, args ...interface{}) (num int64, id int64, err error) {
	result, err := s.stmt.ExecContext(ctx, args...)
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
func (s *Statement) Query(alloc RowAllocFunc, args ...interface{}) (num int64, err error) {
	return s.QueryContext(context.Background(), alloc, args...)
}

// QueryContext
func (s *Statement) QueryContext(ctx context.Context, alloc RowAllocFunc, args ...interface{}) (num int64, err error) {
	rows, err := s.stmt.QueryContext(ctx, args...)
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
