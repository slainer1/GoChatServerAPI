package main

import (
	"main/internal/store"
	"net/http"
)

// GetFeed godoc
//
//	@Summary		Fetches a user's feed
//	@Description	Fetches a user's feed
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit results"		default(20)
//	@Param			offset	query		int		false	"Offset results"	default(0)
//	@Param			sort	query		string	false	"Sort direction"	Enums(asc, desc)	default(desc)
//	@Success		200		{array}		store.Post
//	@Failure		400		{object}	error
//	@Failure		401		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	// Pagination, filters, sorting endpoint
	fq := store.PaginationFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}
	var err error
	fq, err = fq.Parse(r)
	if err != nil {
		app.badRequestResponsee(w, r, err)
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponsee(w, r, err)
	}
	// Get the User's feed based on the created at.
	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(215), fq)
	if err != nil {
		app.statusInternalServerErrorHandler(w, r, err)
	}
	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.statusInternalServerErrorHandler(w, r, err)
	}
}
