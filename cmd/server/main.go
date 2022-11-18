package main

import (
	"log"

	"github.com/Joe-Degs/yocki/internal/server"
)

func main() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
