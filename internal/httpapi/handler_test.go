package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestGetHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// assert.NotNil(t, req)
	if req == nil {
		t.Fatal("Request is nil")
	}

	recorder := httptest.NewRecorder()
	GetHealthHandler(recorder, req)

// 	assert.Equal(t, http.StatusOK, recorder.Code)
// 	assert.Equal(t, "OK", recorder.Body.String())

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", recorder.Code)
	}

	if recorder.Body.String() != "OK" {
		t.Errorf("Expected body OK, got %s", recorder.Body.String())
	}
}