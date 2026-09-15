package identity

//go:generate go tool mockgen -typed -package identity -destination mock_repositories_test.go gitlab.com/orltom/questionnaire/backend/internal/identity UserRepository
