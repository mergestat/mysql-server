package gitdb

import (
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	name      string
	mergestat *sqlx.DB
}

func NewDatabase(name string, mergestat *sqlx.DB) *Database {
	return &Database{name: name, mergestat: mergestat}
}

func (db *Database) Name() string {
	return db.name
}

func (db *Database) GetTableInsensitive(ctx *sql.Context, tblName string) (sql.Table, bool, error) {
	switch strings.ToLower(tblName) {
	case "commits":
		return NewCommitsTable(tblName, db), true, nil
	case "refs":
		return NewRefsTable(tblName, db), true, nil
	case "files":
		return NewFilesTable(tblName, db), true, nil
	case "stats":
		return NewStatsTable(tblName, db), true, nil
	default:
		return nil, false, nil
	}
}

func (db *Database) GetTableNames(ctx *sql.Context) ([]string, error) {
	return []string{"commits", "refs", "files", "stats"}, nil
}

func (db *Database) IsReadOnly() {}
