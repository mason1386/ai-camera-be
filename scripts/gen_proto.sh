#!/bin/bash
# Install plugins if needed:
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc --proto_path=api/proto \
       --go_out=pkg/pb --go_opt=paths=source_relative \
       --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
       api/proto/*.proto
