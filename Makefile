include ./auth.env

export $(shell sed 's/=.*//' ./auth.env)

.PHONY: infra
infra:
	docker compose -p up

.PHONY: gen
gen:
	go generate ./...
	sqlboiler -c ./sqlboiler.pq.toml psql
#	protoc -I=./ --go_out=./ --go-grpc_out=./ api/proto/*/*.proto

.PHONY: tools
tools:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0
	go install go.uber.org/mock/mockgen@latest
	go install github.com/google/wire/cmd/wire@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: migrate-new
migrate-new:
	goose -dir db/migrations/postgres create $(name) sql

.PHONY: migrate-up
migrate-up:
	goose -dir db/migrations/postgres postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	goose -dir db/migrations/postgres postgres "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" down

.PHONY: test
test:
	@go test -v ../...