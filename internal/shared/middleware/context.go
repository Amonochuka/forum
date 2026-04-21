package middleware

import "net/http"

type contextKey string

const userKey contextKey = "user"

func GetUserID(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userKey).(int)
	return userID, ok
}