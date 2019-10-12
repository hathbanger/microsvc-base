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
)

// Claims - struct for jwt claims
type Claims struct {
	stdjwt.StandardClaims
	LastName string `json:"LNAME"`
	UserName string `json:"USERNAME"`
	Email    string `json:"EMAIL"`
}

// AuthMiddleware - returns endpoint for auth validation
func AuthMiddleware(
	groups []string,
	profileURL string,
	signingKey *rsa.PublicKey,
) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			token, ok := ctx.Value(jwt.JWTTokenContextKey).(string)
			if !ok {
				return nil, ErrInvalidToken
			}
			claims, err := parseToken(token, signingKey)
			if err != nil {
				return nil, err
			}
			if claims.ExpiresAt == 0 || len(claims.UserName) <= 0 {
				return nil, ErrInvalidToken
			}
			ctx = context.WithValue(context.Background(), ContextKeyToken, token)
			if len(groups) > 0 {
				req, err := http.NewRequest(http.MethodGet, profileURL, nil)
				if err != nil {
					return nil, err
				}
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
				u := url.Values{}
				u.Set("group", strings.Join(groups, ","))
				req.URL.RawQuery = u.Encode()
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					return nil, errors.New("forbidden")
				}
			}
			return next(ctx, request)
		}
	}
}
func parseToken(tokenString string, key *rsa.PublicKey) (c *Claims, err error) {
	//token, err := stdjwt.ParseWithClaims(
	//	tokenString,
	//	&Claims{},
	//	func(token *stdjwt.Token) (interface{}, error) {
	//		return key, nil
	//	},
	//)
	//if err != nil {
	//	return nil, err
	//}
	// Check validity of token
	//if !token.Valid {
	//	return nil, ErrInvalidToken
	//}
	// Check claims have been parsed
	//if c, ok := token.Claims.(*Claims); ok {
	return c, err
	//}
	//return nil, ErrUnableToParseClaims
}
