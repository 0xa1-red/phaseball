package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/0xa1-red/phaseball/internal/config"
	"github.com/0xa1-red/phaseball/internal/service"
)

var (
	configPath string = "./config.yml"
)

func main() {
	flag.StringVar(&configPath, "cfg", "./config.yml", "Path to config file in YAML format")
	flag.Parse()

	if err := config.Init(configPath); err != nil {
		panic(err)
	}

	s := service.New()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	wg := &sync.WaitGroup{}
Loop:
	for {
		select {
		case <-sigint:
			wg.Add(1)
			if err := s.Stop(wg); err != nil {
				log.Printf("Error while shutting down service: %s", err.Error())
			}
			wg.Wait()
			break Loop
		case err := <-s.Errors():
			if err != http.ErrServerClosed {
				log.Printf("ERROR: %s", err.Error())
			}
		}
	}
}
