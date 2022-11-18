package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type version string

type route struct {
	endpoint string
	methods  []string
	f        func(http.ResponseWriter, *http.Request)
}

type httpServer struct {
	router *mux.Router
	routes map[version][]*route

	Log *Log
}

func newHttpServer() *httpServer {
	srv := &httpServer{
		router: mux.NewRouter(),
		routes: make(map[version][]*route),
		Log:    NewLog(),
	}
	if err := srv.initRoutes(); err != nil {
		log.Fatal(err)
	}
	return srv
}

// initRoutes gives the ability to version the endpoints of the server
func (h *httpServer) initRoutes() error {
	h.routes["/api/v1"] = []*route{
		&route{
			endpoint: "/produce",
			methods:  []string{"POST"},
			f:        h.handleProduce,
		},
		&route{
			endpoint: "/consume",
			methods:  []string{"GET"},
			f:        h.handleConsume,
		},
	}

	if err := h.registerRoutes(); err != nil {
		return err
	}
	return nil
}

func (h *httpServer) registerRoutes() error {
	for v, routes := range h.routes {
		for _, r := range routes {
			versionedEndpoint, err := url.JoinPath(string(v), r.endpoint)
			if err != nil {
				return err
			}
			h.router.HandleFunc(versionedEndpoint, r.f).Methods(r.methods...)
		}
	}
	return nil
}

func NewHTTPServer(addr string) *http.Server {
	srv := newHttpServer()
	return &http.Server{
		Addr:    addr,
		Handler: srv.router,
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func (h *httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	off, err := h.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ProduceResponse{off}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rec, err := h.Log.Read(req.Offset)
	if errors.Is(err, ErrOffsetNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res := ConsumeResponse{rec}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
