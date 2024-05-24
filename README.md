## How to set up

### Install
1. [GOLANG](https://go.dev/doc/install)
   -  v.1.2.1 + is required to install SQLC using GO
2. [Docker (or podman)](https://www.docker.com/) with [postgres image](https://hub.docker.com/_/postgres)
3. [SQLC](https://github.com/sqlc-dev/sqlc)
   - is used to generate go code to interact with postgres database from native sql
4. [Goose](https://github.com/pressly/goose)
    - is used for migrating and versioning database
5. [make]()
   - is used to run predefined script
   - need to install through chocolatey in windows os, most linuxs come with make already


### Then run prepared script in Makefile (edit it as you like if it doesn't work in your environment)
1. `make db-contatiner-create`
2. `make db-create`
3. `make db-migrate-up`

