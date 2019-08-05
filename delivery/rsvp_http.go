package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	rsvp "github.com/faris-arifiansyah/fws-rsvp"
	"github.com/faris-arifiansyah/fws-rsvp/handler"
	"github.com/faris-arifiansyah/fws-rsvp/request"
	"github.com/faris-arifiansyah/fws-rsvp/response"
	"github.com/julienschmidt/httprouter"
)

// RsvpHandler struct
type RsvpHandler struct {
	uc rsvp.Usecase
}

func NewRsvpHandler(uc rsvp.Usecase) RsvpHandler {
	return RsvpHandler{
		uc: uc,
	}
}

func (h *RsvpHandler) Register(router *httprouter.Router) error {
	if router == nil {
		return fmt.Errorf("router cannot be empty")
	}

	router.POST("/rsvps", handler.WithAuth(h.CreateRsvp, handler.Anonymous))
	router.GET("/rsvps", handler.WithAuth(h.RetrieveAllRsvp, handler.Admin))

	return nil
}

func (h *RsvpHandler) CreateRsvp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var rsvpRequest rsvp.Rsvp

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rsvpRequest); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()
	createdRsvp, err := h.uc.CreateRsvp(ctx, rsvpRequest)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	m := response.MetaInfo{HTTPStatus: http.StatusCreated}
	response.Write(w, response.BuildSuccess(createdRsvp, m), http.StatusCreated)
}

func (h *RsvpHandler) RetrieveAllRsvp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	qh := request.NewQueryHelper(r)
	p := rsvp.Parameter{
		Sort:   qh.GetString("sort", ""),
		Limit:  qh.GetInt("limit", 10),
		Offset: qh.GetInt("offset", 0),
	}

	rsvpResult, err := h.uc.GetRsvps(ctx, &p)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return
	}

	m := response.MetaInfo{
		HTTPStatus: http.StatusOK,
		Limit:      p.Limit,
		Offset:     p.Offset,
		Total:      rsvpResult.Total,
		Sort:       p.Sort,
	}
	response.Write(w, response.BuildSuccess(rsvpResult.Data, m), http.StatusOK)
}
