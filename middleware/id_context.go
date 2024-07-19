package middleware

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(userIDKey).(uint)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}
