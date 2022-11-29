package main

import (
	"log"

	v1 "github.com/Joe-Degs/yocki/api/v1"
	"github.com/Joe-Degs/yocki/server"
)

func main() {
	srv := server.NewServer(":8080")
	if err := srv.InitRoutes(v1.NewLogService()); err != nil {
		log.Fatal(err)
	}
	log.Fatal(srv.ListenAndServe())
}
