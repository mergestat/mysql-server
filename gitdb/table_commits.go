package gitdb

import (
	gosql "database/sql"
	"io"
	"time"

	"github.com/dolthub/go-mysql-server/sql"
	"github.com/jmoiron/sqlx"
)

type CommitsTable struct {
	name string
	db   *Database
}

func NewCommitsTable(name string, db *Database) *CommitsTable {
	return &CommitsTable{
		name: name,
		db:   db,
	}
}

func (t *CommitsTable) Name() string   { return t.name }
func (t *CommitsTable) String() string { return t.name }

func (t *CommitsTable) Schema() sql.Schema {
	return []*sql.Column{
		{Name: "hash", Type: sql.Text, Source: t.name},
		{Name: "message", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "author_name", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "author_email", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "author_when", Type: sql.Datetime, Nullable: true, Source: t.name},
		{Name: "committer_name", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "committer_email", Type: sql.Text, Nullable: true, Source: t.name},
		{Name: "committer_when", Type: sql.Datetime, Nullable: true, Source: t.name},
		{Name: "parents", Type: sql.Int64, Nullable: true, Source: t.name},
	}
}

func (t *CommitsTable) Partitions(ctx *sql.Context) (sql.PartitionIter, error) {
	return &singlePartitionIter{}, nil
}

func (t *CommitsTable) PartitionRows(ctx *sql.Context, _ sql.Partition) (sql.RowIter, error) {
	return newCommitsIter(t)
}

type commitsIter struct {
	*CommitsTable
	rows *sqlx.Rows
}

type commit struct {
	Hash           gosql.NullString `db:"hash"`
	Message        gosql.NullString `db:"message"`
	AuthorName     gosql.NullString `db:"author_name"`
	AuthorEmail    gosql.NullString `db:"author_email"`
	AuthorWhen     time.Time        `db:"author_when"`
	CommitterName  gosql.NullString `db:"committer_name"`
	CommitterEmail gosql.NullString `db:"committer_email"`
	CommitterWhen  time.Time        `db:"committer_when"`
	Parents        int              `db:"parents"`
}

const commitsQuery = `
SELECT
	hash, message, author_name, author_email, author_when, committer_name, committer_email, committer_when, parents
FROM commits(?)
`

func newCommitsIter(t *CommitsTable) (*commitsIter, error) {
	iter := &commitsIter{CommitsTable: t}
	if rows, err := iter.db.mergestat.Queryx(commitsQuery, t.db.name); err != nil {
		return nil, err
	} else {
		if err := rows.Err(); err != nil {
			return nil, err
		}
		iter.rows = rows
	}
	return iter, nil
}

func (iter *commitsIter) Next() (sql.Row, error) {
	s := commit{}
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

	if s.Hash.Valid {
		out = append(out, s.Hash.String)
	} else {
		out = append(out, nil)
	}
	if s.Message.Valid {
		out = append(out, s.Message.String)
	} else {
		out = append(out, nil)
	}
	if s.AuthorName.Valid {
		out = append(out, s.AuthorName.String)
	} else {
		out = append(out, nil)
	}
	if s.AuthorEmail.Valid {
		out = append(out, s.AuthorEmail.String)
	} else {
		out = append(out, nil)
	}
	out = append(out, s.AuthorWhen)
	if s.CommitterName.Valid {
		out = append(out, s.CommitterName.String)
	} else {
		out = append(out, nil)
	}
	if s.CommitterEmail.Valid {
		out = append(out, s.CommitterEmail.String)
	} else {
		out = append(out, nil)
	}
	out = append(out, s.CommitterWhen)
	out = append(out, s.Parents)

	return out, nil
}

func (iter *commitsIter) Close(*sql.Context) error {
	return iter.rows.Close()
}
