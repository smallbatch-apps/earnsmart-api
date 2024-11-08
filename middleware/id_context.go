package middleware

import (
	"context"
	"errors"
	"log"
)

type contextKey string

const userIDKey contextKey = "userID"

func GetUserIDFromContext(ctx context.Context) (uint64, error) {
	userID, ok := ctx.Value(userIDKey).(uint64)
	log.Println("userID", userID)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}
