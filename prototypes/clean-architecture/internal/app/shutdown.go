package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Gracefully() {
	quit := make(chan os.Signal, 1)
	defer close(quit)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Gracefully shutting down")
}
