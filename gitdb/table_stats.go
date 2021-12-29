package gitdb

import (
	gosql "database/sql"
	"io"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/jmoiron/sqlx"
)

type StatsTable struct {
	name string
	db   *Database
}

func NewStatsTable(name string, db *Database) *StatsTable {
	return &StatsTable{
		name: name,
		db:   db,
	}
}

func (t *StatsTable) Name() string   { return t.name }
func (t *StatsTable) String() string { return t.name }

func (t *StatsTable) Schema() sql.Schema {
	return []*sql.Column{
		{Name: "commit_hash", Type: sql.Text, Source: t.name},
		{Name: "file_path", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "additions", Type: sql.Int64, Nullable: true, Source: t.name},
		{Name: "deletions", Type: sql.Int64, Nullable: true, Source: t.name},
	}
}

func (t *StatsTable) Partitions(ctx *sql.Context) (sql.PartitionIter, error) {
	return &singlePartitionIter{}, nil
}

func (t *StatsTable) PartitionRows(ctx *sql.Context, _ sql.Partition) (sql.RowIter, error) {
	return newStatsIter(t)
}

type statsIter struct {
	*StatsTable
	rows *sqlx.Rows
}

type stat struct {
	CommitHash gosql.NullString `db:"hash"`
	FilePath   gosql.NullString `db:"file_path"`
	Additions  int              `db:"additions"`
	Deletions  int              `db:"deletions"`
}

const statsQuery = `
SELECT
	commits.hash, file_path, additions, deletions
FROM commits(?), stats(?, commits.hash)
`

func newStatsIter(t *StatsTable) (*statsIter, error) {
	iter := &statsIter{StatsTable: t}
	if rows, err := iter.db.mergestat.Queryx(statsQuery, t.db.name, t.db.name); err != nil {
		return nil, err
	} else {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		iter.rows = rows
	}
	return iter, nil
}

func (iter *statsIter) Next() (sql.Row, error) {
	s := stat{}
	if iter.rows.Next() {
		if err := iter.rows.Err(); err != nil {
			return nil, err
		}
		if err := iter.rows.StructScan(&s); err != nil {
			return nil, err
		}
	} else {
		return nil, io.EOF
	}

	out := make([]interface{}, 0)

	if s.CommitHash.Valid {
		out = append(out, s.CommitHash.String)
	} else {
		out = append(out, nil)
	}
	if s.FilePath.Valid {
		out = append(out, s.FilePath.String)
	} else {
		out = append(out, nil)
	}
	out = append(out, s.Additions)
	out = append(out, s.Deletions)

	return out, nil
}

func (iter *statsIter) Close(*sql.Context) error {
	return iter.rows.Close()
}
