package gitdb

import (
	gosql "database/sql"
	"io"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/jmoiron/sqlx"
)

type FilesTable struct {
	name string
	db   *Database
}

func NewFilesTable(name string, db *Database) *FilesTable {
	return &FilesTable{
		name: name,
		db:   db,
	}
}

func (t *FilesTable) Name() string   { return t.name }
func (t *FilesTable) String() string { return t.name }

func (t *FilesTable) Schema() sql.Schema {
	return []*sql.Column{
		{Name: "path", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "executable", Type: sql.Boolean, Nullable: true, Source: t.name},
		{Name: "contents", Type: sql.LongBlob, Nullable: true, Source: t.name},
	}
}

func (t *FilesTable) Partitions(ctx *sql.Context) (sql.PartitionIter, error) {
	return &singlePartitionIter{}, nil
}

func (t *FilesTable) PartitionRows(ctx *sql.Context, _ sql.Partition) (sql.RowIter, error) {
	return newFilesIter(t)
}

type filesIter struct {
	*FilesTable
	rows *sqlx.Rows
}

type file struct {
	Path       gosql.NullString `db:"path"`
	Executable gosql.NullBool   `db:"executable"`
	Contents   gosql.NullString `db:"contents"`
}

const filesQuery = `
SELECT
	path, executable, contents
FROM files(?)
`

func newFilesIter(t *FilesTable) (*filesIter, error) {
	iter := &filesIter{FilesTable: t}
	if rows, err := iter.db.mergestat.Queryx(filesQuery, t.db.name); err != nil {
		return nil, err
	} else {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		iter.rows = rows
	}
	return iter, nil
}

func (iter *filesIter) Next() (sql.Row, error) {
	s := file{}
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

	if s.Path.Valid {
		out = append(out, s.Path.String)
	} else {
		out = append(out, nil)
	}
	if s.Executable.Valid {
		out = append(out, s.Executable.Bool)
	} else {
		out = append(out, nil)
	}
	if s.Contents.Valid {
		out = append(out, s.Contents.String)
	} else {
		out = append(out, nil)
	}

	return out, nil
}

func (iter *filesIter) Close(*sql.Context) error {
	return iter.rows.Close()
}
