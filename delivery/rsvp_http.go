package delivery

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	rsvp "github.com/faris-arifiansyah/fws-rsvp"
	"github.com/faris-arifiansyah/fws-rsvp/constants"
	"github.com/faris-arifiansyah/fws-rsvp/handler"
	"github.com/faris-arifiansyah/fws-rsvp/middleware"
	"github.com/faris-arifiansyah/fws-rsvp/request"
	"github.com/faris-arifiansyah/fws-rsvp/request/validator"
	"github.com/faris-arifiansyah/fws-rsvp/response"
	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

// RsvpHandler struct
type RsvpHandler struct {
	uc  rsvp.Usecase
	rds *redis.Client
}

func NewRsvpHandler(uc rsvp.Usecase, rds *redis.Client) RsvpHandler {
	return RsvpHandler{
		uc:  uc,
		rds: rds,
	}
}

func (h *RsvpHandler) Register(router *httprouter.Router, ds []middleware.Decorator) error {
	if router == nil {
		return fmt.Errorf("router cannot be empty")
	}

	router.POST("/rsvps", handler.Decorate(handler.WithAuth(h.CreateRsvp, handler.Anonymous), ds...))
	router.GET("/rsvps", handler.Decorate(handler.WithAuth(h.RetrieveAllRsvp, handler.Admin), ds...))
	router.GET("/files/rsvps", handler.Decorate(handler.WithAuth(h.DownloadRsvpCsv, handler.Admin), ds...))

	return nil
}

func (h *RsvpHandler) CreateRsvp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var ctx = r.Context()
	var rsvpRequest rsvp.Rsvp
	var remoteAddr = r.RemoteAddr
	var err error

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&rsvpRequest); err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return err
	}
	defer r.Body.Close()

	errs := validator.Validate(rsvpRequest)
	if len(errs) > 0 {
		errBody := response.BuildErrors(errs)
		response.Write(w, errBody, http.StatusBadRequest)
		return err
	}

	//Check Rate Limit
	host, _, _ := net.SplitHostPort(remoteAddr)
	count, err := h.rds.Get(constants.RedisPrefix + host).Int()
	if count+1 > constants.RateLimit { //Rate Limit Exceeded
		err = response.RateLimitExceededError
	} else if count == 0 { //Set in Redis with Expire
		err = h.rds.Set(constants.RedisPrefix+host, 1, time.Duration(constants.RateLimitExp*time.Second)).Err()
	} else { //Increment Number of Requests
		err = h.rds.Incr(constants.RedisPrefix + host).Err()
	}

	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return err
	}

	//Create RSVP
	createdRsvp, err := h.uc.CreateRsvp(ctx, rsvpRequest)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return err
	}

	m := response.MetaInfo{HTTPStatus: http.StatusCreated}
	response.Write(w, response.BuildSuccess(createdRsvp, m), http.StatusCreated)
	return nil
}

func (h *RsvpHandler) RetrieveAllRsvp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
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
		return err
	}

	m := response.MetaInfo{
		HTTPStatus: http.StatusOK,
		Limit:      p.Limit,
		Offset:     p.Offset,
		Total:      rsvpResult.Total,
		Sort:       p.Sort,
	}

	response.Write(w, response.BuildSuccess(rsvpResult.Data, m), http.StatusOK)
	return nil
}

func (h *RsvpHandler) DownloadRsvpCsv(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	ctx := r.Context()

	qh := request.NewQueryHelper(r)
	p := rsvp.Parameter{
		Sort:   qh.GetString("sort", ""),
		Limit:  qh.GetInt("limit", 10),
		Offset: qh.GetInt("offset", 0),
	}

	file, err := h.uc.WriteRsvpsCsv(ctx, &p)
	if err != nil {
		errBody, httpStatus := response.BuildErrorAndStatus(err, "")
		response.Write(w, errBody, httpStatus)
		return err
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.Name))
	w.Header().Set("Content-Type", "text/csv")
	w.Write(file.Content)

	return nil
}
