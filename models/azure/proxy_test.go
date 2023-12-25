package azure_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/models/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedRequestConverter struct {
	mock.Mock
}

func (m *MockedRequestConverter) Convert(req *http.Request, deployment *define.ModelConfig) (*http.Request, error) {
	args := m.Called(req, deployment)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *MockedRequestConverter) Name() string {
	args := m.Called()
	return args.String(0)
}

func TestProxyMiddlewareWithOptionsMethod(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockReqConverter := new(MockedRequestConverter)
	r.Use(azure.ProxyWithConverter(mockReqConverter))

	req, _ := http.NewRequest(http.MethodOptions, "/", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	// Check for CORS headers here...
}

func TestModelProxySuccess(t *testing.T) {
	// This test would require setting up the expected behavior of sending requests
	// and collecting results, you would mock the external dependencies.
}

func TestModelProxyFailures(t *testing.T) {
	// Similarly, this would test failure scenarios (bad responses, errors in request sending, etc.)
	// by adjusting the mocked behavior accordingly.
}

func TestProxyFunctionality(t *testing.T) {
	// Here you would validate the proxy functionality with a setup similar to
	// 'TestModelProxySuccess' and 'TestModelProxyFailures' tests but focusing on the Proxy function.
}

func TestGetDeploymentByModel(t *testing.T) {
	expectedModel := "test-model"
	expectedConfig := define.ModelConfig{
		Name:     expectedModel,
		Endpoint: "https://example.com",
		Key:      "secret-key",
	}

	// Assuming ModelConfig is a global variable storing configurations, it should be mocked or set appropriately.
	azure.ModelConfig = map[string]define.ModelConfig{
		expectedModel: expectedConfig,
	}

	config, err := azure.GetDeploymentByModel(expectedModel)

	assert.Nil(t, err)
	assert.Equal(t, &expectedConfig, config)
}

func TestGetDeploymentByModelNotFound(t *testing.T) {
	unexpectedModel := "non-existent-model"

	_, err := azure.GetDeploymentByModel(unexpectedModel)

	assert.NotNil(t, err)
	assert.Equal(t, "deployment config for non-existent-model not found", err.Error())
}
