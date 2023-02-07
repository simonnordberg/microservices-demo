package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"net"
	"os"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	pb "simonnordberg.com/demoshop/currencyservice/genproto"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	log  *logrus.Logger
	port = "3550"
)

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
	log.Infof("starting grpc server at :%s", port)
	run(port)
	select {}
}

func run(port string) string {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	var srv *grpc.Server
	srv = grpc.NewServer()
	svc := &currencyService{}

	pb.RegisterCurrencyServiceServer(srv, svc)
	healthpb.RegisterHealthServer(srv, svc)
	reflection.Register(srv)

	go func() {
		if err := srv.Serve(l); err != nil {
			log.Fatal("failed to serve: %v", err)
		}
	}()
	return l.Addr().String()
}

type currencyService struct {
	pb.UnimplementedCurrencyServiceServer
}

func (c *currencyService) Check(_ context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (c *currencyService) Watch(_ *healthpb.HealthCheckRequest, _ healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (*currencyService) GetSupportedCurrencies(_ context.Context, _ *pb.Empty) (*pb.GetSupportedCurrenciesResponse, error) {
	return &pb.GetSupportedCurrenciesResponse{
		CurrencyCodes: []string{"USD", "SEK", "NOK"},
	}, nil
}

func (*currencyService) Convert(_ context.Context, req *pb.CurrencyConversionRequest) (*pb.Money, error) {
	log.Printf("Convert(%v, %v)\n", req.From, req.ToCode)

	return &pb.Money{
		CurrencyCode: req.ToCode,
		Units:        req.From.Units,
		Nanos:        req.From.Nanos,
	}, nil
}
