package gitdb

import (
	gosql "database/sql"
	"io"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/jmoiron/sqlx"
)

type RefsTable struct {
	name string
	db   *Database
}

func NewRefsTable(name string, db *Database) *RefsTable {
	return &RefsTable{
		name: name,
		db:   db,
	}
}

func (t *RefsTable) Name() string   { return t.name }
func (t *RefsTable) String() string { return t.name }

func (t *RefsTable) Schema() sql.Schema {
	return []*sql.Column{
		{Name: "name", Type: sql.Text, Source: t.name},
		{Name: "type", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "remote", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "full_name", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "hash", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "target", Type: sql.Text, Nullable: true, Source: t.name},
	}
}

func (t *RefsTable) Partitions(ctx *sql.Context) (sql.PartitionIter, error) {
	return &singlePartitionIter{}, nil
}

func (t *RefsTable) PartitionRows(ctx *sql.Context, _ sql.Partition) (sql.RowIter, error) {
	return newRefsIter(t)
}

type refsIter struct {
	*RefsTable
	rows *sqlx.Rows
}

type ref struct {
	Name     gosql.NullString `db:"name"`
	Type     gosql.NullString `db:"type"`
	Remote   gosql.NullString `db:"remote"`
	FullName gosql.NullString `db:"full_name"`
	Hash     gosql.NullString `db:"hash"`
	Target   gosql.NullString `db:"target"`
}

const refsQuery = `
SELECT
	name, type, remote, full_name, hash, target
FROM refs(?)
`

func newRefsIter(t *RefsTable) (*refsIter, error) {
	iter := &refsIter{RefsTable: t}
	if rows, err := iter.db.mergestat.Queryx(refsQuery, t.db.name); err != nil {
		return nil, err
	} else {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		iter.rows = rows
	}
	return iter, nil
}

func (iter *refsIter) Next() (sql.Row, error) {
	s := ref{}
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

	if s.Name.Valid {
		out = append(out, s.Name.String)
	} else {
		out = append(out, nil)
	}
	if s.Type.Valid {
		out = append(out, s.Type.String)
	} else {
		out = append(out, nil)
	}
	if s.Remote.Valid {
		out = append(out, s.Remote.String)
	} else {
		out = append(out, nil)
	}
	if s.FullName.Valid {
		out = append(out, s.FullName.String)
	} else {
		out = append(out, nil)
	}
	if s.Hash.Valid {
		out = append(out, s.Hash.String)
	} else {
		out = append(out, nil)
	}
	if s.Target.Valid {
		out = append(out, s.Target.String)
	} else {
		out = append(out, nil)
	}

	return out, nil
}

func (iter *refsIter) Close(*sql.Context) error {
	return iter.rows.Close()
}
