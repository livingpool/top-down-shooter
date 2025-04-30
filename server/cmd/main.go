package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/livingpool/top-down-shooter/server/server"
)

var addr = flag.String("addr", ":42069", "game server address")

func main() {
	flag.Parse()

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen on %v\n", *addr)
	}
	log.Printf("listening on ws://%v\n", listener.Addr())

	server := &http.Server{
		Handler:      server.NewGameServer(),
		Addr:         *addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- server.Serve(listener)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	select {
	case err := <-errChan:
		log.Printf("failed to serve %v\n", err)
	case sig := <-sigChan:
		log.Printf("terminating: %v\n", sig)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error shutting down http server: %v\n", err)
	}
}
