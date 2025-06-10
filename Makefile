## run http server
run-http:
	go run main.go http_server

## install goose for migration
install-goose:
	go install github.com/pressly/goose/v3/cmd/goose@latest

## install openapi codegen
install-openapi:
	go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

## to generate sql migration file, require NAME to run. e.g make migration NAME=test
create-migration:
	$(call check_defined, NAME)
	goose -dir="./db/migrations" create $(NAME) sql

create-seed:
	$(call check_defined, NAME)
	goose -dir="./db/seeds" create $(NAME) sql

## run migration until latest version
migrate:
	go run main.go migrate

## run seed
seed:
	go run main.go seed

## generate openapi code
generate-openapi:
	go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -config ./oapi_codegen.yml ./api/openapi.yml