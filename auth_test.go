package uhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAuthHeaderToContext(t *testing.T) {
	request := httptest.NewRequest("GET", "/", http.NoBody)
	request.Header.Set("Authorization", "Bearer token")

	ctx := AuthHeaderToContext(request.Context(), request)

	got := AuthHeaderFromContext(ctx)
	require.Equal(t, "Bearer token", got)

	ctxVal, ok := ctx.Value(authHeaderKey).(string)
	require.True(t, ok)
	require.Equal(t, "Bearer token", ctxVal)
}

func TestAuthHeaderToContextMux(t *testing.T) {
	request := httptest.NewRequest("GET", "/", http.NoBody)
	request.Header.Set("Authorization", "Bearer token")

	var got string
	var ok bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		got, ok = ctx.Value(authHeaderKey).(string)
	})

	mux := AuthHeaderToContextMux()(handler)
	mux.ServeHTTP(httptest.NewRecorder(), request)

	require.True(t, ok)
	require.Equal(t, "Bearer token", got)
}

func TestAuthToContext(t *testing.T) {
	ctx := context.Background()
	ctx = AuthToContext(ctx, "Bearer token")

	got, ok := ctx.Value(authHeaderKey).(string)
	require.True(t, ok)
	require.Equal(t, "Bearer token", got)
}

func TestAuthHeaderFromContext(t *testing.T) {
	ctx := context.Background()
	ctx = AuthToContext(ctx, "Bearer token")

	got := AuthHeaderFromContext(ctx)
	require.Equal(t, "Bearer token", got)
}

func TestAuthHeaderFromContext_IncorrectType(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, authHeaderKey, 42)

	got := AuthHeaderFromContext(ctx)
	require.Equal(t, "", got)
}

func TestIsInternal_internal(t *testing.T) {
	internalRequest := httptest.NewRequest("GET", "/", http.NoBody)
	internalRequest.Header.Set("X-Your-Internal-Header", "true")

	got := IsInternal(internalRequest)
	require.True(t, got)
}

func TestIsInternal_external(t *testing.T) {
	externalRequest := httptest.NewRequest("GET", "/", http.NoBody)
	externalRequest.Header.Set("X-Forwarded-For", "test")

	got := IsInternal(externalRequest)
	require.False(t, got)
}

func TestInternalOnly(t *testing.T) {
	internalRequest := httptest.NewRequest("GET", "/", http.NoBody)
	internalRequest.Header.Set("X-Your-Internal-Header", "true")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	InternalOnly(handler).ServeHTTP(w, internalRequest)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestInternalOnly_NotInternal(t *testing.T) {
	externalRequest := httptest.NewRequest("GET", "/", http.NoBody)
	externalRequest.Header.Set("X-Forwarded-For", "test")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	InternalOnly(handler).ServeHTTP(w, externalRequest)

	require.Equal(t, http.StatusForbidden, w.Code)
}
