include .env

LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.13.0


get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-user-api

generate-user-api:
	mkdir -p pkg/user_v1
	protoc --proto_path api/user_v1 \
	--go_out=pkg/user_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go.exe \
	--go-grpc_out=pkg/user_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc.exe \
	api/user_v1/user.proto

install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	GOBIN=$(LOCAL_BIN) golangci-lint run ./... --config .golangci.pipeline.yaml

include .env

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD)"

local-migration-status:
	${LOCAL_BIN}/goose.exe -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose.exe -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose.exe -dir ${LOCAL_MIGRATION_DIR} postgres ${PG_DSN} down -v



build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

copy-to-server:
	scp service_linux root@45.12.231.178:

copy-migrations-to-server:
	scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -r ./migrations/** root@45.12.231.178:~/auth/migrations

docker-selectel-build-and-push:
	docker build --no-cache --build-arg grpc_host=localhost --build-arg grpc_port=50051 --build-arg pg_dsn="host=auth-db-1 port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable" --platform linux/amd64 -t cr.selcloud.ru/noskov-sergey/test_sever:v0.0.1 .
	docker login -u token -p CRgAAAAAoUvVJ50Atz2MVMsa09Mi0MNVo9mZRrWD cr.selcloud.ru/noskov-sergey
	docker push cr.selcloud.ru/noskov-sergey/test_sever:v0.0.1

docker-local-pull-and-install:
	docker login -u token -p CRgAAAAAoUvVJ50Atz2MVMsa09Mi0MNVo9mZRrWD cr.selcloud.ru/noskov-sergey
	docker pull cr.selcloud.ru/noskov-sergey/test_sever:v0.0.1
	docker run -p 50051:50051 cr.selcloud.ru/noskov-sergey/test_sever:v0.0.1

docker-selectel-pull-and-install:
	ssh root@45.12.231.178
	docker login -u token -p CRgAAAAAoUvVJ50Atz2MVMsa09Mi0MNVo9mZRrWD cr.selcloud.ru/noskov-sergey
	docker pull cr.selcloud.ru/noskov-sergey/test_sever:v0.0.1
	docker run -p 50051:50051 cr.selcloud.ru/noskov-sergey/test_sever:v0.0.1

docker-auth-local-build-and-compose-up:
	docker build --build-arg grpc_host=localhost --build-arg grpc_port=50051 --build-arg pg_dsn="host=localhost port=54321 dbname=auth user=auth-user password=auth-password sslmode=disable" -t auth:latest .
	docker-compose up -d



#docker run --rm --name auth_db -e POSTGRES_PASSWORD=auth-password -e POSTGRES_USER=auth-users -e POSTGRES_DB=auth -p 54321:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data postgres

#docker run --rm --name auth_db -e POSTGRES_DB=auth -e POSTGRES_USER=auth -e POSTGRES_PASSWORD=authpassword -p 54321:5432 -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data -d postgres