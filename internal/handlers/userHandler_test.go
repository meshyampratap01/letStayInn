package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshyampratap01/letStayInn/internal/mocks"
	"github.com/meshyampratap01/letStayInn/internal/models"
	gomock "go.uber.org/mock/gomock"
)

func TestSignupHTTPHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIUserService(ctrl)
	h := NewUserHandler(mockService)

	tests := []struct {
		name           string
		body           string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:           "invalid JSON body",
			body:           "{invalid-json}",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid email",
			body:           `{"name":"John","email":"bademail","password":"Password@123"}`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid password",
			body:           `{"name":"John","email":"john@example.com","password":"123"}`,
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "signup service error",
			body: `{"name":"John","email":"john@example.com","password":"Password@123"}`,
			mockSetup: func() {
				mockService.EXPECT().
					Signup("John", "john@example.com", "Password@123", int(models.RoleGuest)).
					Return("", errors.New("signup failed"))
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "signup success",
			body: `{"name":"John","email":"john@example.com","password":"Password@123"}`,
			mockSetup: func() {
				mockService.EXPECT().
					Signup("John", "john@example.com", "Password@123", int(models.RoleGuest)).
					Return("signup successful", nil)
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewBufferString(tt.body))

			tt.mockSetup()
			h.SignupHTTPHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func TestLoginHTTPHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIUserService(ctrl)
	h := NewUserHandler(mockService)

	validUser := &models.User{
		ID:   "u1",
		Name: "John",
		Role: models.RoleGuest,
	}

	tests := []struct {
		name           string
		body           string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:           "invalid JSON body",
			body:           "{invalid-json}",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "login service error",
			body: `{"email":"john@example.com","password":"Password@123"}`,
			mockSetup: func() {
				mockService.EXPECT().
					Login("john@example.com", "Password@123").
					Return(nil, errors.New("invalid credentials"))
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "login success",
			body: `{"email":"john@example.com","password":"Password@123"}`,
			mockSetup: func() {
				mockService.EXPECT().
					Login("john@example.com", "Password@123").
					Return(validUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(tt.body))

			tt.mockSetup()
			h.LoginHTTPHandler(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
