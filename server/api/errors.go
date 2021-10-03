package api

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// ApiError holds possible http api error fields
type ApiError struct {
	Err error `json:"-"`

	StatusCode int    `json:"-"`
	StatusText string `json:"status"`

	Location  string      `json:"location,omitempty"`
	AppCode   int64       `json:"code,omitempty"`
	ErrorText string      `json:"error,omitempty"`
	Cause     string      `json:"cause,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

// Error return an error text
func (e *ApiError) Error() string {
	return e.ErrorText
}

// Render sends error message to the client
func (e *ApiError) Render(w http.ResponseWriter, r *http.Request) error {
	pc := make([]uintptr, 5) // maximum 5 levels to go
	runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc)
	next := false
	for {
		frame, more := frames.Next()
		if next {
			e.Location = fmt.Sprintf("%s:%d", frame.File, frame.Line)
		}
		if strings.Contains(frame.File, "api/renderer.go") {
			next = true
		}
		if !more {
			break
		}
	}
	w.WriteHeader(e.StatusCode)
	return nil
}

// ErrUnauthorized is error message for Unauthorized
func ErrUnauthorized(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusUnauthorized,
		StatusText: "Unauthorized",
		ErrorText:  err.Error(),
	}
}

// ErrInvalidRequest is error message for Unauthorized
func ErrInvalidRequest(err error, data ...interface{}) *ApiError {
	v := &ApiError{
		Err:        err,
		StatusCode: http.StatusBadRequest,
		StatusText: "Invalid request.",
		ErrorText:  err.Error(),
	}
	if len(data) > 0 {
		if errText, ok := data[0].(string); ok {
			v.ErrorText = fmt.Sprintf("%s %s", v.ErrorText, errText)
		}
	}
	return v
}

// ErrServiceUnavailable is error message for Service Unavailable
func ErrServiceUnavailable(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusServiceUnavailable,
		StatusText: "Service Unavailable.",
		ErrorText:  err.Error(),
	}
}

// ErrInternalServerError is error message for Internal Server.
func ErrInternalServerError(err error) *ApiError {
	return &ApiError{
		Err:        err,
		StatusCode: http.StatusInternalServerError,
		StatusText: "Internal Server.",
		ErrorText:  err.Error(),
	}
}
