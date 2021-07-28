package middlewares

import "net/http"

func realIP(r *http.Request) string {
	result := r.RemoteAddr
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		result = xForwardedFor
	}

	return result
}
