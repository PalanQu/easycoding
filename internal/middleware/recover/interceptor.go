package recover

import (
	"context"
	"easycoding/pkg/errors"
	"fmt"
	"runtime"

	recovery_middleware "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

func recoverFunc(p interface{}) (err error) {
	_, f1, l1, _ := runtime.Caller(6)
	_, f2, l2, _ := runtime.Caller(7)
	detail := fmt.Sprintf("%v\n%s, %d\n%s,%d", p, f1, l1, f2, l2)
	return errors.ErrInternalRaw(detail)
}

func Interceptor() func(
	ctx context.Context,
	req interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return recovery_middleware.UnaryServerInterceptor(
		recovery_middleware.WithRecoveryHandler(recoverFunc))
}
