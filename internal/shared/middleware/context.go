package middleware

import (
	"context"
	"net/http"
)

type contextKey string

const userKey contextKey = "user"

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userKey).(int)
	return userID, ok
}

func WithUserID(r *http.Request, userID int) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, userKey, userID)
	return r.WithContext(ctx)
}