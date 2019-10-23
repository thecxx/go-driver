package mysql

import (
	"testing"
)

func TestConfigTag(t *testing.T) {
	c := NewDefaultConfig()
	c.Addr = "127.0.0.1:3306"
	c.DBName = "test"
	c.User = "root"
	c.Passwd = "123456"

	t.Logf("Tag: %s\n", c.Tag())
}
