package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Joe-Degs/yocki/api/hello"
	jv1 "github.com/Joe-Degs/yocki/api/logservice/json/v1"
	"github.com/Joe-Degs/yocki/internal/config"
	"github.com/Joe-Degs/yocki/server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "specify path to server toml config")
}

func handleSignals(shutdownc <-chan io.Closer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	for {
		sig, ok := (<-c).(syscall.Signal)
		if !ok {
			continue
		}
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			donec := make(chan struct{})
			go func() {
				cl := <-shutdownc
				log.Print("shutting down server...")
				if err := cl.Close(); err != nil {
					log.Fatal(err)
				}
				donec <- struct{}{}
			}()
			select {
			case <-donec:
				os.Exit(0)
			case <-time.After(3 * time.Second):
				log.Fatal("Server took too long to shutdown, exiting now!")
			}
		default:
			log.Fatal("Recieved another signal, should not happen")
		}
	}
}

func main() {
	flag.Parse()
	if err := config.Init(configPath); err != nil {
		flag.Usage()
		os.Exit(1)
	}

	srv := server.NewServer(config.GetConfig())
	services := []server.Servicer{jv1.NewLogService(), hello.Hello{}}
	for _, service := range services {
		if err := srv.InitService(service); err != nil {
			log.Fatal(err)
		}
	}

	shutdownc := make(chan io.Closer, 1)
	go handleSignals(shutdownc)
	shutdownc <- srv
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

	close(shutdownc)
	os.Exit(0)
}
