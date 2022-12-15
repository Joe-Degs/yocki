package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/Joe-Degs/yocki/internal/config"
	"github.com/gorilla/mux"
)

// var config *config.Config

// Route encapsulates a relative path to a resource in the server. It also handles
// access management of the route it defines.
type Route struct {
	http.HandlerFunc
	Path    string
	Methods []string
}

// Servicer provides an interface to types that provide versioned paths to some
// resource/service. The implementer of Servicer is responsible for managing the
// state of the resource/service that it provides.
type Servicer interface {

	// Routes returns a set of api endpoints. An enpoint contains a relative
	// path to some resource/service, the methods of access (http) and a routine for
	// handling the mechanics of that access.
	Routes() []*Route

	// Version returns the version of the set of api endpoints, any type of
	// versioning can be used as long as it can be represented as a string.
	Version() string
}

// Server encapsuslates an http.Server and provides the mechanics for running
// the http server that allows for access to resources registered on the server.
type Server struct {
	*http.Server
	router   *mux.Router
	services []Servicer
}

func NewServer(config config.Config) *Server {
	srv := &Server{
		Server: &http.Server{
			Addr:              config.Address,
			ReadTimeout:       1 * time.Second,
			WriteTimeout:      1 * time.Second,
			IdleTimeout:       30 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		},
		router: mux.NewRouter(),
	}
	srv.Server.Handler = srv.router
	return srv
}

// InitService takes a servicer and initializes it and registers it to start
// serving requests on the server.
func (l *Server) InitService(s Servicer) error {
	if err := l.registerRoutes(s.Version(), s.Routes()); err != nil {
		return err
	}
	return nil
}

func (l *Server) registerRoutes(version string, routes []*Route) error {
	for _, r := range routes {
		versionedEndpoint, err := url.JoinPath("/api/"+version, r.Path)
		if err != nil {
			return err
		}
		l.router.HandleFunc(versionedEndpoint, r.HandlerFunc).Methods(r.Methods...)
	}
	return nil
}

// Start opens the server up for incoming remote connections and servicing of requests.
func (l *Server) Start() error {
	if err := l.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			<-time.After(2 * time.Second)
			return nil
		}
		return err
	}
	return nil
}

// Close cleans up and shuts down the http server
func (l *Server) Close() error {
	return l.Server.Shutdown(context.TODO())
}
