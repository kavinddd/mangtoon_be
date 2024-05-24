
ifeq ($(shell command -v podman 2> /dev/null),)
  CONTAINER_RUNTIME := docker
else
  CONTAINER_RUNTIME := podman
endif

test:   
	@go test -v -cover ./...

db-create:
	${CONTAINER_RUNTIME} exec -it mangtoon_postgres createdb --username=root --owner=root mangtoon

db-drop:
	${CONTAINER_RUNTIME} exec -it mangtoon_postgres dropdb mangtoon

db-container-create:
	${CONTAINER_RUNTIME} run --name mangtoon_postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=dev -d postgres

db-container-start:
	${CONTAINER_RUNTIME} start mangtoon_postgres

db-status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=root dbname=mangtoon sslmode=disable host=localhost port=5432 password=dev" goose -dir=internal/db/migrations status

db-migrate-up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=root dbname=mangtoon sslmode=disable host=localhost port=5432 password=dev" goose -dir=internal/db/migrations up

db-migrate-down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=root dbname=mangtoon sslmode=disable host=localhost port=5432 password=dev" goose -dir=internal/db/migrations down

db-migrate-reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING="user=root dbname=mangtoon sslmode=disable host=localhost port=5432 password=dev" goose -dir=internal/db/migrations reset

sqlc: 
	sqlc generate


.PHONY: db-create db-drop db-status db-migrate-up db-migrate-down db-migrate-reset db-container-create db-container-run sqlc test
