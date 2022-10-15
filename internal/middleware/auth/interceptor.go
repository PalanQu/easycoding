package auth

import (
	"context"

	pkg_auth "easycoding/pkg/auth"
	"easycoding/pkg/ent"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Interceptor() grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(authInterceptor)
}

func authInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	userInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	newCtx := pkg_auth.ToContext(ctx, userInfo)

	return newCtx, nil
}

func parseToken(token string) (*ent.User, error) {
	// hard code auth func, use your own func in production
	return &ent.User{
		ID:   123,
		Name: "foo",
	}, nil
}
