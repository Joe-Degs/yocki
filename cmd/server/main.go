package main

import (
	"log"

	"github.com/Joe-Degs/yocki/api/hello"
	v1 "github.com/Joe-Degs/yocki/api/v1"
	"github.com/Joe-Degs/yocki/server"
)

func main() {
	srv := server.NewServer(":8080")
	services := []server.Servicer{v1.NewLogService(), hello.Hello{}}
	for _, service := range services {
		if err := srv.InitRoutes(service); err != nil {
			log.Fatal(err)
		}
	}
	log.Fatal(srv.ListenAndServe())
}
