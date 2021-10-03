package post

import "github.com/go-chi/chi"

// Routes sets up the routes for the post service
func Routes() chi.Router {
	r := chi.NewRouter()

	r.Use(PostCtx)
	r.Get("/", GetPosts)

	return r
}
