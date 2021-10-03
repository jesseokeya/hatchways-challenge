package post

import (
	"context"
	"errors"
	"net/http"
	"posts/v1/data"
	"posts/v1/lib/hatchways"
	"posts/v1/server/api"
)

var (
	ErrInvalidSortBy = errors.New("sortBy parameter is invalid")
)

type PostResponse struct {
	*hatchways.PostResponse
}

func (p *PostResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// PostCtx is a context for the post handler. Fetches the posts puts it in context and passes to the next handler.
func PostCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tags, sortBy, direction := r.FormValue("tags"), r.FormValue("sortBy"), r.FormValue("direction")
		q := &hatchways.FilterQueryParams{Tags: tags, SortBy: sortBy, Direction: direction}

		validSortBy := map[string]bool{
			"id":         true,
			"reads":      true,
			"likes":      true,
			"popularity": true,
		}

		_, ok := validSortBy[sortBy]
		if sortBy != "" && !ok {
			api.Render(w, r, api.ErrInvalidRequest(ErrInvalidSortBy))
			return
		}

		posts, err := data.GetPosts(ctx, q)
		if err != nil {
			api.Render(w, r, api.ErrInvalidRequest(err))
		}

		ctx = context.WithValue(ctx, api.PostContextKey, &PostResponse{posts})
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// Post returns the posts from the context.
func GetPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	posts := ctx.Value(api.PostContextKey).(*PostResponse)
	api.Render(w, r, posts)
}
