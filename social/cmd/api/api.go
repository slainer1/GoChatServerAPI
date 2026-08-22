package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"main/docs"
	"main/internal/auth"
	"main/internal/env"
	"main/internal/mailer"
	"main/internal/ratelimiter"
	"main/internal/store"
	"main/internal/store/cache"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config        config
	store         store.Storage
	cacheStorage  cache.Storage
	logger        *zap.SugaredLogger
	mailer        mailer.Client
	authenticator auth.Authenticator
	rateLimiter   ratelimiter.Limiter
}
type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	mail        mailConfig
	frontendURL string
	auth        authConfig
	redisCfg    redisConfig
	rateLimiter ratelimiter.Config
}
type redisConfig struct {
	addr     string
	password string
	db       int
	enabled  bool
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}
type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}
type basicConfig struct {
	user string
	pass string
}
type mailConfig struct {
	sendGrid  sendGridConfig
	fromEmail string
	exp       time.Duration
}
type sendGridConfig struct {
	apiKey string
}

// Handler for middleware
func (app *application) mount() *chi.Mux {
	//mux := http.NewServeMux()
	//mux.HandleFunc("GET /v1/health", app.healthCheckHandler)
	//CHI Library for version control, timeout management, and routing
	//posts, users, auth endpoints
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://localhost:5174")},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(app.rateLimiterMiddleware)

	//HOT RELOADINg

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		//expvar exposes metrics/operations
		r.With(app.BasicAuthMiddleware()).Get("/debug/vars", expvar.Handler().ServeHTTP)

		docsURL := fmt.Sprintf("%s/v1/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		//r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.AuthTokenMiddleware)
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.checkPostOwnership("admin", app.deletePostHandler))
				r.Patch("/", app.checkPostOwnership("moderator", app.updatePostHandler))
				r.Post("/comments", app.createCommentPostHandler)
			})

		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)

				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		//public routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
		})
	})
	return r
}

func (app *application) run(mux *chi.Mux) error {
	//Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}
	//graceful shutdown
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		app.logger.Infow("signal caught", "signal", s)
		shutdown <- srv.Shutdown(ctx)
	}()

	//app.logger.Infow("Server Started on Port:%s", app.config.addr)
	app.logger.Infow("Server Started on: ", app.config.addr, "env: ", app.config.env)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	err = <-shutdown
	if err != nil {
		return err
	}
	app.logger.Infow("Server Stopped on: ", app.config.addr)
	return nil
}
