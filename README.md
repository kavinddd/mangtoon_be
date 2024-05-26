## How to set up

### Install
1. [GOLANG](https://go.dev/doc/install)
   -  v.1.2.1 + is required to install SQLC using GO
2. [Docker (or podman)](https://www.docker.com/) with [postgres image](https://hub.docker.com/_/postgres)
3. [SQLC](https://github.com/sqlc-dev/sqlc)
   - is used to generate go code to interact with postgres database from native sql
4. [Goose](https://github.com/pressly/goose)
    - is used for migrating and versioning database
5. [gomock](https://github.com/golang/mock)
   - is used to generate go code for mocking database
6. [make]()
   - is used to run predefined script
   - need to install through chocolatey in windows os, most linuxs come with make already


### Then run prepared script in Makefile (edit it as you like if it doesn't work in your environment)
1. `make db-contatiner-create`
   - this creates a postgres container from image named postgres
   - if you want to just run the container then use `make db-container-run`
2. `make db-create`
   - this creates a mangtoon database inside the postgres container (from the 1.)
3. `make db-migrate-up`
   - this applies all *.sql in /internal/db/migrations/ to the postgres container (from the 1.)

## Development Practices
### DB Migrations
You want to change DB schema, do the following
1. run `make db-new-migrate`, this will generate a new SQL file in the /internal/db/migrations/20241203123052_new.sql
2. You write native sql (DDL) to the file, change the name from new to something more descriptive
3. run `make db-migrate-up` to apply the new SQL file to your postgres db container
### New queries (think of repository in MVC)
1. You write native sql (DML) inside /internal/db/sqlc/
2. run `make sqlc` to generate go code which will be appeared in /internal/db/db.go and *.sql.go
### Testing
1. always do unit testing in api using gomock
