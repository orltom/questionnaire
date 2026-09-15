package application

//go:generate go tool mockgen -typed -package application -destination mock_repositories_test.go gitlab.com/orltom/questionnaire/backend/internal/quiz/application QuestionRepository,QuizRepository,QuestionLookupService
