package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"

	pet_pb "easycoding/api/pet"
	ping_pb "easycoding/api/ping"
	pet_svc "easycoding/internal/service/pet"
	ping_svc "easycoding/internal/service/ping"
	"easycoding/pkg/ent"
)

const (
	// 500 M
	// TODO: add to config
	maxMsgSize = 500 * 1024 * 1024
)

type RegisterHandlerFromEndpoint func(
	ctx context.Context,
	gwmux *runtime.ServeMux,
	endpoint string,
	opts []grpc.DialOption) (err error)

var endpointFuncs = []RegisterHandlerFromEndpoint{
	ping_pb.RegisterPingSvcHandlerFromEndpoint,
	pet_pb.RegisterPetStoreSvcHandlerFromEndpoint,
}

// RegisterServers register grpc services.
func RegisterServers(grpcServer *grpc.Server, logger *logrus.Logger, db *ent.Client, tracer trace.Tracer) {
	ping_pb.RegisterPingSvcServer(grpcServer, ping_svc.New(logger, tracer))
	pet_pb.RegisterPetStoreSvcServer(grpcServer, pet_svc.New(logger, db))
}

// RegisterHandlers register grpc gateway handlers.
func RegisterHandlers(gwmux *runtime.ServeMux, grpcAddr string) {
	ctx := context.Background()
	dopts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	}
	for _, f := range endpointFuncs {
		f(ctx, gwmux, grpcAddr, dopts)
	}
}
