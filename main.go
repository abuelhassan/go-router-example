package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/abuelhassan/go-router-example/handler"
	"github.com/abuelhassan/go-router-example/router"
)

var routes = []router.Route{
	{
		Method:  http.MethodGet,
		Pattern: "/health",
		Handler: handler.HealthCheck,
	},
}

func main() {
	const (
		defaultPort         = 8080
		gracefulShutdownSec = 60
	)

	port := *flag.Uint("port", defaultPort, "server port")
	flag.Parse()

	rtr := router.New(routes, handler.NotFound)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: rtr,
	}

	go func() {
		log.Printf("starting server at port :%d\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Supports graceful shutdown for SIGINT, but not for SIGQUIT or SIGTERM.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownSec*time.Second)
	defer cancel()

	log.Printf("Shutting down...\nThis may take up to %d seconds.\n", gracefulShutdownSec)

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	os.Exit(0)
}
