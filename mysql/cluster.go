package mysql

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

const (
	// Group names
	GroupWriter = "w"
	GroupReader = "r"
	// Default value
	DefaultWeight int64 = 1
)

var (
	//
	commands = map[string]string{
		"SELECT": GroupReader,
		"INSERT": GroupWriter,
		"UPDATE": GroupWriter,
		"DELETE": GroupWriter,
	}
)

type Cluster struct {
	query  *regexp.Regexp
	groups map[string]*Group
}

// New cluster
func NewCluster() *Cluster {
	return NewClusterWithGroups(nil, nil)
}

// New cluster with groups
func NewClusterWithGroups(w, r *Group) *Cluster {
	c := new(Cluster)
	c.query = regexp.MustCompile(`^(?:/\*([a-z]+)\*/)?([A-Za-z]+)\s+.+$`)
	c.groups = make(map[string]*Group)
	// writer/reader groups
	if w == nil {
		w = NewGroup()
	}
	if r == nil {
		r = NewGroup()
	}
	c.groups[GroupWriter] = w
	c.groups[GroupReader] = r

	return c
}

// New cluster with configs
func NewClusterWithConfigs(w, r *Config) (*Cluster, error) {
	wg := NewGroup()
	rg := NewGroup()
	// writer
	if w != nil {
		db, err := NewDatabaseWithConfig(w)
		if err != nil {
			return nil, err
		}
		wg.add(db, DefaultWeight)
	}
	// reader
	if w != nil {
		db, err := NewDatabaseWithConfig(r)
		if err != nil {
			return nil, err
		}
		rg.add(db, DefaultWeight)
	}
	return NewClusterWithGroups(wg, rg), nil
}

// Add a new backend
func (c *Cluster) AddBackend(group string, config *Config, weight int64, healthCheck bool) error {
	g, ok := c.groups[group]
	if !ok {
		return fmt.Errorf("No group usable")
	}
	db, err := NewDatabaseWithConfig(config)
	if err != nil {
		return err
	}
	// Health check
	if healthCheck {
		err = db.Ping()
		if err != nil {
			return err
		}
	}
	// Add
	g.add(db, weight)

	return nil
}

// Ping
func (c *Cluster) Ping(group string) error {
	g, ok := c.groups[group]
	if !ok {
		return fmt.Errorf("No group usable")
	}
	return g.ping()
}

// InTransaction
func (c *Cluster) InTransaction() bool {
	return false
}

// BeginTransaction
func (c *Cluster) BeginTransaction(ctx context.Context) (*Transaction, *TransactionCloser, error) {
	g := c.groups[GroupWriter]
	if g.isEmpty() {
		return nil, nil, fmt.Errorf("No backend usable")
	}
	return g.schedule().BeginTransaction(ctx)
}

// Prepare
// In the case of multiple backends, it will be bound to the same backend
func (c *Cluster) Prepare(query string) (stmt *Statement, err error) {
	return c.PrepareContext(context.Background(), query)
}

// PrepareContext
// In the case of multiple backends, it will be bound to the same backend
func (c *Cluster) PrepareContext(ctx context.Context, query string) (stmt *Statement, err error) {
	g, err := c.parseGroup(query, false)
	if err != nil {
		return nil, err
	}
	if g.isEmpty() {
		return nil, fmt.Errorf("No backend usable")
	}
	return g.schedule().PrepareContext(ctx, query)
}

// Exec
func (c *Cluster) Exec(query string, args ...interface{}) (num int64, id int64, err error) {
	return c.ExecContext(context.Background(), query, args...)
}

// ExecContext
func (c *Cluster) ExecContext(ctx context.Context, query string, args ...interface{}) (num int64, id int64, err error) {
	g, err := c.parseGroup(query, false)
	if err != nil {
		return 0, 0, err
	}
	if g.isEmpty() {
		return 0, 0, fmt.Errorf("No backend usable")
	}
	return g.schedule().ExecContext(ctx, query, args...)
}

// Query
func (c *Cluster) Query(alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error) {
	return c.QueryContext(context.Background(), alloc, query, args...)
}

// QueryContext
func (c *Cluster) QueryContext(ctx context.Context, alloc RowAllocFunc, query string, args ...interface{}) (num int64, err error) {
	g, err := c.parseGroup(query, true)
	if err != nil {
		return 0, err
	}
	if g.isEmpty() {
		return 0, fmt.Errorf("No backend usable")
	}
	return g.schedule().QueryContext(ctx, alloc, query, args...)
}

// Parse group for query statement
func (c *Cluster) parseGroup(query string, useGroupFlag bool) (*Group, error) {
	submatch := c.query.FindStringSubmatch(query)
	if len(submatch) != 3 {
		return nil, fmt.Errorf("Invalid query statment")
	}
	group := submatch[1]
	command := submatch[2]
	// Check command
	if command == "" {
		return nil, fmt.Errorf("Command cannot be empty")
	}
	command = strings.ToUpper(command)
	// Check group
	// e. '/*group*/SELECT * FROM `xxx`'
	if useGroupFlag && group != "" {
		g, ok := c.groups[group]
		if !ok {
			return nil, fmt.Errorf("No group usable")
		}
		return g, nil
	}
	group, ok := commands[command]
	if !ok {
		return nil, fmt.Errorf("Command not supported")
	}
	return c.groups[group], nil
}
