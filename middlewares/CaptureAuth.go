package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
CaptureAuth captures an authorization token from an Authorization
header and stored it in a context variable named "authtoken". This
middleware expect the header to be in the format of:

  Authorization: Bearer <token here>

If the header format is invalid, the provided error method is called.
Here is an example:

  onInvalidHeader = func(logger *logrus.Entry, w http.ResponseWriter) {
    result := map[string]string{
      "error": "invalid JWT header!",
    }

    nerdweb.WriteJSON(logger, w, http.StatusBadRequest, result)
  }

  // Now, in your handler definition
  http.HandleFunc("/endpoint", middlewares.CaptureAuth(handlerFunc, logger, onInvalidHeader))
*/
func CaptureAuth(next http.HandlerFunc, logger *logrus.Entry, onInvalidHeader func(logger *logrus.Entry, w http.ResponseWriter)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		auth := strings.SplitN(authHeader, " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" {
			logger.Error("invalid JWT authorization header. Expected 'Bearer <token here>'")
			onInvalidHeader(logger, w)
			return
		}

		token := auth[1]
		ctx := context.WithValue(r.Context(), "authtoken", token)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
