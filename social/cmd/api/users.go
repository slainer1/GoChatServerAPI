package main

import (
	_ "errors"
	"main/internal/store"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}	[get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponsee(w, r, err)
		return
	}
	ctx := r.Context()

	user, err := app.getUser(ctx, userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.statusInternalServerErrorHandler(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.statusInternalServerErrorHandler(w, r, err)
	}
}

// FollowUser godoc
//
//	@Summary		Follows a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Target ID to follow"
//	@Success		200	{object}	store.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow	[put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromContext(r)
	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponsee(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.statusInternalServerErrorHandler(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Follow(ctx, followerUser.ID, followedID); err != nil {
		switch err {
		case store.ErrorConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.statusInternalServerErrorHandler(w, r, err)
		}

	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollows a user profile
//	@Description	Unfollows a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int		true	"Target User ID to unfollow"
//	@Success		204		{string}	string	"No Content"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromContext(r)
	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponsee(w, r, err)
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, unfollowedID); err != nil {
		switch err {
		case store.ErrorConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.statusInternalServerErrorHandler(w, r, err)
		}
	}
}

// ActivateUser godoc
//
//	@Summary		Activates/registers a user
//	@Description	Activates/registers a user by invitation token
//	@Tags			users
//	@Produce		json
//	@Param			token	path		string	true	"Invitation Token"
//	@Success		204		{string}	string	"No Content"
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	err := app.store.Users.Activate(r.Context(), token)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.statusInternalServerErrorHandler(w, r, err)
		}
		return
	}
	if err := app.jsonResponse(w, http.StatusNoContent, ""); err != nil {
		app.statusInternalServerErrorHandler(w, r, err)
	}
}

func (app *application) getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
