package main

import (
	"context"
	"google.golang.org/grpc"
	pb "simonnordberg.com/demoshop/frontend/genproto"
)

func (fe *frontendServer) getCurrencies(ctx context.Context) ([]string, error) {
	currs, err := pb.NewCurrencyServiceClient(fe.currencyService).GetSupportedCurrencies(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	var out []string
	for _, c := range currs.CurrencyCodes {
		if _, ok := whitelistedCurrencies[c]; ok {
			out = append(out, c)
		}
	}
	return out, nil
}

func (fe *frontendServer) getProducts(ctx context.Context) ([]*pb.Product, error) {
	resp, err := pb.NewProductCatalogServiceClient(fe.productCatalogService).ListProducts(ctx, &pb.Empty{})
	return resp.GetProducts(), err
}

func (fe *frontendServer) convertCurrency(ctx context.Context, money *pb.Money, currency string) (*pb.Money, error) {
	if money.GetCurrencyCode() == currency {
		return money, nil
	}

	return pb.NewCurrencyServiceClient(fe.currencyService).
		Convert(ctx, &pb.CurrencyConversionRequest{
			From:   money,
			ToCode: currency})
}

func (fe *frontendServer) getEnvironments(ctx context.Context) error {
	services := []*grpc.ClientConn{fe.currencyService, fe.productCatalogService}
	for _, s := range services {
		env, err := pb.NewRuntimeServiceClient(s).GetEnvironment(ctx, &pb.Empty{})
		if err != nil {
			return err
		}
		log.Debugf("Environment for: %s", s.Target())
		for k, v := range env.Variables {
			log.Debugf("%s = %s", k, v)
		}
	}

	return nil
}
