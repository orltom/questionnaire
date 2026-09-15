package rest

import (
	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
)

var _ api.ServerInterface = (*HTTPServer)(nil)

type HTTPServer struct {
	*questionHandler
	*quizHandler
}

func NewHTTPServer(
	questionSvc *application.QuestionService,
	quizSvc *application.QuizService,
) *HTTPServer {
	return &HTTPServer{
		&questionHandler{
			service: questionSvc,
		},
		&quizHandler{
			service: quizSvc,
		},
	}
}
