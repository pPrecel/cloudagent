.PHONY: build
build:
	go build -o .out/gardenagent cmd/main.go

.PHONY: protobuf
protobuf:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		internal/grpc/proto/route.proto
