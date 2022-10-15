package auth

import (
	"context"
	"easycoding/pkg/ent"
	"easycoding/pkg/errors"
)

type ctxAuthMarker struct{}

var (
	ctxAuthKey = &ctxAuthMarker{}
)

func ToContext(ctx context.Context, userInfo *ent.User) context.Context {
	return context.WithValue(ctx, ctxAuthKey, userInfo)
}

func ExtractContext(ctx context.Context) (*ent.User, error) {
	user, ok := ctx.Value(ctxAuthKey).(*ent.User)
	if !ok || user == nil {
		return nil, errors.ErrUnauthorizedRaw("no authInfo in ctx")
	}
	return user, nil
}
