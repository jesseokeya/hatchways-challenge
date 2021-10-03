package api

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// Status customizes http response status
type Status struct {
	Code int    `json:"code"`
	Text string `json:"text,omitempty"`
}

// renderResponder is a custom render.Responder which returns a custom response
func renderResponder(w http.ResponseWriter, r *http.Request, v interface{}) {
	if v, ok := v.(Status); ok {
		w.WriteHeader(v.Code)
	}
	render.DefaultResponder(w, r, v)
}

// Render customizes http response status
func Render(w http.ResponseWriter, r *http.Request, v render.Renderer) {
	if err := render.Render(w, r, v); err != nil {
		log.Error().Timestamp().
			Str("error", err.Error()).
			Send()
	}
}

// init registers custom render.Responder
func init() {
	// inject and override defaults in the render package
	render.Respond = renderResponder
}
