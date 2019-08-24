package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/faris-arifiansyah/fws-rsvp/middleware"
	"github.com/faris-arifiansyah/fws-rsvp/response"
	"github.com/julienschmidt/httprouter"
)

// AuthType is custom type for decorator key type
type AuthType string

const (
	Anonymous AuthType = "anonymous"
	Admin     AuthType = "admin"
)

type Registration interface {
	Register(r *httprouter.Router, ds []middleware.Decorator) error
}

func Healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "ok")
}

func NotFound(w http.ResponseWriter, _ *http.Request) {
	meta := response.MetaInfo{HTTPStatus: 404}
	res := response.BuildSuccess("path not found", meta)
	response.Write(w, res, meta.HTTPStatus)
}

func NewHandler(registrations ...Registration) (http.Handler, error) {
	router := httprouter.New()
	router.HandleMethodNotAllowed = false

	router.HandlerFunc("GET", "/healthz", Healthz)

	// decorator for delivery
	sd := middleware.StandardDecorators()

	// start route
	for _, reg := range registrations {
		reg.Register(router, sd)
	}

	router.NotFound = http.HandlerFunc(NotFound)

	return router, nil
}

func WithAuth(h func(http.ResponseWriter, *http.Request, httprouter.Params) error, authType AuthType) middleware.HandleWithError {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
		if authType == Admin {
			headerUsername, headerPass, ok := r.BasicAuth()

			if !ok {
				response.Write(w, response.BuildError([]error{response.UserUnauthorizedError}), response.UserUnauthorizedError.HTTPCode)
				return response.UserUnauthorizedError
			}

			username := os.Getenv("FWS_RSVP_USERNAME")
			pass := os.Getenv("FWS_RSVP_PASSWORD")

			if username != headerUsername || pass != headerPass {
				response.Write(w, response.BuildError([]error{response.UserUnauthorizedError}), response.UserUnauthorizedError.HTTPCode)
				return response.UserUnauthorizedError
			}
		}

		return h(w, r, params)
	}
}

// Decorate util to simplify combining middleware
func Decorate(handle middleware.HandleWithError, ds ...middleware.Decorator) httprouter.Handle {
	return middleware.HTTP(middleware.ApplyDecorators(handle, ds...))
}
