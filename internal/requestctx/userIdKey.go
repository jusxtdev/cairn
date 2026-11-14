package requestctx

import (
	"context"
)

type userIdKey struct {}

func WithUserId(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userIdKey{}, userId)
}

func GetUserId(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(userIdKey{}).(string)
	return userId, ok
}