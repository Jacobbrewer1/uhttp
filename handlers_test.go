package uhttp

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotFoundHandler(t *testing.T) {
	tests := []struct {
		name      string
		w         *httptest.ResponseRecorder
		r         *http.Request
		requestId string
		status    int
		want      string
	}{
		{
			name:      "NotFound",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId: "",
			status:    http.StatusNotFound,
			want:      `{"detail":"not found","details":["method: GET","path: /"],"request_id":"","status":404,"title":"Not Found"}`,
		},
		{
			name:      "NotFoundWithQuery",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/?foo=bar", http.NoBody),
			requestId: "",
			status:    http.StatusNotFound,
			want:      `{"detail":"not found","details":["method: GET","path: /","query: foo=bar"],"request_id":"","status":404,"title":"Not Found"}`,
		},
		{
			name:      "NotFoundWithRequestID",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId: "123",
			status:    http.StatusNotFound,
			want:      `{"detail":"not found","details":["method: GET","path: /"],"request_id":"123","status":404,"title":"Not Found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requestId != "" {
				tt.r.Header.Set(requestIDHeader, tt.requestId)
				ctx := RequestIDToContext(tt.r.Context(), tt.r)
				tt.r = tt.r.WithContext(ctx)
			}

			NotFoundHandler().ServeHTTP(tt.w, tt.r)
			require.Equal(t, tt.status, tt.w.Code)

			resp := new(HTTPError)
			require.NoError(t, DecodeJSON(tt.w.Result().Body, resp))

			require.Equal(t, "not found", resp.Detail)

			expDetails := []any{"method: " + tt.r.Method, "path: " + tt.r.URL.Path}
			if tt.r.URL.RawQuery != "" {
				expDetails = append(expDetails, "query: "+tt.r.URL.RawQuery)
			}
			require.Equal(t, expDetails, resp.Details)

			require.Equal(t, tt.status, resp.Status)
			require.Equal(t, "Not Found", resp.Title)

			if tt.requestId != "" {
				require.Equal(t, tt.requestId, resp.RequestId)
			} else {
				require.NotEmpty(t, resp.RequestId)
			}
		})
	}
}

func TestMethodNotAllowedHandler(t *testing.T) {
	tests := []struct {
		name      string
		w         *httptest.ResponseRecorder
		r         *http.Request
		requestId string
		status    int
		want      string
	}{
		{
			name:      "MethodNotAllowed",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId: "",
			status:    http.StatusMethodNotAllowed,
			want:      `{"detail":"method not allowed","details":["method: GET","path: /"],"request_id":"","status":405,"title":"Method Not Allowed"}`,
		},
		{
			name:      "MethodNotAllowedWithQuery",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/?foo=bar", http.NoBody),
			requestId: "",
			status:    http.StatusMethodNotAllowed,
			want:      `{"detail":"method not allowed","details":["method: GET","path: /"],"request_id":"","status":405,"title":"Method Not Allowed"}`,
		},
		{
			name:      "MethodNotAllowedWithRequestID",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId: "123",
			status:    http.StatusMethodNotAllowed,
			want:      `{"detail":"method not allowed","details":["method: GET","path: /"],"request_id":"123","status":405,"title":"Method Not Allowed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requestId != "" {
				tt.r.Header.Set(requestIDHeader, tt.requestId)
				ctx := RequestIDToContext(tt.r.Context(), tt.r)
				tt.r = tt.r.WithContext(ctx)
			}

			MethodNotAllowedHandler().ServeHTTP(tt.w, tt.r)
			require.Equal(t, tt.status, tt.w.Code)

			resp := new(HTTPError)
			require.NoError(t, DecodeJSON(tt.w.Result().Body, resp))

			require.Equal(t, "method not allowed", resp.Detail)

			expDetails := []any{"method: " + tt.r.Method, "path: " + tt.r.URL.Path}
			if tt.r.URL.RawQuery != "" {
				expDetails = append(expDetails, "query: "+tt.r.URL.RawQuery)
			}
			require.Equal(t, expDetails, resp.Details)

			require.Equal(t, tt.status, resp.Status)
			require.Equal(t, "Method Not Allowed", resp.Title)

			if tt.requestId != "" {
				require.Equal(t, tt.requestId, resp.RequestId)
			} else {
				require.NotEmpty(t, resp.RequestId)
			}
		})
	}
}

