package mysql

import (
	"testing"
)

func TestQueryFromWriter(t *testing.T) {
	wr := NewDefaultConfig()
	wr.Addr = "127.0.0.1:3306"
	wr.DBName = "test"
	wr.User = "root"
	wr.Passwd = "123456"

	wr.AllowNativePasswords = true

	c := NewCluster()

	// writer
	err := c.AddBackend(ClusterWriter, wr, DefaultWeight, false)
	if err != nil {
		t.Errorf("AddBackend failed: %s\n", err.Error())
		return
	}

	// reader
	wr.Addr = "127.0.0.1:3307"

	err = c.AddBackend(ClusterReader, wr, DefaultWeight, false)
	if err != nil {
		t.Errorf("AddBackend failed: %s\n", err.Error())
		return
	}

	var (
		id     int
		field1 string
		field2 string
	)

	num, err := c.Query(func(index int) []interface{} {
		return []interface{}{&id, &field1, &field2}
	}, "/*w*/SELECT `id`, `field1`, `field2` FROM `table1` LIMIT 1")
	if err != nil {
		t.Errorf("Query failed: %s\n", err.Error())
		return
	}

	t.Logf("OK\n num=%d, id=%d, field1=%s, field2=%s\n", num, id, field1, field2)
}
