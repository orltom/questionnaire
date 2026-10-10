package test

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"
	"uuid"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"gitlab.com/orltom/questionnaire/backend/api"
	"gitlab.com/orltom/questionnaire/backend/internal/identity"
	"gitlab.com/orltom/questionnaire/backend/internal/server"
)

func TestComponentTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Questionnaire Component Tests")
}

var _ = Describe("Quiz Challenge Workflow", Ordered, func() {
	var (
		client          *TestClient
		testServer      *httptest.Server
		quizID          uuid.UUID
		questionID      uuid.UUID
		challengeID     uuid.UUID
		correctAnswerID uuid.UUID
		wrongAnswerID   uuid.UUID
		userA           identity.User
		userB           identity.User
	)

	BeforeAll(func() {
		repos := server.BuildInMemoryRepositories()

		userA, _ = identity.NewUser("https://accounts.google.com", "foo", "foo")
		err := repos.UserRepository.Create(context.Background(), userA)
		Expect(err).To(BeNil())

		userB, _ = identity.NewUser("https://accounts.google.com", "boo", "boo")
		err = repos.UserRepository.Create(context.Background(), userB)
		Expect(err).To(BeNil())

		services := server.BuildServerDependency(repos)
		handler := server.BuildMux(services)
		testServer = httptest.NewServer(handler)
		client = NewTestClient(testServer.URL, "admin")
	})

	AfterAll(func() {
		if testServer != nil {
			testServer.Close()
		}
	})

	Context("Setup", func() {
		It("creates quiz with question", func() {
			var err error
			quizID, err = client.CreateQuiz("Science Quiz", "Basic science", "public")
			Expect(err).NotTo(HaveOccurred())

			answers := []api.Answer{
				{Description: "Mitochondria", Position: 1, Correct: true},
				{Description: "Nucleus", Position: 2, Correct: false},
			}
			questionID, err = client.CreateQuestion("Powerhouse of cell?", "public", answers)
			Expect(err).NotTo(HaveOccurred())

			question, err := client.GetQuestion(questionID)
			Expect(err).NotTo(HaveOccurred())
			Expect(question.Answers).NotTo(BeEmpty())

			for _, answer := range question.Answers {
				if answer.Correct {
					correctAnswerID = answer.Id
				} else {
					wrongAnswerID = answer.Id
				}
			}
			Expect(correctAnswerID).NotTo(Equal(uuid.UUID{}))

			err = client.AddQuestionToQuiz(quizID, questionID, 1)
			Expect(err).NotTo(HaveOccurred())
		})

		It("creates public challenge", func() {
			var err error
			challengeID, err = client.CreateChallenge(quizID, time.Now(), "public", nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(challengeID).NotTo(Equal(uuid.UUID{}))
		})
	})

	Context("Participation", func() {
		It("user foo can submit answers", func() {
			c := client.AsUser(userA.Name())

			err := c.SubmitAnswer(challengeID, questionID, correctAnswerID)
			Expect(err).NotTo(HaveOccurred())
		})
		It("user boo can submit answers", func() {
			c := client.AsUser(userB.Name())

			err := c.SubmitAnswer(challengeID, questionID, wrongAnswerID)
			Expect(err).NotTo(HaveOccurred())
		})

		It("owner can view results", func() {
			results, err := client.GetChallengeResults(challengeID)

			Expect(err).NotTo(HaveOccurred())
			Expect(results).NotTo(BeEmpty())
			Expect(results).To(HaveLen(2))
			Expect(results).To(ContainElement(api.Result{Points: 1, UserId: uuid.UUID(userA.ID())}))
			Expect(results).To(ContainElement(api.Result{Points: 0, UserId: uuid.UUID(userB.ID())}))
		})
	})
})
