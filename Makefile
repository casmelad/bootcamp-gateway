GRPC_GATEWAY_DIR := ~/Documents/repos/Bootcamp-GRPCGateway
GO_INSTALLED := $(shell which go)
PROTOC_INSTALLED := $(shell which protoc)
PGGG_INSTALLED := $(shell which protoc-gen-grpc-gateway 2> /dev/null)
PGG_INSTALLED := $(shell which protoc-gen-go 2> /dev/null)

all: build

generate:
	buf generate;

run:
	go run main.go;

update:
	buf mod update;
	go mod tidy;