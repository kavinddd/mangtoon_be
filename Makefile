
ifeq ($(shell command -v podman 2> /dev/null),)
  CONTAINER_RUNTIME := docker
else
  CONTAINER_RUNTIME := podman
endif

test:   
	@go test -v -cover ./...

server:
	@go run ./cmd/mangtoon/main.go

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

db-new-migrate:
	goose -dir internal/db/migrations create new sql


sqlc: 
	sqlc generate

mock:
	mockgen -package mockdb -destination internal/db/mock/store.go github.com/kavinddd/mangtoon_be/internal/db Store




.PHONY: db-create db-drop db-status db-migrate-up db-migrate-down db-migrate-reset db-new-migrate db-container-create db-container-run sqlc test server mock
