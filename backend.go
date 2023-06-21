package main

import (
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"digitales-filmmanagement-backend/globals"
	"digitales-filmmanagement-backend/middleware"
)
import chiMiddleware "github.com/go-chi/chi/v5/middleware"

// this function configures and starts the http server
func main() {
	// create a main router handling the different routes for the backend
	router := chi.NewRouter()
	// now enable some middleware globally which is used to identify requests
	// better in case of abuse or debugging
	router.Use(chiMiddleware.RealIP)
	router.Use(chiMiddleware.RequestID)
	router.Use(chiMiddleware.Logger)
	router.Use(middleware.UserInfo(globals.Configuration.OIDC))

	// FIXME: remove preliminary testing route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Second * 600,
		ReadTimeout:  time.Second * 600,
		IdleTimeout:  time.Second * 600,
		Handler:      router,
	}

	// Start the server and log errors that happen while running it
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal().Err(err).Msg("An error occurred while starting the http server")
		}
	}()

	// Set up the signal handling to allow the server to shut down gracefully

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)

	// Block further code execution until the shutdown signal was received
	<-cancelSignal
}
