## mysql-server

This is an *experimental* implementation of a [`MergeStat`](https://github.com/mergestat/mergestat) backend for [`go-mysql-server`](https://github.com/dolthub/go-mysql-server) from DoltHub.
The `go-mysql-server` is a "frontend" SQL engine based on a MySQL syntax and wire protocol implementation.
It allows for pluggable "backends" as database and table providers.

### Usage

The easiest way to get started (for now) is probably by building and running a Docker container locally.
`MergeStat` has some build dependencies (such as [`libgit2`](https://libgit2.org/)), which must be available on your system when compiling (see the `Makefile` for details).

For now, running `docker build . -t mergestat/mysql-server` will produce a docker image you can run like so:

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
