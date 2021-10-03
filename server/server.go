package server

import (
	"net/http"
	"os"
	"posts/v1/server/api"
	post "posts/v1/server/posts"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

// Server holds refernces to other interfaces to be used
type Server struct {
	lg   *zerolog.Logger
	opts ServerOptions
}

// ServerOptions holds the options for the server
type ServerOptions struct {
	debug bool
}

// ServerOption is a function that configures the server
type ServerOption func(*ServerOptions) error

// Debug sets the debug mode
func Debug(b bool) ServerOption {
	return func(o *ServerOptions) error {
		o.debug = b
		return nil
	}
}

// New initializes the server
func New(options ...ServerOption) (*Server, error) {
	logger := zerolog.New(os.Stdout)

	s := &Server{}
	for _, opt := range options {
		if err := opt(&s.opts); err != nil {
			return nil, err
		}
	}
	if s.opts.debug {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	s.lg = &logger
	return s, nil
}

// HealthCheck is a health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`👾`))
}

// Ping is a ping endpoint
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{ "success": true }`))
}

// NewStructuredLogger returns a new structured logger
func (s *Server) NewStructuredLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&api.StructuredLogger{Logger: s.lg})
}

// Routes initializes the middlewares and routes to be used in the application
func (s *Server) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.NoCache)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(s.NewStructuredLogger())

	r.Get("/", HealthCheck)

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", Ping)
		r.Mount("/posts", post.Routes())
	})

	return r
}
