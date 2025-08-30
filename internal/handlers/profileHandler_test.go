package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/dto"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	gomock "go.uber.org/mock/gomock"
)

func setupProfileHandler(t *testing.T) (*ProfileHandler, *mocks.MockIUserService, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockUserService := mocks.NewMockIUserService(ctrl)
	h := NewProfileHandler(mockUserService)
	return h, mockUserService, ctrl
}

func TestViewProfileHTTP(t *testing.T) {
	t.Run("missing userID in context", func(t *testing.T) {
		h, _, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil)
		rr := httptest.NewRecorder()

		h.ViewProfileHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rr.Code)
		}
	})

	t.Run("GetUserByID error", func(t *testing.T) {
		h, mockUserService, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil).
			WithContext(context.WithValue(context.Background(), contextkeys.UserIDKey, "123"))
		rr := httptest.NewRecorder()

		mockUserService.EXPECT().GetUserByID("123").Return(nil, errors.New("db error"))

		h.ViewProfileHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		h, mockUserService, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		user := &models.User{ID: "123", Name: "Alice", Email: "alice@test.com", Role: models.RoleManager, Available: true}

		req := httptest.NewRequest(http.MethodGet, "/api/v1/profile", nil).
			WithContext(context.WithValue(context.Background(), contextkeys.UserIDKey, "123"))
		rr := httptest.NewRecorder()

		mockUserService.EXPECT().GetUserByID("123").Return(user, nil)

		h.ViewProfileHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rr.Code)
		}

		var wrapper struct {
			Code    int                `json:"code"`
			Message string             `json:"message"`
			Data    dto.UserProfileDTO `json:"data"`
		}
		if err := json.NewDecoder(rr.Body).Decode(&wrapper); err != nil {
			t.Fatal("failed to decode response")
		}
		if wrapper.Data.ID != "123" {
			t.Errorf("expected ID=123, got %s", wrapper.Data.ID)
		}
	})
}

func TestUpdateProfileHTTP(t *testing.T) {
	t.Run("missing userID", func(t *testing.T) {
		h, _, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", nil)
		rr := httptest.NewRecorder()

		h.UpdateProfileHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", rr.Code)
		}
	})

	t.Run("GetUserByID error", func(t *testing.T) {
		h, mockUserService, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", nil).
			WithContext(context.WithValue(context.Background(), contextkeys.UserIDKey, "123"))
		rr := httptest.NewRecorder()

		mockUserService.EXPECT().GetUserByID("123").Return(nil, errors.New("db error"))

		h.UpdateProfileHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		h, mockUserService, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewBufferString("{invalid")).
			WithContext(context.WithValue(context.Background(), contextkeys.UserIDKey, "123"))
		rr := httptest.NewRecorder()

		mockUserService.EXPECT().GetUserByID("123").Return(&models.User{ID: "123"}, nil)

		h.UpdateProfileHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("no valid fields", func(t *testing.T) {
		h, mockUserService, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		body := `{"name":"A","email":"abc","password":"123"}`
		req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewBufferString(body)).
			WithContext(context.WithValue(context.Background(), contextkeys.UserIDKey, "123"))
		rr := httptest.NewRecorder()

		mockUserService.EXPECT().GetUserByID("123").Return(&models.User{ID: "123"}, nil)

		h.UpdateProfileHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", rr.Code)
		}
	})

	t.Run("UpdateUser error", func(t *testing.T) {
		h, mockUserService, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		newName := "Alice Updated"
		body, _ := json.Marshal(map[string]string{"name": newName})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewBuffer(body)).
			WithContext(context.WithValue(context.Background(), contextkeys.UserIDKey, "123"))
		rr := httptest.NewRecorder()

		user := &models.User{ID: "123", Name: "Alice"}
		mockUserService.EXPECT().GetUserByID("123").Return(user, nil)
		mockUserService.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("update failed"))

		h.UpdateProfileHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected 500, got %d", rr.Code)
		}
	})

	t.Run("success with name + email + password", func(t *testing.T) {
		h, mockUserService, ctrl := setupProfileHandler(t)
		defer ctrl.Finish()

		newName := "Alice Updated"
		newEmail := "alice@updated.com"
		newPassword := "securepass"
		body, _ := json.Marshal(map[string]string{
			"name":     newName,
			"email":    newEmail,
			"password": newPassword,
		})
		req := httptest.NewRequest(http.MethodPut, "/api/v1/profile", bytes.NewBuffer(body)).
			WithContext(context.WithValue(context.Background(), contextkeys.UserIDKey, "123"))
		rr := httptest.NewRecorder()

		user := &models.User{ID: "123", Name: "Alice", Email: "alice@test.com"}
		mockUserService.EXPECT().GetUserByID("123").Return(user, nil)
		mockUserService.EXPECT().UpdateUser(gomock.Any()).Return(nil)

		h.UpdateProfileHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	})
}

func TestContainsAt(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"user@example.com", true},
		{"foo@bar", true},
		{"foobar.com", false},
		{"", false},
		{"@leading", true},
		{"trailing@", true},
	}

	for _, tt := range tests {
		got := containsAt(tt.email)
		if got != tt.expected {
			t.Errorf("containsAt(%q) = %v, want %v", tt.email, got, tt.expected)
		}
	}
}
