package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

/*
Allow verifies if the caller method matches the provided method.

Example:

  mux := nerdweb.NewServeMux()
  mux.HandleFunc("/endpoint", middlewares.Allow(myHandler, http.MethodPost))

If the caller's method does not match what is allowed, the string
"method not allowed" is written back to the caller.
*/
func Allow(next http.HandlerFunc, allowedMethod string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) != strings.ToLower(allowedMethod) {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintf(w, "%s", "method not allowed")

			return
		}

		if r.Method == http.MethodOptions {
			return
		}

		next.ServeHTTP(w, r)
	})
}
