package hello

import (
	"net/http"

	"github.com/Joe-Degs/yocki/server"
)

type Hello struct{}

func (Hello) Routes() []*server.Route {
	return []*server.Route{
		{
			Path:    "",
			Methods: []string{"GET"},
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Hello, World!\n"))
			},
		},
	}
}

func (Hello) Version() string { return "hello" }
