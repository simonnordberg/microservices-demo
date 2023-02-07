FROM golang:1.19-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /productcatalogservice .

FROM alpine:3.17.1 AS without-grpc-health-probe-bin
WORKDIR /src
COPY --from=build /productcatalogservice ./server
COPY products.json .

EXPOSE 3550
ENTRYPOINT ["/src/server"]

FROM without-grpc-health-probe-bin
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.14 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe \
