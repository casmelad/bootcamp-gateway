generate:
	buf generate;

run:
	go run main.go;

update:
	buf mod update;
	go mod tidy;

install:
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway;
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2;
	go install google.golang.org/protobuf/cmd/protoc-gen-go;
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc;
	go install github.com/envoyproxy/protoc-gen-validate;