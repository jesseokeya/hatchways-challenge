package api

// ContextKey is a type alias for string.
// This technique for defining context keys was copied from Go 1.7's new use of context in net/http.
type ContextKey struct {
	Name string
}

var (
	// PostContextKey is the context.Context key to store hatchway post request context.
	PostContextKey = &ContextKey{Name: "Posts"}
)
