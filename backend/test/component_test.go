package test

import (
	"net/http/httptest"
	"testing"
	"time"
	"uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/server"
)

func TestComponentTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Questionnaire Component Tests")
}

var _ = Describe("Quiz Challenge Workflow", Ordered, func() {
	var (
		ownerClient     *TestClient
		testServer      *httptest.Server
		quizID          uuid.UUID
		questionID      uuid.UUID
		challengeID     uuid.UUID
		correctAnswerID uuid.UUID
	)

	BeforeAll(func() {
		deps := server.SetupServiceLayer()
		handler := server.BuildMux(deps)
		testServer = httptest.NewServer(handler)
		ownerClient = NewTestClient(testServer.URL, uuid.New())
	})

	AfterAll(func() {
		if testServer != nil {
			testServer.Close()
		}
	})

	Context("Setup", func() {
		It("creates quiz with question", func() {
			var err error
			quizID, err = ownerClient.CreateQuiz("Science Quiz", "Basic science", "public")
			Expect(err).NotTo(HaveOccurred())

			answers := []api.Answer{
				{Description: "Mitochondria", Position: 1, Correct: true},
				{Description: "Nucleus", Position: 2, Correct: false},
			}
			questionID, err = ownerClient.CreateQuestion("Powerhouse of cell?", "public", answers)
			Expect(err).NotTo(HaveOccurred())

			question, err := ownerClient.GetQuestion(questionID)
			Expect(err).NotTo(HaveOccurred())
			Expect(question.Answers).NotTo(BeEmpty())

			for _, answer := range question.Answers {
				if answer.Correct {
					correctAnswerID = answer.Id
					break
				}
			}
			Expect(correctAnswerID).NotTo(Equal(uuid.UUID{}))

			err = ownerClient.AddQuestionToQuiz(quizID, questionID, 1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates public challenge", func() {
			var err error
			challengeID, err = ownerClient.CreateChallenge(quizID, time.Now(), "public", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(challengeID).NotTo(Equal(uuid.UUID{}))
		})
	})

	Context("Participation", func() {
		It("different user can submit answers", func() {
			participantClient := ownerClient.AsUser(uuid.NewV7())

			err := participantClient.SubmitAnswer(challengeID, questionID, correctAnswerID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("owner can view results", func() {
			results, err := ownerClient.GetChallengeResults(challengeID)
			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			Expect(results).To(HaveLen(1))
			Expect(results[0].Points).To(BeNumerically(">", 0))
		})
	})
})
