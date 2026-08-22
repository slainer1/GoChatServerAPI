package main

import (
	"log"
	"main/internal/store/cache"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
)

// go test
func TestGetUser(t *testing.T) {
	withRedis := config{
		redisCfg: redisConfig{
			enabled: true,
		},
	}
	app := newTestApplication(t, withRedis)
	mux := app.mount()
	testToken, err := app.authenticator.GenerateToken(nil)
	if err != nil {
		t.Fatal(err)
	}
	////
	//Unauthenticated request
	///
	t.Run("should not allow unauthenticated requests", func(t *testing.T) {
		//check for the 401 code
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	////
	//Authenticated requests
	///
	t.Run("should allow authenticated requests", func(t *testing.T) {
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		mockCacheStore.On("Get", mock.Anything, int64(42)).Return(nil, nil)
		mockCacheStore.On("Get", mock.Anything, int64(1)).Return(nil, nil)
		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)
		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testToken)

		rr := executeRequest(req, mux)

		checkResponseCode(t, http.StatusOK, rr.Code)
		log.Println(rr.Body)
		mockCacheStore.Calls = nil
	})
	//spies using testify
	t.Run("should hit the cache first and if not exists it sets the user on the cache", func(t *testing.T) {
		withRedis := config{
			redisCfg: redisConfig{
				enabled: true,
			},
		}
		app := newTestApplication(t, withRedis)
		mux := app.mount()
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		mockCacheStore.On("Get", mock.Anything, int64(42)).Return(nil, nil)
		mockCacheStore.On("Get", mock.Anything, int64(1)).Return(nil, nil)

		mockCacheStore.On("Set", mock.Anything, mock.Anything).Return(nil)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
		mockCacheStore.AssertNumberOfCalls(t, "Get", 2)

		mockCacheStore.Calls = nil
	})
	t.Run("should NOT hit the cache if it is not enabled", func(t *testing.T) {
		withRedis := config{
			redisCfg: redisConfig{
				enabled: false,
			},
		}
		app := newTestApplication(t, withRedis)
		mux := app.mount()
		mockCacheStore := app.cacheStorage.Users.(*cache.MockUserStore)

		req, err := http.NewRequest(http.MethodGet, "/v1/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Authorization", "Bearer "+testToken)
		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
		mockCacheStore.AssertNotCalled(t, "Get")
		mockCacheStore.Calls = nil
	})
}
