package nerdweb

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/app-nerds/nerdweb/v2/middlewares"
	"github.com/gorilla/mux"
)

/*
SPAConfig is used to configure a single page application router
*/
type SPAConfig struct {
	AppDirectory  string
	AppFileSystem embed.FS
	Endpoints     Endpoints
	Host          string
	IdleTimeout   int
	IndexHTML     []byte
	MainJS        []byte
	ManifestJSON  []byte
	ReadTimeout   int
	Version       string
	WriteTimeout  int
}

/*
DefaultSPAConfig creates a single page application configuration with default
values. In this configuration the directory holding the front-end application
is "app". The HTTP server is configured with an idle timeout of 60 seconds,
and a read and write timeout of 30 seconds.
*/
func DefaultSPAConfig(host, version string, appFileSystem embed.FS, indexHTML, mainJS, manifestJSON []byte) SPAConfig {
	return SPAConfig{
		AppDirectory:  "app",
		AppFileSystem: appFileSystem,
		Endpoints:     make(Endpoints, 0, 20),
		Host:          host,
		IdleTimeout:   60,
		IndexHTML:     indexHTML,
		MainJS:        mainJS,
		ManifestJSON:  manifestJSON,
		ReadTimeout:   30,
		Version:       version,
		WriteTimeout:  30,
	}
}

/*
NewSPARouterAndServer creates a new Gorilla router and HTTP server with
some preconfigured defaults for single page applications. The HTTP server
is setup to use the resulting router.
*/
func NewSPARouterAndServer(config SPAConfig) (*mux.Router, *http.Server) {
	router := mux.NewRouter()
	server := &http.Server{
		Addr:         config.Host,
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
		Handler:      router,
	}

	fs := http.FileServer(getClientAppFileSystem(config))
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
	router.HandleFunc(`/{path:[a-zA-Z0-9\-_\/\.]*}`, getRootHandler(config))

	return router, server
}

/*
NewSPAServer accepts an existing Gorilla router and returns an HTTP server
with some preconfigured defaults for a single page application. The HTTP
server is setup to use the resulting router.
*/
func NewSPAServer(router *mux.Router, config SPAConfig) *http.Server {
	server := &http.Server{
		Addr:         config.Host,
		WriteTimeout: time.Second * time.Duration(config.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.IdleTimeout),
		Handler:      router,
	}

	fs := http.FileServer(getClientAppFileSystem(config))
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
	router.HandleFunc(`/{path:[a-zA-Z0-9\-_\/\.]*}`, getRootHandler(config))

	return server
}

func getClientAppFileSystem(spaConfig SPAConfig) http.FileSystem {
	if spaConfig.Version == "development" {
		return http.FS(os.DirFS(spaConfig.AppDirectory))
	}

	fsys, err := fs.Sub(spaConfig.AppFileSystem, spaConfig.AppDirectory)

	if err != nil {
		panic("unable to load application static assets: " + err.Error())
	}

	return http.FS(fsys)
}

func getRootHandler(spaConfig SPAConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.Contains(path, "/main.js") {
			w.Header().Set("Content-Type", "text/javascript")
			_, _ = w.Write(getFile("main.js", spaConfig))
			return
		}

		if strings.Contains(path, "/manifest.json") {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(getFile("manifest.json", spaConfig))
			return
		}

		if strings.Index(path, ".") > -1 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write(getFile("index.html", spaConfig))
	}
}

func getFile(fileName string, spaConfig SPAConfig) []byte {
	if spaConfig.Version == "development" {
		f, _ := os.Open("app/" + fileName)
		b, _ := io.ReadAll(f)
		return b
	}

	if fileName == "main.js" {
		return spaConfig.MainJS
	}

	if fileName == "manifest.json" {
		return spaConfig.ManifestJSON
	}

	return spaConfig.IndexHTML
}
