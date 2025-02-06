package uhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNotFoundHandler(t *testing.T) {
	tests := []struct {
		name   string
		w      *httptest.ResponseRecorder
		r      *http.Request
		status int
		want   string
	}{
		{
			name:   "NotFound",
			w:      httptest.NewRecorder(),
			r:      httptest.NewRequest(http.MethodGet, "/", http.NoBody),
			status: http.StatusNotFound,
			want:   "{\"detail\":\"not found\",\"details\":[\"method: GET\",\"path: /\"],\"request_id\":\"\",\"status\":404,\"title\":\"Not Found\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NotFoundHandler().ServeHTTP(tt.w, tt.r)
			require.Equal(t, tt.status, tt.w.Code)
			require.Equal(t, tt.want, tt.w.Body.String())
		})
	}
}
