package handlers

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"

// 	gomock "go.uber.org/mock/gomock"

// 	contextkeys "github.com/meshyampratap01/letStayInn/internal/contextKeys"
// 	"github.com/meshyampratap01/letStayInn/internal/mocks"
// )

// func makeRequest(method, path string, body []byte, role, userID interface{}) *http.Request {
// 	req := httptest.NewRequest(method, path, bytes.NewReader(body))
// 	ctx := req.Context()
// 	if role != nil {
// 		ctx = context.WithValue(ctx, contextkeys.UserRoleKey, role)
// 	}
// 	if userID != nil {
// 		ctx = context.WithValue(ctx, contextkeys.UserIDKey, userID)
// 	}
// 	return req.WithContext(ctx)
// }

// func TestSubmitFeedbackHTTP(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockSvc := mocks.NewMockIFeedbackService(ctrl)
// 	handler := NewFeedbackHandler(mockSvc)

// 	tests := []struct {
// 		name          string
// 		req           *http.Request
// 		mockExpect    func()
// 		wantStatus    int
// 		wantBodyMatch string
// 	}{
// 		{
// 			name:       "invalid role in context",
// 			req:        makeRequest(http.MethodPost, "/api/v1/feedback", nil, 123, "user-1"),
// 			mockExpect: func() {},
// 			wantStatus: http.StatusUnauthorized,
// 			wantBodyMatch: "invalid role",
// 		},
// 		{
// 			name:       "forbidden role (not Guest)",
// 			req:        makeRequest(http.MethodPost, "/api/v1/feedback", nil, "Manager", "user-1"),
// 			mockExpect: func() {},
// 			wantStatus: http.StatusForbidden,
// 			wantBodyMatch: "guest access required",
// 		},
// 		{
// 			name:       "invalid userID type",
// 			req:        makeRequest(http.MethodPost, "/api/v1/feedback", nil, "Guest", 123),
// 			mockExpect: func() {},
// 			wantStatus: http.StatusUnauthorized,
// 			wantBodyMatch: "invalid user ID",
// 		},
// 		{
// 			name:       "invalid JSON body",
// 			req:        makeRequest(http.MethodPost, "/api/v1/feedback", []byte("{bad json"), "Guest", "user-1"),
// 			mockExpect: func() {},
// 			wantStatus: http.StatusBadRequest,
// 			wantBodyMatch: "Invalid request body",
// 		},
// 		{
// 			name: "service error while submitting feedback",
// 			req: makeRequest(
// 				http.MethodPost,
// 				"/api/v1/feedback",
// 				[]byte(`{"message":"great","rating":5}`),
// 				"Guest",
// 				"user-1",
// 			),
// 			mockExpect: func() {
// 				mockSvc.EXPECT().
// 					SubmitFeedback(gomock.Any(), "great", 5).
// 					Return(errors.New("db error"))
// 			},
// 			wantStatus:    http.StatusInternalServerError,
// 			wantBodyMatch: "db error",
// 		},
// 		{
// 			name: "success",
// 			req: makeRequest(
// 				http.MethodPost,
// 				"/api/v1/feedback",
// 				[]byte(`{"message":"nice stay","rating":4}`),
// 				"Guest",
// 				"user-1",
// 			),
// 			mockExpect: func() {
// 				mockSvc.EXPECT().
// 					SubmitFeedback(gomock.Any(), "nice stay", 4).
// 					Return(nil)
// 			},
// 			wantStatus:    http.StatusCreated,
// 			wantBodyMatch: "Feedback submitted",
// 		},
// 		{
// 			name: "success with empty message (edge case)",
// 			req: makeRequest(
// 				http.MethodPost,
// 				"/api/v1/feedback",
// 				[]byte(`{"message":"","rating":0}`),
// 				"Guest",
// 				"user-1",
// 			),
// 			mockExpect: func() {
// 				mockSvc.EXPECT().
// 					SubmitFeedback(gomock.Any(), "", 0).
// 					Return(nil)
// 			},
// 			wantStatus:    http.StatusCreated,
// 			wantBodyMatch: "Feedback submitted",
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			tc.mockExpect()
// 			rec := httptest.NewRecorder()
// 			handler.SubmitFeedbackHTTP(rec, tc.req)

// 			if rec.Code != tc.wantStatus {
// 				t.Errorf("[%s] expected status %d, got %d", tc.name, tc.wantStatus, rec.Code)
// 			}
// 			// Check Content-Type header is always application/json
// 			if got := rec.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
// 				t.Errorf("[%s] expected Content-Type application/json, got %s", tc.name, got)
// 			}
// 			if tc.wantBodyMatch != "" {
// 				if !strings.Contains(rec.Body.String(), tc.wantBodyMatch) {
// 					t.Errorf("[%s] expected body to contain %q, got %q", tc.name, tc.wantBodyMatch, rec.Body.String())
// 				}
// 			}
// 			// Also validate response is valid JSON
// 			var js map[string]interface{}
// 			if err := json.Unmarshal(rec.Body.Bytes(), &js); err != nil {
// 				t.Errorf("[%s] response is not valid JSON: %v", tc.name, err)
// 			}
// 		})
// 	}
// }
