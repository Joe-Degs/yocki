package server

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

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

func NewServer(addr string) *Server {
	srv := &Server{
		Server: &http.Server{
			Addr: addr,
		},
		router: mux.NewRouter(),
	}
	srv.Server.Handler = srv.router
	return srv
}

// InitRoutes takes a router interface and registers it on the server as a path
// for the access of some resources.
func (l *Server) InitRoutes(s Servicer) error {
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
