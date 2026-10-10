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

type Repositories struct {
	QuestionRepository      application.QuestionRepository
	QuestionLookup          application.QuestionLookup
	QuizRepository          application.QuizRepository
	ChallengeRepository     application.ChallengeRepository
	ParticipationRepository application.ParticipationRepository
	UserRepository          identity.UserRepository
}

type ServerDependency struct {
	QuestionSvc      *application.QuestionService
	QuizSvc          *application.QuizService
	ChallengeSvc     *application.ChallengeService
	ParticipationSvc *application.ParticipationService
	IdentitySvc      *identity.UserService
}

func BuildInMemoryRepositories() *Repositories {
	questionRepository := persistence.NewInMemoryQuestionRepository()
	quizRepository := persistence.NewInMemoryQuizRepository()
	challengeRepository := persistence.NewInMemoryChallengeRepository()
	participationRepository := persistence.NewInMemoryParticipationRepository()
	userRepository := persistence.NewInMemoryUserRepository()
	return &Repositories{
		QuestionRepository:      questionRepository,
		QuestionLookup:          questionRepository,
		QuizRepository:          quizRepository,
		ChallengeRepository:     challengeRepository,
		ParticipationRepository: participationRepository,
		UserRepository:          userRepository,
	}
}

func BuildServerDependency(repos *Repositories) *ServerDependency {
	identitySvc := identity.NewUserService(repos.UserRepository)
	questionSvc := application.NewQuestionService(repos.QuestionRepository)
	quizSvc := application.NewQuizService(repos.QuizRepository, repos.QuestionLookup)
	challengeSvc := application.NewChallengeService(repos.ChallengeRepository, repos.QuizRepository, repos.QuestionRepository)
	participationSvc := application.NewParticipationService(repos.ParticipationRepository, repos.ChallengeRepository)

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
