## mysql-server

This is an *experimental* implementation of a [`MergeStat`](https://github.com/mergestat/mergestat) backend for [`go-mysql-server`](https://github.com/dolthub/go-mysql-server) from DoltHub.
The `go-mysql-server` is a "frontend" SQL engine based on a MySQL syntax and wire protocol implementation.
It allows for pluggable "backends" as database and table providers.

### Database

When you connect with a MySQL client, the **database name** you specify will be a reference to a git repository.
Currently, this can either be an HTTP(s) URL (`https://github.com/mergestat/mergestat`) or a path to a repository on disk.

For instance, the following will list all commits from [`mergestat/mergestat`](https://github.com/mergestat/mergestat).
The repo is cloned to a temporary directory in the container before the query is executed.

```
mysql --host=127.0.0.1 --port=3306 "https://github.com/mergestat/mergestat" -u root -proot -e "select * from commits"
```

#### Tables

- Commits
- Refs
- Files

### Usage

The easiest way to get started (for now) is probably by building and running a Docker container locally.
`MergeStat` has some build dependencies (such as [`libgit2`](https://libgit2.org/)), which must be available on your system when compiling (see the `Makefile` for details).

You can use the included `docker-compose.yaml` file in this repository by running `docker compose up`.
This will start a MySQL server on port `3306`.

You may also run `docker build . -t mergestat/mysql-server` to produce a docker image you can run like so:

```
docker run -p 3306:3306 -v ${PWD}:/repo mergestat/mysql-server
```

Note the `-v` flag, in the above we're mounting the current directory (assumed to be a git repo) into `/repo` in the container.

The repo path is specified by the **database name** used when connecting to the MySQL server.
So, in the above, you would connect like so:

```
mysql --host=127.0.0.1 --port=3306 "/repo" -u root -proot -e "select hash from commits"
```

Note the password defaults to `root`.
The user/password can be set by supplying the `MYSQL_USER` and `MYSQL_PWD` env vars.

#### TODO
- Publish pre-built binaries
- Publish a Docker image
