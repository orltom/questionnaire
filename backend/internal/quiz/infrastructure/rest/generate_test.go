package rest

//go:generate go tool mockgen -typed -package rest -destination mock_service_test.go gitlab.com/orltom/questionnaire/backend/internal/quiz/infrastructure/rest catalogService,quizService,questionService
