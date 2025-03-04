package uhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacobbrewer1/uhttp/common"
	"github.com/stretchr/testify/require"
)

func TestMustSendMessageWithStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		message string
	}{
		{
			name:    "success",
			status:  http.StatusOK,
			message: "hello",
		},
		{
			name:    "error",
			status:  http.StatusInternalServerError,
			message: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			MustSendMessageWithStatus(w, tt.status, tt.message)

			res := w.Result()
			t.Cleanup(func() {
				require.NoError(t, res.Body.Close())
			})

			require.Equal(t, tt.status, res.StatusCode)
			var msg common.Message
			err := DecodeJSON(res.Body, &msg)
			require.NoError(t, err)
			require.Equal(t, tt.message, msg.Message)
		})
	}
}

func TestSendMessageWithStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  int
		message string
		wantErr bool
	}{
		{
			name:    "success",
			status:  http.StatusOK,
			message: "hello",
			wantErr: false,
		},
		{
			name:    "error",
			status:  http.StatusInternalServerError,
			message: "error",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := SendMessageWithStatus(w, tt.status, tt.message)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				res := w.Result()
				t.Cleanup(func() {
					require.NoError(t, res.Body.Close())
				})

				require.Equal(t, tt.status, res.StatusCode)
				var msg common.Message
				err := DecodeJSON(res.Body, &msg)
				require.NoError(t, err)
				require.Equal(t, tt.message, msg.Message)
			}
		})
	}
}

func TestMustSendMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "success",
			message: "hello",
		},
		{
			name:    "error",
			message: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			MustSendMessage(w, tt.message)

			res := w.Result()
			t.Cleanup(func() {
				require.NoError(t, res.Body.Close())
			})

			require.Equal(t, http.StatusOK, res.StatusCode)
			var msg common.Message
			err := DecodeJSON(res.Body, &msg)
			require.NoError(t, err)
			require.Equal(t, tt.message, msg.Message)
		})
	}
}

func TestSendMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "success",
			message: "hello",
			wantErr: false,
		},
		{
			name:    "error",
			message: "error",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := SendMessage(w, tt.message)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				res := w.Result()
				t.Cleanup(func() {
					require.NoError(t, res.Body.Close())
				})

				require.Equal(t, http.StatusOK, res.StatusCode)
				var msg common.Message
				err := DecodeJSON(res.Body, &msg)
				require.NoError(t, err)
				require.Equal(t, tt.message, msg.Message)
			}
		})
	}
}
