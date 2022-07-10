package prometheus

import (
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

func Interceptor() func(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return grpc_prometheus.UnaryServerInterceptor
}
