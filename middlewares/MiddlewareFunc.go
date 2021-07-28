package middlewares

import "net/http"

type MiddlewareFunc func(http.Handler) http.Handler
