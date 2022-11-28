package ping

import (
	"context"
	ping_pb "easycoding/api/ping"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Service implements ping_pb.pingSrvServer.
type service struct {
	Logger *logrus.Logger
	Tracer trace.Tracer
}

// AuthFuncOverride overrides global AuthFunc, this is used to escape from Auth
// Interceptor.
func (*service) AuthFuncOverride(
	ctx context.Context, _ string) (context.Context, error) {
	return ctx, nil
}

var _ ping_pb.PingSvcServer = (*service)(nil)

func New(logger *logrus.Logger, tracer trace.Tracer) *service {
	return &service{
		Logger: logger,
		Tracer: tracer,
	}
}

// Ping implements ping_pb.pingSrv.Pong
func (s *service) Ping(
	ctx context.Context,
	_ *ping_pb.PingRequest,
) (*ping_pb.PingResponse, error) {
	_, span := s.Tracer.Start(ctx, "ping")
	span.SetAttributes(attribute.Key("method").String("ping"))
	defer span.End()
	return &ping_pb.PingResponse{Res: "pong"}, nil
}
