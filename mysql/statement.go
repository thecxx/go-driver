package mysql

import (
	"database/sql"
)

type Statement struct {
	*sql.Stmt
}
