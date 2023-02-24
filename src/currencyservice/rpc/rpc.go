package rpc

import (
	"context"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	pb "simonnordberg.com/demoshop/currencyservice/genproto"
	"simonnordberg.com/demoshop/shared/env"
)

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

type CurrencyService struct {
	pb.UnimplementedCurrencyServiceServer
}

func (*CurrencyService) GetSupportedCurrencies(_ context.Context, _ *pb.Empty) (*pb.GetSupportedCurrenciesResponse, error) {
	return &pb.GetSupportedCurrenciesResponse{
		CurrencyCodes: []string{"USD", "SEK", "NOK"},
	}, nil
}

func (*CurrencyService) Convert(_ context.Context, req *pb.CurrencyConversionRequest) (*pb.Money, error) {
	return &pb.Money{
		CurrencyCode: req.ToCode,
		Units:        req.From.Units,
		Nanos:        req.From.Nanos,
	}, nil
}
