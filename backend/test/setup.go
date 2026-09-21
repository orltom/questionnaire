package test

import (
	"net/http/httptest"
	"uuid"

	//nolint:staticcheck
	. "github.com/onsi/ginkgo/v2"

	"gitlab.com/orltom/questionnaire/backend/internal/server"
)

func StartTestServer() *TestClient {
	deps := server.SetupServiceLayer()
	handler := server.BuildMux(deps)

	ts := httptest.NewServer(handler)
	DeferCleanup(ts.Close)

	return NewTestClient(ts.URL, uuid.New())
}
