package log

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	logMethodBlackList = []string{}
)

func Interceptor(logger *logrus.Logger) func(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	ops := []grpc_logrus.Option{
		grpc_logrus.WithLevels(levelFunc),
		grpc_logrus.WithDecider(decider),
	}
	entry := logrus.NewEntry(logger)

	logInterceptorBefore := createBeforeInterceptor(entry)
	logInterceptorAfter := createAfterInterceptor(entry)

	return grpc_middleware.ChainUnaryServer(
		logInterceptorBefore,
		grpc_logrus.UnaryServerInterceptor(entry, ops...),
		logInterceptorAfter,
	)
}

func createBeforeInterceptor(entry *logrus.Entry) func(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		newCtx := ctxlogrus.ToContext(ctx, entry)
		fields := logrus.Fields{}
		grpc_logrus.AddFields(newCtx, fields)
		return handler(newCtx, req)
	}
}

func createAfterInterceptor(entry *logrus.Entry) func(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(ctx context.Context,
		req interface{},
		_ *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		fields := make(logrus.Fields)
		res, err := handler(ctx, req)
		s, ok := status.FromError(err)
		if ok {
			code := runtime.HTTPStatusFromCode(s.Code())
			fields["code"] = code
		}
		if err != nil && ok {
			msg, _ := new(runtime.JSONPb).Marshal(s.Proto())
			fields["resp"] = string(msg)
		}
		grpc_logrus.AddFields(ctx, fields)
		return res, err
	}
}

func levelFunc(c codes.Code) logrus.Level {
	switch c {
	case codes.Internal:
		return logrus.ErrorLevel
	case codes.InvalidArgument,
		codes.Unauthenticated,
		codes.PermissionDenied:
		return logrus.WarnLevel
	default:
		return logrus.InfoLevel
	}
}

func decider(fullMethodName string, err error) bool {
	if err != nil {
		return true
	}
	for _, name := range logMethodBlackList {
		if fullMethodName == name {
			return false
		}
	}
	return true
}
