package validate

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	Validate() error
}

func Interceptor() func(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(
		ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		if v, ok := req.(validator); ok {
			if err := v.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, err.Error())
			}
		}
		res, err := handler(ctx, req)
		if err == nil {
			if v, ok := res.(validator); ok {
				if err := v.Validate(); err != nil {
					return res, status.Errorf(codes.Internal, err.Error())
				}
			}
		}
		return res, err
	}
}
