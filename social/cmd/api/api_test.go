package main

import (
	"main/internal/ratelimiter"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config{
		rateLimiter: ratelimiter.Config{
			Enabled:              true,
			RequestsPerTimeFrame: 20,
			TimeFrame:            5 * time.Second,
		},
		addr: ":8080",
	}
	app := newTestApplication(t, cfg)

	//set rate limiter used
	if app.rateLimiter == nil {
		app.rateLimiter = ratelimiter.NewFixedWindowRateLimiter(20, 5*time.Second)
	}
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	mockIP := "192.168.1.1"
	marginOfError := 2

	for i := 0; i < cfg.rateLimiter.RequestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest("GET", ts.URL+"/v1/health", nil)
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}
		req.Header.Set("X-Real-IP", mockIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not make request: %v", err)
		}
		defer resp.Body.Close()

		if i < cfg.rateLimiter.RequestsPerTimeFrame {
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected Status 200 but got %v", resp.StatusCode)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("Expected Status too many requests; got %v", resp.StatusCode)
			}
		}
		resp.Body.Close()
	}

}
