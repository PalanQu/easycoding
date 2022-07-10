package service

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	pet_pb "easycoding/api/pet"
	ping_pb "easycoding/api/ping"
	pet_svc "easycoding/internal/service/pet"
	ping_svc "easycoding/internal/service/ping"
)

const (
	// 500 M
	// TODO: add to config
	maxMsgSize = 500 * 1024 * 1024
)

// RegisterServers register grpc services.
func RegisterServers(grpcServer *grpc.Server, logger *logrus.Logger, db *gorm.DB) {
	ping_pb.RegisterPingSvcServer(grpcServer, ping_svc.New(logger))
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
	type RegisterHandlerFromEndpoint func(
		ctx context.Context,
		gwmux *runtime.ServeMux,
		endpoint string,
		opts []grpc.DialOption) (err error)

	funcs := []RegisterHandlerFromEndpoint{
		ping_pb.RegisterPingSvcHandlerFromEndpoint,
		pet_pb.RegisterPetStoreSvcHandlerFromEndpoint,
	}

	for _, f := range funcs {
		f(ctx, gwmux, grpcAddr, dopts)
	}
}
