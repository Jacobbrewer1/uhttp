package uhttp

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpError_Error(t *testing.T) {
	type fields struct {
		error error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Golden Path",
			fields: fields{
				error: errors.New("test error"),
			},
			want: "test error",
		},
		{
			name: "No Error",
			fields: fields{
				error: nil,
			},
			want: "Bad Request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewHTTPError(http.StatusBadRequest, tt.fields.error)
			got := e.Error()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestHttpError_StatusCode(t *testing.T) {
	type fields struct {
		Status int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Golden Path",
			fields: fields{
				Status: http.StatusNotFound,
			},
			want: http.StatusNotFound,
		},
		{
			name: "No Status",
			fields: fields{
				Status: 0,
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewHTTPError(tt.fields.Status, nil)
			got := e.StatusCode()
			require.Equal(t, tt.want, got)
		})
	}
}

func TestHttpError_SetRequestId(t *testing.T) {
	type fields struct {
		RequestID string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Golden Path",
			fields: fields{
				RequestID: "test-id",
			},
			want: "test-id",
		},
		{
			name: "No Request ID",
			fields: fields{
				RequestID: "",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewHTTPError(http.StatusBadRequest, nil)
			e.SetRequestId(tt.fields.RequestID)
			require.Equal(t, tt.want, e.RequestId)
		})
	}
}
