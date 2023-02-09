package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log                   *logrus.Logger
	port                  = "3550"
	whitelistedCurrencies = map[string]bool{
		"USD": true,
		"EUR": true,
		"CAD": true,
		"JPY": true,
		"GBP": true,
		"TRY": true}
)

type frontendServer struct {
	productCatalogSvcAddr string
	productCatalogSvcConn *grpc.ClientConn

	currencySvcAddr string
	currencySvcConn *grpc.ClientConn
}

func init() {
	log = logrus.New()
	log.Formatter = &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	}
	log.Out = os.Stdout
}

func main() {
	flag.Parse()
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	log.Infof("starting server at 0.0.0.0:%s", port)
	ctx := context.Background()

	svc := new(frontendServer)

	svc.productCatalogSvcAddr = getEnvOrDefault("PRODUCT_CATALOG_SERVICE_ADDR", "localhost:3551")
	svc.currencySvcAddr = getEnvOrDefault("CURRENCY_SERVICE_ADDR", "localhost:3550")

	mustConnGRPC(ctx, &svc.productCatalogSvcConn, svc.productCatalogSvcAddr)
	mustConnGRPC(ctx, &svc.currencySvcConn, svc.currencySvcAddr)

	handler := mux.NewRouter()
	handler.HandleFunc("/", svc.homeHandler).Methods(http.MethodGet, http.MethodHead)

	if err := http.ListenAndServe("0.0.0.0:"+port, handler); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func mustConnGRPC(ctx context.Context, conn **grpc.ClientConn, addr string) {
	var err error
	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	*conn, err = grpc.DialContext(ctx, addr, grpc.WithInsecure())
	if err != nil {
		panic(errors.Wrapf(err, "grpc: failed to connect %s", addr))
	}
}