func TestUnauthorizedHandler(t *testing.T) {
	tests := []struct {
		name      string
		w         *httptest.ResponseRecorder
		r         *http.Request
		requestId string
		status    int
		want      string
	}{
		{
			name:      "Unauthorized",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId: "",
			status:    http.StatusUnauthorized,
			want:      `{"detail":"unauthorized","details":["method: GET","path: /"],"request_id":"","status":401,"title":"Unauthorized"}`,
		},
		{
			name:      "UnauthorizedWithQuery",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/?foo=bar", http.NoBody),
			requestId: "",
			status:    http.StatusUnauthorized,
			want:      `{"detail":"unauthorized","details":["method: GET","path: /","query: foo=bar"],"request_id":"","status":401,"title":"Unauthorized"}`,
		},
		{
			name:      "UnauthorizedWithRequestID",
			w:         httptest.NewRecorder(),
			r:         httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId: "123",
			status:    http.StatusUnauthorized,
			want:      `{"detail":"unauthorized","details":["method: GET","path: /"],"request_id":"123","status":401,"title":"Unauthorized"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requestId != "" {
				tt.r.Header.Set(requestIDHeader, tt.requestId)
				ctx := RequestIDToContext(tt.r.Context(), tt.r)
				tt.r = tt.r.WithContext(ctx)
			}

			UnauthorizedHandler().ServeHTTP(tt.w, tt.r)
			require.Equal(t, tt.status, tt.w.Code)

			resp := new(HTTPError)
			require.NoError(t, DecodeJSON(tt.w.Result().Body, resp))

			require.Equal(t, "unauthorized", resp.Detail)

			expDetails := []any{"method: " + tt.r.Method, "path: " + tt.r.URL.Path}
			if tt.r.URL.RawQuery != "" {
				expDetails = append(expDetails, "query: "+tt.r.URL.RawQuery)
			}
			require.Equal(t, expDetails, resp.Details)

			require.Equal(t, tt.status, resp.Status)
			require.Equal(t, "Unauthorized", resp.Title)

			if tt.requestId != "" {
				require.Equal(t, tt.requestId, resp.RequestId)
			} else {
				require.NotEmpty(t, resp.RequestId)
			}
		})
	}
}

func TestGenericErrorHandler(t *testing.T) {
	tests := []struct {
		name       string
		w          *httptest.ResponseRecorder
		r          *http.Request
		requestId  string
		inputError error
		status     int
		want       string
	}{
		{
			name:       "GenericError",
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId:  "",
			inputError: errors.New("some error"),
			status:     http.StatusBadRequest,
			want:       `{"detail":"some error","details":["method: GET","path: /"],"request_id":"","status":400,"title":"Bad Request"}`,
		},
		{
			name:       "GenericErrorWithQuery",
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest(http.MethodGet, "/?foo=bar", http.NoBody),
			requestId:  "",
			inputError: errors.New("some unknown error"),
			status:     http.StatusBadRequest,
			want:       `{"detail":"internal server error","details":["method: GET","path: /","query: foo=bar"],"request_id":"","status":500,"title":"Internal Server Error"}`,
		},
		{
			name:       "GenericErrorWithRequestID",
			w:          httptest.NewRecorder(),
			r:          httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			requestId:  "123",
			inputError: errors.New("some error test with a request ID"),
			status:     http.StatusBadRequest,
			want:       `{"detail":"internal server error","details":["method: GET","path: /"],"request_id":"123","status":500,"title":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.requestId != "" {
				tt.r.Header.Set(requestIDHeader, tt.requestId)
				ctx := RequestIDToContext(tt.r.Context(), tt.r)
				tt.r = tt.r.WithContext(ctx)
			}

			GenericErrorHandler(tt.w, tt.r, tt.inputError)
			require.Equal(t, tt.status, tt.w.Code)

			resp := new(HTTPError)
			require.NoError(t, DecodeJSON(tt.w.Result().Body, resp))

			require.Equal(t, tt.inputError.Error(), resp.Detail)
			require.Nil(t, resp.Details)

			require.Equal(t, tt.status, resp.Status)
			require.Equal(t, "Bad Request", resp.Title)

			if tt.requestId != "" {
				require.Equal(t, tt.requestId, resp.RequestId)
			} else {
				require.NotEmpty(t, resp.RequestId)
			}
		})
	}
}
