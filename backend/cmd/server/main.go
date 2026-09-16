package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gitlab.com/orltom/questionnaire/backend/api"
	middelware "gitlab.com/orltom/questionnaire/backend/internal/http"
	"gitlab.com/orltom/questionnaire/backend/internal/identity"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/infrastructure/persistence"
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
	jsonHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(jsonHandler)
	slog.SetDefault(logger)

	// initialize business components...
	questionRepository := persistence.NewInMemoryQuestionRepository()
	quizRepository := persistence.NewInMemoryQuizRepository()
	challengeRepository := persistence.NewInMemoryChallengeRepository()
	participationRepository := persistence.NewInMemoryParticipationRepository()
	userRepository := persistence.NewInMemoryUserRepository()

	identitySvc := identity.NewUserService(userRepository)
	questionSvc := application.NewQuestionService(questionRepository)
	quizSvc := application.NewQuizService(quizRepository, questionRepository)
	challengeSvc := application.NewChallengeService(challengeRepository, quizRepository, questionRepository)
	participationSvc := application.NewParticipationService(participationRepository, challengeRepository)
	server := rest.NewHTTPServer(questionSvc, quizSvc, challengeSvc, participationSvc)

	// start HTTP Server ...
	mux := http.NewServeMux()
	api.HandlerFromMux(server, mux)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           middelware.ExtractUser(middelware.WriteAccessLog(mux), identitySvc),
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
