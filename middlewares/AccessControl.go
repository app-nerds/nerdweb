package middlewares

import "net/http"

const (
	AllowAllOrigins           string = "*"
	AllowAllMethods           string = "POST, GET, OPTIONS, PUT, DELETE"
	AllowAllHeaders           string = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
	AllowHeaderAccept         string = "Accept"
	AllowHeaderContentType    string = "Content-Type"
	AllowHeaderContentLength  string = "Content-Length"
	AllowHeaderAcceptEncoding string = "Accept-Encoding"
	AllowHeaderCSRF           string = "X-CSRF-Token"
	AllowHeaderAuthorization  string = "Authorization"
)

type accessControl struct {
	handler      http.Handler
	allowOrigin  string
	allowMethods string
	allowHeaders string
}

func (ac *accessControl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", ac.allowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", ac.allowMethods)
	w.Header().Set("Access-Control-Allow-Headers", ac.allowHeaders)

	ac.handler.ServeHTTP(w, r)
}

/*
AccessControl wraps an HTTP mux with a middleware that sets
headers for access control and allowed headers.

Example:

  mux := nerdweb.NewServeMux()
  mux.HandleFunc("/endpoint", handler)

  mux.Use(middlewares.AccessControl(middlewares.AllowAllOrigins, middlewares.AllowAllMethods, middlewares.AllowAllHeaders)
*/
func AccessControl(allowOrigin, allowMethods, allowHeaders string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := &accessControl{
				handler:      next,
				allowOrigin:  allowOrigin,
				allowMethods: allowMethods,
				allowHeaders: allowHeaders,
			}

			handler.ServeHTTP(w, r)
		})
	}
}
