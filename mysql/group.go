package mysql

import (
	"fmt"
	"sync/atomic"
)

type Backend struct {
	server *Database
	weight int64
}

type Group struct {
	backends []Backend
	//
	tw int64
}

// New Group
func NewGroup() *Group {
	return new(Group)
}

// Add
func (g *Group) add(db *Database, weight int64) {
	// Total weight
	atomic.AddInt64(&g.tw, int64(weight))
	// Add
	g.backends = append(g.backends, Backend{db, weight})
}

// Count
func (g *Group) count() int {
	return len(g.backends)
}

// Schedule
func (g *Group) schedule() *Database {
	return g.backends[0].server
}

// Check if it is empty
func (g *Group) isEmpty() bool {
	return len(g.backends) == 0
}

// Check if it is available
func (g *Group) ping() error {
	if g.isEmpty() {
		return fmt.Errorf("No backend usable")
	}
	var err error
	for _, backend := range g.backends {
		err = backend.server.Ping()
		if err == nil {
			break
		}
	}
	return err
}
