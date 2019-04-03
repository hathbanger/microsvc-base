package microsvc

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/hathbanger/microsvc-base/pkg/microsvc/models"
)

type authorizor struct {
	groups     []string
	profileURL string
	signingKey *rsa.PublicKey
	logger     log.Logger
	client     *http.Client
}

type profile struct {
	Name     string   `json:"name"`
	MemberOf []string `json:"auth_groups"`
}

type claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	stdjwt.StandardClaims
}

// ValidateAuth - validates the token
func (a *authorizor) ValidateAuth() endpoint.Middleware {
	logger := log.With(a.logger, "event", "validate", "method", "ValidateAuth")
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			token, ok := ctx.Value(jwt.JWTTokenContextKey).(string)
			if !ok {
				err := errors.New("Token Not Found")
				logger.Log("err", err)
				return nil, errors.New("Token not found")
			}

			claims, err := a.parseToken(token)
			if err != nil {
				logger.Log("err", err)
				return nil, errors.New("Invalid token")
			}

			logger.Log("UserLogin", claims.Username)

			if claims.ExpiresAt == 0 {
				logger.Log("err", "Token has expired")
				return nil, errors.New("Token has expired")
			}

			if claims.Username == "" {
				logger.Log("err", "Username is invalid")
				return nil, errors.New("Username invalid")
			}

			if len(a.groups) > 0 {
				// get user profile info to verify groups
				logger.Log("message", "checking group membership")
				req, _ := http.NewRequest(http.MethodGet, a.profileURL, nil)
				req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))

				u := url.Values{}
				u.Set("group", strings.Join(a.groups, ","))
				req.URL.RawQuery = u.Encode()

				resp, err := a.client.Do(req)
				if err != nil {
					logger.Log("err", "unable to retrieve groups authorization")
					return nil, errors.New("Unable to retrieve groups authorization")
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					return nil, errors.New("user is forbidden")
				}
			}
			type contextKey string
			var (
				contextKeyUser = contextKey("user")
			)
			ctx = context.WithValue(
				context.Background(),
				contextKeyUser,
				claims.Username,
			)
			return next(ctx, request)
		}
	}
}

func (a *authorizor) parseToken(tokenString string) (c *claims, err error) {
	token, err := stdjwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *stdjwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Token is invalid")
	}

	if c, ok := token.Claims.(*claims); ok {
		return c, err
	}
	return nil, errors.New("Unable to parse claims")
}

// NewAuth - returns a new authorizor function
func NewAuth(logger log.Logger, c *models.Config) *authorizor {
	return &authorizor{
		groups:     c.Auth.Groups,
		profileURL: c.Auth.ProfileURL,
		logger:     logger,
		signingKey: c.PublicKey,
		client:     &http.Client{},
	}
}
