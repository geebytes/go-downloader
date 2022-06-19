GOPATH:=$(shell go env GOPATH)
.PHONY: init
init:
    go install google.golang.org/grpc

	go get -u github.com/golang/protobuf/proto
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
.PHONY: proto
proto:
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative pb/downloader.proto

	
# .PHONY: build
# build:
# 	go build -o registerConfiguration *.go

# .PHONY: test
# test:
# 	go test -v ./... -cover

# .PHONY: docker
# docker:
# 	docker build . -t registerConfiguration:latest