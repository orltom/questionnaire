package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/infrastructure/database"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/infrastructure/rest"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	// initialize logger...
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// initialize business components...
	questionRepository := database.NewInMemoryQuestionRepository()
	quizRepository := database.NewInMemoryQuizRepository()

	questionSvc := application.NewQuestionService(questionRepository)
	quizSvc := application.NewQuizService(quizRepository, questionRepository)
	server := rest.NewHTTPServer(questionSvc, quizSvc)

	// start HTTP Server ...
	mux := http.NewServeMux()
	api.HandlerFromMux(server, mux)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
