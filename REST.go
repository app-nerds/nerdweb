package nerdweb

import (
	"net/http"
	"sort"
	"time"

	"github.com/app-nerds/nerdweb/v2/middlewares"
	"github.com/gorilla/mux"
)

/*
RESTConfig is used to configure a router for basic REST servers
*/
type RESTConfig struct {
	Endpoints    Endpoints
	Host         string
	IdleTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

/*
DefaultRESTConfig creates a REST configuration with default
values. In this configuration the HTTP server is configured with an idle timeout of 60 seconds,
and a read and write timeout of 30 seconds.
*/
func DefaultRESTConfig(host string) RESTConfig {
	return RESTConfig{
		Endpoints:    make(Endpoints, 0, 20),
		Host:         host,
		IdleTimeout:  60,
		ReadTimeout:  30,
		WriteTimeout: 30,
	}
}

/*
NewRESTRouterAndServer creates a new Gorilla router and HTTP server with
some preconfigured defaults for REST applications. The HTTP server
is setup to use the resulting router.
*/
func NewRESTRouterAndServer(config RESTConfig) (*mux.Router, *http.Server) {
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         config.Host,
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
		Handler:      router,
	}

	router.Use(middlewares.AccessControl(middlewares.AllowAllOrigins, middlewares.AllowAllMethods, middlewares.AllowAllHeaders))

	sort.Sort(config.Endpoints)

	for _, e := range config.Endpoints {
		if e.HandlerFunc != nil {
			router.HandleFunc(e.Path, e.HandlerFunc).Methods(e.Methods...)
		} else {
			router.Handle(e.Path, e.Handler).Methods(e.Methods...)
		}
	}

	return router, server
}
