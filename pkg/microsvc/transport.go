package microsvc

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"

	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

type contextKey string

const (
	// ContextKeyUser - context key for user
	ContextKeyUser = contextKey("user")
	// ContextKeyEmail - context key for email
	ContextKeyEmail = contextKey("email")
	// ContextKeyToken - context key for token
	ContextKeyToken = contextKey("token")
)

var (
	// ErrBadRequest - invalid or malformed request
	ErrBadRequest = errors.New("bad request")
	// ErrBadRouting - bad routing error
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
	// ErrInvalidArgument - invalid argument error
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrForbidden - invalid argument error
	ErrForbidden = errors.New("you're forbidden to do this")
	// ErrInvalidToken - invalid argument error
	ErrInvalidToken = errors.New("token is invalid")
	// ErrUnauthorized - unauthorized
	ErrUnauthorized = errors.New("unauthorized")
	// ErrServiceUnhealthy - the service has become unhealthy
	ErrServiceUnhealthy = errors.New("service unhealthy")
	// ErrInvalidRequestBody - the request body was invalid and could not be
	// decoded
	ErrInvalidRequestBody = errors.New("invalid request body, could not decode")
	// ErrUnableToParseClaims - the claims could not be parsed
	ErrUnableToParseClaims = errors.New("unable to parse claims")
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
	c *models.Config,
) *mux.Router {

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerBefore(
			[]kithttp.RequestFunc{
				jwt.HTTPToContext(),
				xRequestIDToContext,
			}...,
		),
	}
	// create the router
	router := mux.NewRouter().StrictSlash(false)
	// / and /version handler for version

	router.Methods(http.MethodGet).Path("/health").Handler(
		kithttp.NewServer(
			MakeHealthEndpoint(s),
			decodeHealthRequest,
			encodeResponse,
			options...,
		),
	)

	api := router.PathPrefix("/api").Subrouter()

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

// xRequestIDToConext populates the context with a request id and sets the
// X-Request-ID http.Header
func xRequestIDToContext(ctx context.Context, r *http.Request) context.Context {
	id := uuid.NewV4().String()
	h := r.Header.Get("X-Request-Id")
	if len(h) > 0 {
		id = h
	} else {
		if v, ok := ctx.Value(kithttp.ContextKeyRequestXRequestID).(string); ok {
			id = v
		}
	}
	return context.WithValue(
		ctx,
		kithttp.ContextKeyRequestXRequestID,
		id,
	)
}

// version - returns the service version information
func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"arch":        Arch,
		"api_version": APIVersion,
		"build_time":  BuildTime,
		"git_commit":  GitCommit,
		"name":        Name,
		"version":     Ver,
	}); err != nil {
		return
	}
}

func decodeHealthRequest(
	_ context.Context,
	_ *http.Request,
) (interface{}, error) {
	return models.HealthRequest{}, nil
}

// decodeRequest.txt
