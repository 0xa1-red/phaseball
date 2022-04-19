package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/0xa1-red/phaseball/internal/service"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	done := make(chan struct{}, 1)

	s := service.Start(done)

	go func(c chan os.Signal) {
		<-c
		log.Println("Stopping gracefully")
		s.GracefulStop()
	}(c)

	<-done
}
