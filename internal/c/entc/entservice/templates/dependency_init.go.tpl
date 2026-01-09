// @file: internal/infra/dependency/init.go

package dependency

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/syralon/coconut/toolkit/sqlite3"
)