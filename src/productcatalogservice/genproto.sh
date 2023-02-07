#!/usr/bin/env bash

PATH=$PATH:$GOPATH/bin
protodir=../../proto

mkdir -p genproto
protoc \
  --go_out=./genproto --go_opt=paths=source_relative \
  --go-grpc_out=./genproto --go-grpc_opt=paths=source_relative \
  --proto_path=$protodir $protodir/demoshop.proto
