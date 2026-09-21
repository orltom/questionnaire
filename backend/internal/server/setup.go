package server

import (
	"net/http"

	"gitlab.com/orltom/questionnaire/backend/api"
	middleware "gitlab.com/orltom/questionnaire/backend/internal/http"
	"gitlab.com/orltom/questionnaire/backend/internal/identity"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/infrastructure/persistence"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/infrastructure/rest"
)

type ServerDependency struct {
	QuestionSvc      *application.QuestionService
	QuizSvc          *application.QuizService
	ChallengeSvc     *application.ChallengeService
	ParticipationSvc *application.ParticipationService
	IdentitySvc      *identity.UserService
}

func SetupServiceLayer() *ServerDependency {
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

	return &ServerDependency{
		QuestionSvc:      questionSvc,
		QuizSvc:          quizSvc,
		ChallengeSvc:     challengeSvc,
		ParticipationSvc: participationSvc,
		IdentitySvc:      identitySvc,
	}
}

func BuildMux(deps *ServerDependency) http.Handler {
	server := rest.NewHTTPServer(deps.QuestionSvc, deps.QuizSvc, deps.ChallengeSvc, deps.ParticipationSvc)
	mux := http.NewServeMux()
	api.HandlerFromMux(server, mux)
	return middleware.ExtractUser(middleware.WriteAccessLog(mux), deps.IdentitySvc)
}
