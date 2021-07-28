package nerdweb

import (
	"net/http"

	"github.com/app-nerds/nerdweb/middlewares"
)

type ServeMux struct {
	middlewares []middlewares.MiddlewareFunc
	mux         *http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		middlewares: make([]middlewares.MiddlewareFunc, 0, 10),
		mux:         http.NewServeMux(),
	}
}

func (sm *ServeMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	sm.mux.HandleFunc(pattern, handler)
}

func (sm *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		handler http.Handler
	)

	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sm.mux.ServeHTTP(w, r)
	})

	for _, m := range sm.middlewares {
		handler = m(handler)
	}

	handler.ServeHTTP(w, r)
}

func (sm *ServeMux) Use(middlewares ...middlewares.MiddlewareFunc) {
	sm.middlewares = append(sm.middlewares, middlewares...)
}
