package main

import (
	"net/http"
	"social/internal/store"
)

// get user feed godoc
//
//	@Summary		Get User Feed
//	@Description	Retrieves the feed for a specific user with pagination and sorting options
//	@Tags			feed
//	@Produce		json
//	@Param			limit	query		int		false	"Number of posts to return"	default(20)
//	@Param			offset	query		int		false	"Number of posts to skip"	default(0)
//	@Param			sort	query		string	false	"Sort order: asc or desc"	default(desc)	Enum(asc, desc)
//	@Success		200		{object}	[]store.PostWithMetaData
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/feed [get]
func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	// pagination, filtering and sorting

	fq := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return

	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeed(ctx, int64(5), fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)

	}

}
