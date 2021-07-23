package nerdweb

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

/*
RealIP attempts to return the IP address of the caller. The result
will default to the RemoteAddr from http.Request. It will also
check the request headers for an "X-Forwarded-For" value and use
that. This is useful for when requests come through proxies or other
non-direct means.
*/
func RealIP(r *http.Request) string {
	result := r.RemoteAddr
	xForwardedFor := r.Header.Get("X-Forwarded-For")

	if xForwardedFor != "" {
		result = xForwardedFor
	}

	return result
}

/*
ValidateHTTPMethod checks the request METHOD against expectedMethod. If
they do not match an error message is written back to the client.
*/
func ValidateHTTPMethod(r *http.Request, w http.ResponseWriter, expectedMethod string, logger *logrus.Entry) error {
	if r.Method != expectedMethod {

		WriteJSON(logger, w, http.StatusMethodNotAllowed, struct {
			Message string `json:"message"`
		}{
			Message: "method not allowed",
		})

		return fmt.Errorf("invalid method")
	}

	return nil
}
