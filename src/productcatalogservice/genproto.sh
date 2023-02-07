#!/usr/bin/env bash

PATH=$PATH:$GOPATH/bin
protodir=../../proto

# Source: https://github.com/grpc/grpc-go/blob/master/cmd/protoc-gen-go-grpc/README.md
# Consider using --go-grpc_out=require_unimplemented_servers=false

mkdir -p genproto
protoc \
  --go_out=./genproto --go_opt=paths=source_relative \
  --go-grpc_out=./genproto --go-grpc_opt=paths=source_relative \
  --proto_path=$protodir $protodir/demoshop.proto
