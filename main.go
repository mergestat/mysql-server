package main

import (
	"fmt"
	"log"
	"os"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/auth"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mergestat/mergestat/extensions"
	"github.com/mergestat/mergestat/extensions/options"
	"github.com/mergestat/mergestat/pkg/locator"
	_ "github.com/mergestat/mergestat/pkg/sqlite"
	"github.com/mergestat/mysql-server/gitdb"
	"go.riyazali.net/sqlite"
)

var (
	user     = "root"
	password = "root"
)

func init() {
	sqlite.Register(
		extensions.RegisterFn(
			options.WithRepoLocator(locator.CachedLocator(locator.MultiLocator())),
		),
	)

	if u := os.Getenv("MYSQL_USER"); u != "" {
		user = u
	}
	if p := os.Getenv("MYSQL_PWD"); p != "" {
		password = p
	}
}

func main() {
	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to initialize database connection: %v", err))
	}

	engine := sqle.NewDefault(gitdb.NewProvider(db))
	config := server.Config{
		Protocol: "tcp",
		Address:  "0.0.0.0:3306",
		Auth:     auth.NewNativeSingle(user, password, auth.ReadPerm),
	}

	s, err := server.NewDefaultServer(config, engine)
	if err != nil {
		log.Fatal(err)
	}

	s.Listener.AllowClearTextWithoutTLS = true

	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}
}
