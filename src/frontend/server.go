package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"simonnordberg.com/demoshop/shared/env"
	"simonnordberg.com/demoshop/shared/logging"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log                   *logrus.Logger
	port                  string
	whitelistedCurrencies = map[string]bool{
		"USD": true,
		"EUR": true,
		"CAD": true,
		"JPY": true,
		"GBP": true,
		"TRY": true}
)

type frontendServer struct {
	productCatalogService *grpc.ClientConn
	currencyService       *grpc.ClientConn
}

func init() {
	log = logging.Init()
}

func main() {
	flag.Parse()
	port = env.GetEnvOrDefault("PORT", "8080")
	log.Infof("starting server at 0.0.0.0:%s", port)

	svc := new(frontendServer)

	connectGRPC(&svc.productCatalogService, env.GetEnvOrDefault("PRODUCT_CATALOG_SERVICE_ADDR", "productcatalogservice:8502"))
	connectGRPC(&svc.currencyService, env.GetEnvOrDefault("CURRENCY_SERVICE_ADDR", "currencyservice:8502"))

	handler := mux.NewRouter()
	handler.HandleFunc("/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)
	handler.HandleFunc("/debug", svc.debugHandler).Methods(http.MethodGet, http.MethodHead)

	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	if err := http.ListenAndServe("0.0.0.0:"+port, handler); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}

func connectGRPC(conn **grpc.ClientConn, addr string) {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	*conn, err = grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(errors.Wrapf(err, "grpc: failed to connect %s", addr))
	}
}
