package rest

import (
	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
)

var _ api.ServerInterface = (*HTTPServer)(nil)

type HTTPServer struct {
	*catalogHandler
	*questionHandler
	*quizHandler
}

func NewHTTPServer(
	catalogSvc *application.CatalogService,
	questionSvc *application.QuestionService,
	quizSvc *application.QuizService,
) *HTTPServer {
	return &HTTPServer{
		&catalogHandler{
			service: catalogSvc,
		},
		&questionHandler{
			service: questionSvc,
		},
		&quizHandler{
			service: quizSvc,
		},
	}
}
