package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gomock "go.uber.org/mock/gomock"

	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
	"github.com/meshyampratap01/letStayInn/internal/mocks"
)

// helper to make request with body
func makeRequest(method, path string, body []byte, role, userID interface{}) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	ctx := req.Context()
	if role != nil {
		ctx = context.WithValue(ctx, contextkeys.UserRoleKey, role)
	}
	if userID != nil {
		ctx = context.WithValue(ctx, contextkeys.UserIDKey, userID)
	}
	return req.WithContext(ctx)
}

func TestSubmitFeedbackHTTP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := mocks.NewMockIFeedbackService(ctrl)
	handler := NewFeedbackHandler(mockSvc)

	tests := []struct {
		name          string
		req           *http.Request
		mockExpect    func()
		wantStatus    int
		wantBodyMatch string
	}{
		{
			name:       "invalid role in context",
			req:        makeRequest(http.MethodPost, "/api/v1/feedback", nil, 123, "user-1"),
			mockExpect: func() {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "forbidden role (not Guest)",
			req:        makeRequest(http.MethodPost, "/api/v1/feedback", nil, "Manager", "user-1"),
			mockExpect: func() {},
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "invalid userID type",
			req:        makeRequest(http.MethodPost, "/api/v1/feedback", nil, "Guest", 123),
			mockExpect: func() {},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid JSON body",
			req:        makeRequest(http.MethodPost, "/api/v1/feedback", []byte("{bad json"), "Guest", "user-1"),
			mockExpect: func() {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "service error while submitting feedback",
			req: makeRequest(
				http.MethodPost,
				"/api/v1/feedback",
				[]byte(`{"message":"great","rating":5}`),
				"Guest",
				"user-1",
			),
			mockExpect: func() {
				mockSvc.EXPECT().SubmitFeedback(gomock.Any(), "user-1", 5).Return(errors.New("db error"))
			},
			wantStatus:    http.StatusInternalServerError,
			wantBodyMatch: "db error",
		},
		{
			name: "success",
			req: makeRequest(
				http.MethodPost,
				"/api/v1/feedback",
				[]byte(`{"message":"great","rating":5}`),
				"Guest",
				"user-1",
			),
			mockExpect: func() {
				mockSvc.EXPECT().SubmitFeedback(gomock.Any(), "user-1", 5).Return(nil)
			},
			wantStatus:    http.StatusCreated,
			wantBodyMatch: "Feedback submitted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			tt.mockExpect()

			handler.SubmitFeedbackHTTP(rr, tt.req)

			if rr.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}
			if tt.wantBodyMatch != "" && !strings.Contains(rr.Body.String(), tt.wantBodyMatch) {
				t.Errorf("expected body to contain %q, got %q", tt.wantBodyMatch, rr.Body.String())
			}
		})
	}
}
