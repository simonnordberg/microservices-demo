FROM golang:1.19-alpine AS build

WORKDIR /src

RUN apk update && apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

COPY src/frontend ./
COPY proto ./proto

RUN ./genproto.sh
RUN go mod download

RUN go build -o /server .

FROM alpine:3.17.1 AS without-grpc-health-probe-bin
WORKDIR /src
COPY --from=build /server ./server
COPY ./src/frontend/templates ./templates

EXPOSE 8080
ENTRYPOINT ["/src/server"]

FROM without-grpc-health-probe-bin
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.14 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe \
