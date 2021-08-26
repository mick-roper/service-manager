package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	sigChan := make(chan os.Signal)
	doneChan := make(chan struct{})

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	r := httprouter.New()

	server := &http.Server{
		Addr:    ":45350",
		Handler: r,
	}

	go func() {
		<-sigChan

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelFunc()

		if err := server.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			log.Println(err)
		}

		doneChan <- struct{}{}
	}()

	<-doneChan

	log.Println("quitting!")
}
