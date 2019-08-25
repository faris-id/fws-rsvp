package middleware

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// HandleWithError is an httprouter.Handle that returns an error.
type HandleWithError func(http.ResponseWriter, *http.Request, httprouter.Params) error

// Decorator decorates HandleWithError.
type Decorator func(HandleWithError) HandleWithError

type ctxKey string

// HTTP runs HandleWithError and converts it to httprouter.Handle.
// The conversion is needed because httprouter.Router needs httprouter.Handle
// in its signature.
func HTTP(handle HandleWithError) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		handle(w, r, params)
	}
}

// WithLogging decorates Decorator with logging.
func WithLogging(logger *zap.Logger) Decorator {
	return func(handle HandleWithError) HandleWithError {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
			reqID := r.Header.Get("X-Request-ID")
			start := time.Now()

			err := handle(w, r, params)

			// elapsed time in milliseconds
			elapsed := time.Since(start).Seconds() * 1000
			elapsedStr := strconv.FormatFloat(elapsed, 'f', -1, 64)

			if err != nil {
				logger.Error(err.Error(),
					zap.String("request_id", reqID),
					zap.String("duration", elapsedStr),
					zap.Strings("tags", []string{r.URL.Path, r.Method}),
				)
			} else {
				logger.Info("everything is fine",
					zap.String("request_id", reqID),
					zap.String("duration", elapsedStr),
					zap.Strings("tags", []string{r.URL.Path, r.Method}),
				)
			}

			return err
		}
	}
}

// WithStandardContext decorates Decorator with standard context.
func WithStandardContext() Decorator {
	return func(handle HandleWithError) HandleWithError {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
			ctx := r.Context()

			if r.Header.Get("X-Request-ID") == "" {
				reqID, err := createRequestID(r)
				if err != nil {
					return err
				}
				r.Header.Set("X-Request-ID", reqID)
			}

			ctx = context.WithValue(ctx, ctxKey("X-Request-ID"), r.Header.Get("X-Request-ID"))
			ctx = context.WithValue(ctx, ctxKey("Authorization"), r.Header.Get("Authorization"))
			ctx = context.WithValue(ctx, ctxKey("Retry"), r.Header.Get("Retry"))

			return handle(w, r.WithContext(ctx), params)
		}
	}
}

// StandardDecorators returns standard decorators.
//
// WithLogging(),
// WithStandardContext()
func StandardDecorators() []Decorator {
	l, _ := zap.NewProduction()

	ds := []Decorator{
		WithLogging(l),
		WithStandardContext(),
	}
	return ds
}

// ApplyDecorators returns decorated HandleWithError.
func ApplyDecorators(handle HandleWithError, ds ...Decorator) HandleWithError {
	for _, d := range ds {
		handle = d(handle)
	}
	return handle
}

func createRequestID(r *http.Request) (string, error) {
	reqID := r.Header.Get("X-Request-ID")
	if reqID == "" {
		temp, err := uuid.NewRandom()
		if err != nil {
			return "", err
		}
		return temp.String(), nil
	}
	return reqID, nil
}
