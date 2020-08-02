package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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
	const defaultPort = 8080
	port := *flag.Uint("port", defaultPort, "server port")

	rtr := router.New(routes, handler.NotFound)

	log.Printf("starting server at port :%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), rtr)
	if err != nil {
		log.Fatal(err)
	}
}
