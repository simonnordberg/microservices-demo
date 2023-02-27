package main

import (
	"flag"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	pb "simonnordberg.com/demoshop/productcatalogservice/genproto"
	"simonnordberg.com/demoshop/productcatalogservice/rpc"
	"simonnordberg.com/demoshop/shared/env"
	"simonnordberg.com/demoshop/shared/logging"
)

var (
	log  *logrus.Logger
	port string
)

func init() {
	log = logging.Init()
}

func main() {
	flag.Parse()
	port = env.GetEnvOrDefault("PORT", "8502")
	log.Infof("starting grpc server at :%s", port)
	run(port)
	select {}
}

func run(port string) string {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	var srv = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(log)),
			),
		),
	)

	pb.RegisterProductCatalogServiceServer(srv, &rpc.ProductCatalog{})
	pb.RegisterRuntimeServiceServer(srv, &rpc.RuntimeService{})
	healthpb.RegisterHealthServer(srv, &rpc.HealthService{})
	reflection.Register(srv)

	go func() {
		if err := srv.Serve(l); err != nil {
			log.Fatal("failed to serve: %v", err)
		}
	}()
	return l.Addr().String()
}
