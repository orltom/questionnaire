package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"uuid"

	"gitlab.com/orltom/questionnaire/backend/api"
)

type TestClient struct {
	*http.Client
	baseURL string
	userID  uuid.UUID
}

func NewTestClient(baseURL string, userID uuid.UUID) *TestClient {
	return &TestClient{
		Client:  &http.Client{},
		baseURL: baseURL,
		userID:  userID,
	}
}

func (tc *TestClient) AsUser(userID uuid.UUID) *TestClient {
	return &TestClient{
		Client:  tc.Client,
		baseURL: tc.baseURL,
		userID:  userID,
	}
}

func (tc *TestClient) CreateQuiz(title, description string, visibility string) (uuid.UUID, error) {
	body := api.CreateQuiz{
		Title:       title,
		Description: description,
		Visibility:  api.Visibility(visibility),
	}
	resp, err := tc.post("/quizzes", body)
	if err != nil {
		return uuid.UUID{}, err
	}

	var quiz api.Quiz
	if err := json.Unmarshal(resp, &quiz); err != nil {
		return uuid.UUID{}, err
	}

	return quiz.Id, nil
}

func (tc *TestClient) GetQuiz(quizID uuid.UUID) (*api.Quiz, error) {
	resp, err := tc.get(fmt.Sprintf("/quizzes/%s", quizID))
	if err != nil {
		return nil, err
	}

	var quiz api.Quiz
	if err := json.Unmarshal(resp, &quiz); err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (tc *TestClient) CreateQuestion(description, visibility string, answers []api.Answer) (uuid.UUID, error) {
	body := api.CreateQuestion{
		Description: description,
		Visibility:  api.Visibility(visibility),
		Answers:     answers,
	}
	resp, err := tc.post("/questions", body)
	if err != nil {
		return uuid.UUID{}, err
	}

	var question api.Question
	if err := json.Unmarshal(resp, &question); err != nil {
		return uuid.UUID{}, err
	}

	return question.Id, nil
}

func (tc *TestClient) GetQuestion(questionID uuid.UUID) (*api.Question, error) {
	resp, err := tc.get(fmt.Sprintf("/questions/%s", questionID))
	if err != nil {
		return nil, err
	}

	var question api.Question
	if err := json.Unmarshal(resp, &question); err != nil {
		return nil, err
	}

	return &question, nil
}

func (tc *TestClient) AddQuestionToQuiz(quizID, questionID uuid.UUID, position int) error {
	body := api.AddQuestionToQuiz{
		QuestionId: questionID,
		Position:   position,
	}
	_, err := tc.post(fmt.Sprintf("/quizzes/%s/questions", quizID), body)
	return err
}

func (tc *TestClient) CreateChallenge(quizID uuid.UUID, startsAt time.Time, access string, invitees []uuid.UUID) (uuid.UUID, error) {
	apiInvitees := make([]api.QuizId, len(invitees))
	copy(apiInvitees, invitees)

	body := api.CreateChallenge{
		QuizId:   quizID,
		StartsAt: startsAt,
		Access:   api.Access(access),
	}
	if invitees != nil {
		body.Invitees = &apiInvitees
	}

	resp, err := tc.post("/challenges", body)
	if err != nil {
		return uuid.UUID{}, err
	}

	var challenge api.Challenge
	if err := json.Unmarshal(resp, &challenge); err != nil {
		return uuid.UUID{}, err
	}

	return challenge.Id, nil
}

func (tc *TestClient) SubmitAnswer(challengeID, questionID, answerID uuid.UUID) error {
	body := api.SubmitAnswer{
		QuestionId: questionID,
		AnswerId:   answerID,
	}
	_, err := tc.post(fmt.Sprintf("/challenges/%s/answers", challengeID), body)
	return err
}

func (tc *TestClient) GetChallengeResults(challengeID uuid.UUID) ([]api.Result, error) {
	resp, err := tc.get(fmt.Sprintf("/challenges/%s/results", challengeID))
	if err != nil {
		return nil, err
	}

	var results []api.Result
	if err := json.Unmarshal(resp, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (tc *TestClient) post(path string, body any) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, tc.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Test-User-ID", tc.userID.String())

	resp, err := tc.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (tc *TestClient) get(path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tc.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Test-User-ID", tc.userID.String())

	resp, err := tc.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}
