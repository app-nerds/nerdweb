package nerdweb

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/app-nerds/nerdweb/v2/middlewares"
	"github.com/gorilla/mux"
)

/*
BasicWebAppConfig is used to configure a Go template web application router
*/
type BasicWebAppConfig struct {
	AppDirectory  string
	AppFileSystem embed.FS
	Endpoints     Endpoints
	Host          string
	IdleTimeout   int
	ReadTimeout   int
	Version       string
	WriteTimeout  int
}

/*
DefaultBasicWebAppConfig creates a basic web application configuration with default
values. In this configuration the directory holding the front-end JavaScript and
CSS is "app". The HTTP server is configured with an idle timeout of 60 seconds,
and a read and write timeout of 30 seconds.
*/
func DefaultBasicWebAppConfig(host, version string, appFileSystem embed.FS) BasicWebAppConfig {
	return BasicWebAppConfig{
		AppDirectory:  "app",
		AppFileSystem: appFileSystem,
		Endpoints:     make(Endpoints, 0, 20),
		Host:          host,
		IdleTimeout:   60,
		ReadTimeout:   30,
		Version:       version,
		WriteTimeout:  30,
	}
}

/*
NewBasicWebAppRouterAndServer creates a new Gorilla router and HTTP server with
some preconfigured defaults for basic web applications. The HTTP server
is setup to use the resulting router.
*/
func NewBasicWebAppRouterAndServer(config BasicWebAppConfig) (*mux.Router, *http.Server) {
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         config.Host,
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
		Handler:      router,
	}

	fs := http.FileServer(getBasicWebAppFileSystem(config))
	router.Use(middlewares.AccessControl(middlewares.AllowAllOrigins, middlewares.AllowAllMethods, middlewares.AllowAllHeaders))

	sort.Sort(config.Endpoints)

	for _, e := range config.Endpoints {
		if e.HandlerFunc != nil {
			router.HandleFunc(e.Path, e.HandlerFunc).Methods(e.Methods...)
		} else {
			router.Handle(e.Path, e.Handler).Methods(e.Methods...)
		}
	}

	router.PathPrefix("/static/").Handler(fs).Methods(http.MethodGet)
	return router, server
}

/*
NewBasicWebAppServer accepts an existing Gorilla router and returns an HTTP server
with some preconfigured defaults for a basic web application. The HTTP
server is setup to use the resulting router.
*/
func NewBasicWebAppServer(router *mux.Router, config BasicWebAppConfig) *http.Server {
	server := &http.Server{
		Addr:         config.Host,
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
		Handler:      router,
	}

	fs := http.FileServer(getBasicWebAppFileSystem(config))
	router.Use(middlewares.AccessControl(middlewares.AllowAllOrigins, middlewares.AllowAllMethods, middlewares.AllowAllHeaders))

	sort.Sort(config.Endpoints)

	for _, e := range config.Endpoints {
		if e.HandlerFunc != nil {
			router.HandleFunc(e.Path, e.HandlerFunc).Methods(e.Methods...)
		} else {
			router.Handle(e.Path, e.Handler).Methods(e.Methods...)
		}
	}

	router.PathPrefix("/static/").Handler(fs).Methods(http.MethodGet)
	return server
}

func getBasicWebAppFileSystem(config BasicWebAppConfig) http.FileSystem {
	if config.Version == "development" {
		return http.FS(os.DirFS(config.AppDirectory))
	}

	fsys, err := fs.Sub(config.AppFileSystem, config.AppDirectory)

	if err != nil {
		panic("unable to load application static assets: " + err.Error())
	}

	return http.FS(fsys)
}

// func getBasicWebAppRootHandler(config BasicWebAppConfig) func(w http.ResponseWriter, r *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		path := r.URL.Path

// 		if strings.Contains(path, "/main.js") {
// 			w.Header().Set("Content-Type", "text/javascript")
// 			_, _ = w.Write(getFile("main.js", spaConfig))
// 			return
// 		}

// 		if strings.Contains(path, "/manifest.json") {
// 			w.Header().Set("Content-Type", "application/json")
// 			_, _ = w.Write(getFile("manifest.json", spaConfig))
// 			return
// 		}

// 		if strings.Index(path, ".") > -1 {
// 			http.Error(w, "Not found", http.StatusNotFound)
// 			return
// 		}

// 		w.Header().Set("Content-Type", "text/html")
// 		_, _ = w.Write(getFile("index.html", spaConfig))
// 	}
// }

// func getBasicWebappFile(fileName string, config BasicWebAppConfig) []byte {
// 	if config.Version == "development" {
//     path := filepath.Join(config.AppDirectory, fileName)
// 		f, _ := os.Open(path)
// 		b, _ := io.ReadAll(f)
// 		return b
// 	}

// 	if fileName == "main.js" {
// 		return spaConfig.MainJS
// 	}

// 	if fileName == "manifest.json" {
// 		return spaConfig.ManifestJSON
// 	}

// 	return spaConfig.IndexHTML
// }
