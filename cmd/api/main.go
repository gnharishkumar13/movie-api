package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type application struct {
	config config
	logger *log.Logger
}

type config struct {
	port int
	env  string
}

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 3000, " port to run the application")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Llongfile|log.Ldate|log.Ltime)

	app:= &application{
		config: cfg,
		logger: logger,
	}

	router := app.routes()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
