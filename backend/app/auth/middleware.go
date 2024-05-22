package auth

import (
	"strings"

	"firebase.google.com/go/auth"
	"github.com/labstack/echo/v4"

	"github.com/Mire0726/unibox/backend/infrastructure/firebase"
	"github.com/Mire0726/unibox/backend/internal/cerror"
)

const (
	authHeader = "Authorization"
	authScheme = "Bearer "
	skipPath   = "/v1/anonymous-user"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (r *AuthMiddleware) AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			if req.RequestURI == skipPath {
				return next(c)
			}

			client, err := r.getFirebaseClient(c)
			if err != nil {
				return err
			}

			if err := r.authByAccessToken(c, client); err != nil {
				return err
			}

			return next(c)
		}
	}
}

func (r *AuthMiddleware) getFirebaseClient(c echo.Context) (*auth.Client, error) {
	ctx := c.Request().Context()
	client, err := firebase.NewClientWithoutLogger(ctx)
	if err != nil {
		return nil, cerror.Wrap(err, "auth", cerror.WithInternalCode(), cerror.WithReasonCode(cerror.RC20001))
	}

	return client, nil
}

func (r *AuthMiddleware) authByAccessToken(c echo.Context, client *auth.Client) error {
	req := c.Request()
	ctx := req.Context()

	token := req.Header.Get(authHeader)
	if len(token) == 0 {
		return cerror.New("Authorization header is not found",
			cerror.WithUnauthorizedCode(),
			cerror.WithReasonCode(cerror.RC20001),
		)
	}

	idToken := strings.TrimSpace(strings.TrimPrefix(token, authScheme))
	authToken, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return cerror.Wrap(err, "auth", cerror.WithUnauthorizedCode(), cerror.WithReasonCode(cerror.RC20004))
	}

	ctx, err = setContext(ctx, authToken, client)
	if err != nil {
		return err
	}

	c.SetRequest(req.WithContext(ctx))

	return nil
}
