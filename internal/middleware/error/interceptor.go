package error

import (
	"context"
	"easycoding/pkg/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Interceptor(logger *logrus.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		reqinfo *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		res, err := handler(ctx, req)
		if err == nil {
			return res, err
		}
		var code codes.Code

		switch {
		case errors.ErrorIs(err, errors.InternalError):
			code = codes.Internal
		case errors.ErrorIs(err, errors.InvalidError):
			code = codes.InvalidArgument
		case errors.ErrorIs(err, errors.NotFoundError):
			code = codes.NotFound
		case errors.ErrorIs(err, errors.PermissionError):
			code = codes.PermissionDenied
		case errors.ErrorIs(err, errors.UnauthorizedError):
			code = codes.Unauthenticated
		default:
			logger.WithError(err).WithField("method", reqinfo.FullMethod).
				Warn("invalid err, without using easycoding/pkg/errors")
			return res, err
		}
		s := status.New(code, err.Error())
		return res, s.Err()
	}
}
