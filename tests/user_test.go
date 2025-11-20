package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/advn1/backend/internal/handlers"
	"github.com/advn1/backend/internal/models/user"
)

func TestGetAllUsers(t *testing.T) {
	w := httptest.NewRecorder()

	handlers.GetAllUsers(w)

	if w.Code != http.StatusOK {
		t.Errorf("GET /users returned status code: %v. wanted %v", w.Code, http.StatusOK)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Response Content-Type is %v. wanted %v",contentType, "application/json")
	}
	
	var users []user.User
	err := json.Unmarshal(w.Body.Bytes(), &users)

	if err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}
}

func TestGetUserById(t *testing.T) {

}