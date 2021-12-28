package gitdb

import (
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	"github.com/jmoiron/sqlx"
)

type Provider struct {
	mergestat         *sqlx.DB
	informationSchema sql.Database
}

func NewProvider(mergestat *sqlx.DB) *Provider {
	return &Provider{informationSchema: information_schema.NewInformationSchemaDatabase(), mergestat: mergestat}
}

func (p *Provider) Database(name string) (sql.Database, error) {
	if name == "information_schema" {
		return p.informationSchema, nil
	}
	return NewDatabase(name, p.mergestat), nil
}

func (p *Provider) HasDatabase(name string) bool {
	return true
}

func (p *Provider) AllDatabases() []sql.Database {
	return []sql.Database{p.informationSchema}
}
