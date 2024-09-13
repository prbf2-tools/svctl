generate:
	go generate ./...

build: build-linux build-win

build-linux: generate
	GOOS=linux CGO_ENABLED=0 go build -o bin/svctl ./

build-win: generate
	GOOS=windows CGO_ENABLED=0 go build -o bin/svctl.exe ./

test: generate
	go run gotest.tools/gotestsum@v1.11.0 -- -count=1 ./...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run

grpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    svctl/svctl.proto

