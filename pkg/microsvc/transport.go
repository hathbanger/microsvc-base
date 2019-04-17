package microsvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

var (
	// ErrBadRequest - bad request
	ErrBadRequest = errors.New("bad request")
	// ErrBadRouting - bad routing
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
	// ErrInvalidArgument - invalid argument
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrUnauthorized - unauthorized
	ErrUnauthorized = errors.New("unauthorized")
	// ErrServiceUnhealthy - service unhealthy
	ErrServiceUnhealthy = errors.New("service unhealthy")
)

// Failure - should be implemented by response types so encoders can check that
// a sresponse can be asserted to Failure
type Failure interface {
	Failed() error
}

// MakeRoutes - makes all the routes
func MakeRoutes(
	s Service,
	logger log.Logger,
	mw []endpoint.Middleware,
) *mux.Router {

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(
			[]kithttp.RequestFunc{
				jwt.HTTPToContext(),
			}...,
		),
	}

	router := mux.NewRouter().StrictSlash(false)

	router.Methods(http.MethodGet).Path("/health").Handler(
		kithttp.NewServer(
			MakeHealthEndpoint(s),
			decodeHealthRequest,
			encodeResponse,
			options...,
		),
	)

	api := router.PathPrefix("/api").Subrouter()

	// routes - start
	foo := kithttp.NewServer(
		MakeFooEndpoint(s, logger, mw),
		decodeFooRequest,
		encodeResponse,
		options...,
	)
	api.Methods(http.MethodPost).Path("/v1/foo").Handler(foo)
	// routes - finish

	// transport.txt

	// plug in metrics here:
	// router.Handle("/metrics", promhttp.Handler())

	return router
}

func encodeResponse(
	ctx context.Context,
	w http.ResponseWriter,
	response interface{},
) error {
	if err, ok := response.(Failure); ok && err.Failed() != nil {
		encodeError(ctx, err.Failed(), w)
		return err.Failed()
	}

	if health, ok := response.(bool); ok {
		if !health {
			encodeError(ctx, ErrServiceUnhealthy, w)
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	var code int
	switch err {
	case ErrBadRouting:
		code = http.StatusBadRequest
	case ErrBadRequest:
		code = http.StatusBadRequest
	case ErrInvalidArgument:
		code = http.StatusBadRequest
	case ErrBadRouting:
		code = http.StatusBadRequest
	case ErrBadRouting:
		code = http.StatusBadRequest
	default:
		code = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err.Error(),
		"http_code":   code,
		"http_status": http.StatusText(code),
	})
}

func decodeHealthRequest(
	_ context.Context,
	_ *http.Request,
) (interface{}, error) {
	return models.HealthRequest{}, nil
}

func decodeFooRequest(
	_ context.Context,
	r *http.Request,
) (interface{}, error) {
	var fooRequest models.FooRequest
	if err := json.NewDecoder(r.Body).Decode(&fooRequest); err != nil {
		return nil, ErrBadRouting
	}
	return fooRequest, nil
}

// decodeRequest.txt
