package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"main/internal/store"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// the point of this middleware is to recieve a html header (or handler) and process edge cases and check whether conditions are right
// if conditions are right, then we return the same header with context, or the page we are trying to go to (in this case authentication or authorization)
func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//read header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("missing Authorization header"))
			return
		}
		//parse header
		parts := strings.Split(authHeader, " ") // authorization: Bearer <token>
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("Authorization header is malformed"))
			return
		}
		//validate header
		token := parts[1]
		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}
		//parse the claims
		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.getUser(ctx, userID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//read the auth header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("Authorization header is missing"))
				return
			}

			//parse it -> get the base64
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("Authorization header is malformed"))
				return
			}
			//decode it

			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicResponse(w, r, err)
				return
			}

			username := app.config.auth.basic.user
			pass := app.config.auth.basic.pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				app.unauthorizedBasicResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}
			//check the credentials (config)
			next.ServeHTTP(w, r)
		})
	}
}

func (app *application) checkPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.getUserFromContext(r)
		post := getPostFromCtx(r)
		// if it is the user's post, then allow handle
		if post.UserID == user.ID {
			next.ServeHTTP(w, r)
			return
		}
		// role check
		allowed, err := app.checkRolePrecedence(r.Context(), user, role)
		if err != nil {
			app.statusInternalServerErrorHandler(w, r, err)
			return
		}
		if !allowed {
			app.forbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) checkRolePrecedence(ctx context.Context, user *store.User, roleName string) (bool, error) {
	role, err := app.store.Roles.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}
	return user.Role.Level >= role.Level, nil

}

func (app *application) getUser(ctx context.Context, userID int64) (*store.User, error) {
	if !app.config.redisCfg.enabled {
		return app.store.Users.GetByID(ctx, userID)
	}
	app.logger.Infow("cache hit", "key", "user", "id", userID)
	//check cache response
	user, err := app.cacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.logger.Infow("fetching from DB", "id", userID)
		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if err := app.cacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (app *application) rateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.rateLimiter.Enabled {
			if allow, retryAfter := app.rateLimiter.Allow(r.RemoteAddr); !allow {
				app.rateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
