package rpc

import (
	"bytes"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"os"
	"os/signal"
	pb "simonnordberg.com/demoshop/productcatalogservice/genproto"
	"simonnordberg.com/demoshop/shared/env"
	"simonnordberg.com/demoshop/shared/logging"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	extraLatency  = time.Duration(1) * time.Second
	cat           pb.ListProductsResponse
	catalogMutex  *sync.Mutex
	log           *logrus.Logger
	reloadCatalog bool
)

func init() {
	log = logging.Init()
	catalogMutex = &sync.Mutex{}
	err := readCatalogFile(&cat)
	if err != nil {
		log.Warnf("could not parse product catalog")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for {
			sig := <-sigs
			log.Printf("Received signal: %s", sig)
			if sig == syscall.SIGUSR1 {
				reloadCatalog = true
				log.Infof("Enable catalog reloading")
			} else {
				reloadCatalog = false
				log.Infof("Disable catalog reloading")
			}
		}
	}()
}

type HealthService struct {
}

func (r *HealthService) Check(_ context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (r *HealthService) Watch(_ *healthpb.HealthCheckRequest, _ healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

type RuntimeService struct {
	pb.UnimplementedRuntimeServiceServer
}

func (r *RuntimeService) GetEnvironment(_ context.Context, _ *pb.Empty) (*pb.Environment, error) {
	return &pb.Environment{Variables: env.GetEnvironmentMap()}, nil
}

type ProductCatalog struct {
	pb.UnimplementedProductCatalogServiceServer
}

func (p *ProductCatalog) Check(_ context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (p *ProductCatalog) Watch(_ *healthpb.HealthCheckRequest, _ healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}

func (p *ProductCatalog) ListProducts(context.Context, *pb.Empty) (*pb.ListProductsResponse, error) {
	time.Sleep(extraLatency)
	return &pb.ListProductsResponse{Products: parseCatalog()}, nil
}

func (p *ProductCatalog) GetProduct(_ context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	time.Sleep(extraLatency)
	var found *pb.Product
	for i := 0; i < len(parseCatalog()); i++ {
		if req.Id == parseCatalog()[i].Id {
			found = parseCatalog()[i]
		}
	}
	if found == nil {
		return nil, status.Errorf(codes.NotFound, "no product with ID %s", req.Id)
	}
	return found, nil
}

func (p *ProductCatalog) SearchProducts(_ context.Context, req *pb.SearchProductsRequest) (*pb.SearchProductsResponse, error) {
	time.Sleep(extraLatency)
	var ps []*pb.Product
	for _, p := range parseCatalog() {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(req.Query)) ||
			strings.Contains(strings.ToLower(p.Description), strings.ToLower(req.Query)) {
			ps = append(ps, p)
		}
	}
	return &pb.SearchProductsResponse{Results: ps}, nil
}

func readCatalogFile(catalog *pb.ListProductsResponse) error {
	catalogMutex.Lock()
	defer catalogMutex.Unlock()
	catalogJSON, err := os.ReadFile("products.json")
	if err != nil {
		log.Fatalf("failed to open product catalog json file: %v", err)
		return err
	}
	if err := jsonpb.Unmarshal(bytes.NewReader(catalogJSON), catalog); err != nil {
		log.Warnf("failed to parse the catalog JSON: %v", err)
		return err
	}
	log.Info("successfully parsed product catalog json")
	return nil
}

func parseCatalog() []*pb.Product {
	if reloadCatalog || len(cat.Products) == 0 {
		err := readCatalogFile(&cat)
		if err != nil {
			return []*pb.Product{}
		}
	}
	return cat.Products
}
