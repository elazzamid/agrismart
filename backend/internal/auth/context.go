package auth

import "context"

type contextKey struct{}

var userContextKey contextKey

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userContextKey, userID)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userContextKey).(string)
	return userID, ok && userID != ""
}
