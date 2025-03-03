package uhttp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jacobbrewer1/uhttp/common"
	"github.com/stretchr/testify/require"
)

func TestMustEncode(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name       string
		status     int
		v          response
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			status:     http.StatusOK,
			v:          response{Message: "hello"},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello"}` + "\n",
		},
		{
			name:       "error",
			status:     http.StatusInternalServerError,
			v:          response{Message: "error"},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"error"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			MustEncode(w, tt.status, tt.v)

			res := w.Result()
			t.Cleanup(func() {
				require.NoError(t, res.Body.Close())
			})

			require.Equal(t, tt.wantStatus, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestEncode(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name       string
		status     int
		v          response
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			status:     http.StatusOK,
			v:          response{Message: "hello"},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello"}` + "\n",
		},
		{
			name:       "error",
			status:     http.StatusInternalServerError,
			v:          response{Message: "error"},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"error"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := Encode(w, tt.status, tt.v)
			require.NoError(t, err)

			res := w.Result()
			t.Cleanup(func() {
				require.NoError(t, res.Body.Close())
			})

			require.Equal(t, tt.wantStatus, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestEncodeJSON(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name       string
		status     int
		v          response
		wantStatus int
		wantBody   string
	}{
		{
			name:       "success",
			status:     http.StatusOK,
			v:          response{Message: "hello"},
			wantStatus: http.StatusOK,
			wantBody:   `{"message":"hello"}` + "\n",
		},
		{
			name:       "error",
			status:     http.StatusInternalServerError,
			v:          response{Message: "error"},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"message":"error"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := EncodeJSON(w, tt.status, tt.v)
			require.NoError(t, err)

			res := w.Result()
			t.Cleanup(func() {
				require.NoError(t, res.Body.Close())
			})

			require.Equal(t, tt.wantStatus, res.StatusCode)
			body, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			require.Equal(t, tt.wantBody, string(body))
		})
	}
}

func TestDecodeJSONBody(t *testing.T) {
	type args struct {
		r *http.Request
	}
	type testCase struct {
		name    string
		args    args
		want    any
		wantErr error
	}
	tests := []testCase{
		{
			name: "success",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"message": "hello"}`)),
			},
			want: &common.Message{Message: "hello"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := new(common.Message)
			err := DecodeRequestJSON(tt.args.r, got)
			require.Equal(t, tt.wantErr, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeRequestJSON(t *testing.T) {
	type request struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name    string
		body    string
		want    request
		wantErr bool
	}{
		{
			name:    "success",
			body:    `{"message": "hello"}`,
			want:    request{Message: "hello"},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			body:    `{"message": "hello"`,
			want:    request{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			var got request
			err := DecodeRequestJSON(r, &got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestDecodeJSON(t *testing.T) {
	type request struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name    string
		body    string
		want    request
		wantErr bool
	}{
		{
			name:    "success",
			body:    `{"message": "hello"}`,
			want:    request{Message: "hello"},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			body:    `{"message": "hello"`,
			want:    request{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := io.NopCloser(strings.NewReader(tt.body))
			var got request
			err := DecodeJSON(reader, &got)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
