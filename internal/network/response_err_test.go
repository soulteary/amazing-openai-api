package network_test

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/soulteary/amazing-openai-api/internal/network"
)

func TestSendError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	testError := errors.New("internal server error")
	network.SendError(c, testError)

	if w.Code != 500 {
		t.Errorf("Expected status code 500, got %d", w.Code)
	}

	var apiResponse network.ApiResponse
	err := json.Unmarshal(w.Body.Bytes(), &apiResponse)
	if err != nil {
		t.Fatalf("Error unmarshalling response: %v", err)
	}

	if apiResponse.Error.Code != "500" {
		t.Errorf("Expected error code '500', got '%s'", apiResponse.Error.Code)
	}

	expectedErrorMessage := testError.Error()
	if apiResponse.Error.Message != expectedErrorMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedErrorMessage, apiResponse.Error.Message)
	}
}
