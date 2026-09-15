package rest

import (
	"time"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/quiz/application"
)

var _ api.ServerInterface = (*HTTPServer)(nil)

type HTTPServer struct {
	*questionHandler
	*quizHandler
	*challengeHandler
}

func NewHTTPServer(
	questionSvc *application.QuestionService,
	quizSvc *application.QuizService,
	challengeSvc *application.ChallengeService,
	participationSvc *application.ParticipationService,
) *HTTPServer {
	return &HTTPServer{
		&questionHandler{
			service: questionSvc,
		},
		&quizHandler{
			service: quizSvc,
		},
		&challengeHandler{
			service:       challengeSvc,
			participation: participationSvc,
			now:           time.Now,
		},
	}
}
