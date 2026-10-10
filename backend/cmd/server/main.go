package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gitlab.com/orltom/questionnaire/backend/internal/server"
)

func main() {
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	repos := server.BuildInMemoryRepositories()
	services := server.BuildServerDependency(repos)
	handler := server.BuildMux(services)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
