package handler

import (
	"fmt"
	"net/http"

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
	Register(r *httprouter.Router) error
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

	// start route
	for _, reg := range registrations {
		reg.Register(router)
	}

	router.NotFound = http.HandlerFunc(NotFound)

	return router, nil
}

func WithAuth(h func(http.ResponseWriter, *http.Request, httprouter.Params), authType AuthType) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		h(w, r, params)
	}
}
