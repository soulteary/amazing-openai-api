package yi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/soulteary/amazing-openai-api/internal/define"
	"github.com/soulteary/amazing-openai-api/models/yi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
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

// Tests
func TestProxyWithConverter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockReqConverter := new(MockedRequestConverter)
	r.Use(yi.ProxyWithConverter(mockReqConverter))

	req, _ := http.NewRequest(http.MethodOptions, "/", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGetDeploymentByModel(t *testing.T) {
	// Assuming ModelConfig has been defined with at least one key "test-model"
	modelName := "test-model"
	expectedConfig := &define.ModelConfig{
		Key: "some-key",
		// ... other fields
	}

	// Set up the global variable for testing
	yi.ModelConfig = map[string]define.ModelConfig{
		modelName: *expectedConfig,
	}

	config, err := yi.GetDeploymentByModel(modelName)
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)

	// Test with a non-existing model
	_, err = yi.GetDeploymentByModel("non-existing-model")
	assert.Error(t, err)
}
