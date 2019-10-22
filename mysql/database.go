package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	driver "github.com/go-sql-driver/mysql"
)

const (
	DefaultMaxOpenConns = 50
	DefaultMaxIdleConns = 10
	DefaultMaxLifetime  = 30 * time.Second
	DefaultTimeout      = 2 * time.Second
	DefaultReadTimeout  = 0 * time.Second
	DefaultWriteTimeout = 0 * time.Second
)

type Config struct {
	*driver.Config
	// Extention
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// New a default config
func NewDefaultConfig() *Config {
	c := new(Config)
	c.Config = new(driver.Config)

	c.Net = "tcp"
	c.Timeout = DefaultTimeout
	c.ReadTimeout = DefaultReadTimeout
	c.WriteTimeout = DefaultWriteTimeout
	c.MaxOpenConns = DefaultMaxOpenConns
	c.MaxIdleConns = DefaultMaxIdleConns
	c.MaxLifetime = DefaultMaxLifetime

	return c
}

// Generate an unique tag
func (c *Config) Tag() string {
	return fmt.Sprintf("%s://%s/%s", c.Net, c.Addr, c.DBName)
}

type Database struct {
	// Attributes
	id string
	db *sql.DB
}

// New Database
func NewDatabase(addr, dbname, user, passwd string, options ...DatabaseOption) (*Database, error) {
	// Default
	config := NewDefaultConfig()
	// Auth
	config.Addr = addr
	config.DBName = dbname
	config.User = user
	config.Passwd = passwd
	// Apply options
	if len(options) > 0 {
		for _, handler := range options {
			handler(config)
		}
	}
	return NewDatabaseWithConfig(config)
}

// New database with config
func NewDatabaseWithConfig(config *Config) (*Database, error) {
	// Convert to DSN
	dsn := config.FormatDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Setting
	db.SetConnMaxLifetime(config.MaxLifetime)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return &Database{config.Tag(), db}, nil
}

// Stats
func (d *Database) Stats() Statistics {
	return Statistics{d.db.Stats()}
}

// Ping
func (d *Database) Ping() error {
	return d.db.PingContext(context.Background())
}

// Ping
func (d *Database) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

// Close
func (d *Database) Close() error {
	return d.db.Close()
}

// Check if it's in a transaction
func (d *Database) InTransaction() bool {
	return false
}

// Begin a new transaction
func (d *Database) BeginTransaction(ctx context.Context) (*Transaction, *TransactionCloser, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	return &Transaction{tx}, &TransactionCloser{tx}, nil
}

// Prepare
func (d *Database) Prepare(query string) (stmt *Statement, err error) {
	return d.PrepareContext(context.Background(), query)
}

// Prepare
func (d *Database) PrepareContext(ctx context.Context, query string) (stmt *Statement, err error) {
	s, err := d.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return &Statement{s}, nil
}

// Exec
func (d *Database) Exec(query string, args ...interface{}) (num int64, id int64, err error) {
	return d.ExecContext(context.Background(), query, args...)
}

// Exec
func (d *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (num int64, id int64, err error) {
	result, err := d.db.ExecContext(ctx, query, args...)
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
func (d *Database) Query(alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error) {
	return d.QueryContext(context.Background(), alloc, query, args...)
}

// Query
func (d *Database) QueryContext(ctx context.Context, alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error) {
	rows, err := d.db.QueryContext(ctx, query, args...)
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
