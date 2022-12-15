package log_v1

import (
	"encoding/json"
	"errors"
	"net/http"

	. "github.com/Joe-Degs/yocki/api/logservice/json"
	srv "github.com/Joe-Degs/yocki/server"
	ylog "github.com/Joe-Degs/yocki/server/log"
)

// LogService is an implementation of the Servicer interface that provides a
// log aggregator that communicate using json over http.
type LogService struct {
	l *ylog.Log
}

func (l *LogService) handleProduce(w http.ResponseWriter, r *http.Request) {
	var req ProduceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	off, err := l.l.Append(req.Record)
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

func (l *LogService) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req ConsumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rec, err := l.l.Read(req.Offset)
	if errors.Is(err, ylog.ErrOffsetNotFound) {
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

func (l *LogService) Routes() []*srv.Route {
	return []*srv.Route{
		&srv.Route{
			Path:        "/produce",
			Methods:     []string{"POST"},
			HandlerFunc: l.handleProduce,
		},
		&srv.Route{
			Path:        "/consume",
			Methods:     []string{"GET"},
			HandlerFunc: l.handleConsume,
		},
	}
}

func (LogService) Version() string { return "v1" }

func NewLogService() *LogService {
	return &LogService{l: ylog.NewLog()}
}
