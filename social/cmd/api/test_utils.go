package main

import (
	"main/internal/auth"
	"main/internal/store"
	"main/internal/store/cache"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

// ***************************************************
// Testing Utils, abstraction methods for testing, and init mock application
// ***************************************************
func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	//logger := zap.NewNop().Sugar()
	logger := zap.Must(zap.NewProduction()).Sugar()
	mockStore := store.NewMockStore()
	mockCacheStorage := cache.NewMockStore()
	testAuth := &auth.TestAuthenticator{}

	return &application{
		logger:        logger,
		store:         mockStore,
		cacheStorage:  mockCacheStorage,
		authenticator: testAuth,
		config:        cfg,
	}
}
func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code to be %d and we got %d", expected, actual)
	}
}
