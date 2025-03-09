FROM golang:1.22 as builder

WORKDIR /

RUN apt-get update && apt-get install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="${PATH}:/go/bin"

WORKDIR /build 