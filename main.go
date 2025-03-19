package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/AdamBrutsaert/basic-go-http-server/internal/store"
)

func main() {
	store := store.New()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: newRouter(store),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, os.Kill)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
