package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
)

var addr = fmt.Sprint(":", os.Getenv("PORT"))

func main() {
	sigChan := make(chan os.Signal)
	doneChan := make(chan struct{})

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	r := httprouter.New()

	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		<-sigChan

		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelFunc()

		if err := server.Shutdown(ctx); err != http.ErrServerClosed {
			log.Println(err)
		}

		doneChan <- struct{}{}
	}()

	go func() {
		log.Println("server listening at ", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	<-doneChan

	log.Println("quitting!")
}
