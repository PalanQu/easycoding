package ping

import (
	"context"
	ping_pb "easycoding/api/ping"

	"github.com/sirupsen/logrus"
)

// Service implements ping_pb.pingSrvServer.
type service struct {
	Logger *logrus.Logger
}

// AuthFuncOverride overrides global AuthFunc, this is used to escape from Auth
// Interceptor.
func (*service) AuthFuncOverride(
	ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}

var _ ping_pb.PingSvcServer = (*service)(nil)

func New(logger *logrus.Logger) *service {
	return &service{
		Logger: logger,
	}
}

// Ping implements ping_pb.pingSrv.Pong
func (s *service) Ping(
	ctx context.Context,
	_ *ping_pb.PingRequest,
) (*ping_pb.PingResponse, error) {
	return &ping_pb.PingResponse{Res: "pong"}, nil
}
