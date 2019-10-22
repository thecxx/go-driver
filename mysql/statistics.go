package mysql

import (
	"database/sql"
)

type Statistics struct {
	sql.DBStats
}
