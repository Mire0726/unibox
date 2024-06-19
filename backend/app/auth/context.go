package auth

import (
	"context"

	"firebase.google.com/go/auth"
	"github.com/Mire0726/unibox/backend/internal/cerror"
)

type contextKey string

var userTokenKey = contextKey("uerToken")

func SetUserToken(ctx context.Context, userToken string) context.Context {
	return context.WithValue(ctx, userTokenKey, userToken)
}

func GetUserTokenFromContext(ctx context.Context) string {
	var userToken string
	if ctx.Value(userTokenKey) != nil {
		userToken = ctx.Value(userTokenKey).(string)
	}
	return userToken
}

type uIDKey struct{}

type orderIDKey struct{}

type storeIDKey struct{}

func setContext(ctx context.Context, authToken *auth.Token, client *auth.Client) (context.Context, error) {
	user, err := client.GetUser(ctx, authToken.UID)
	if err != nil {
		return nil, cerror.Wrap(err, "auth",
			cerror.WithUnauthorizedCode(),
			cerror.WithReasonCode(cerror.RC20005),
		)
	}

	claims := user.CustomClaims
	if orderID, ok := claims["orderID"]; ok {
		ctx = setOrderID(ctx, orderID)
	} else {
		return nil, cerror.New("OrderID is not found in claims",
			cerror.WithUnauthorizedCode(),
			cerror.WithReasonCode(cerror.RC20006),
		)
	}

	if storeID, ok := claims["storeID"]; ok {
		ctx = setStoreID(ctx, storeID)
	} else {
		return nil, cerror.New("storeID is not found in claims",
			cerror.WithUnauthorizedCode(),
			cerror.WithReasonCode(cerror.RC20006),
		)
	}

	ctx = setUID(ctx, user.UID)

	return ctx, nil
}

func setUID(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, uIDKey{}, value)
}

func GetUID(ctx context.Context) (string, error) {
	uid, ok := ctx.Value(uIDKey{}).(string)
	if !ok {
		return "", cerror.New("UID is not found in context", cerror.WithNotFoundCode())
	}

	return uid, nil
}

func setOrderID(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, orderIDKey{}, value)
}

func GetOrderID(ctx context.Context) (string, error) {
	orderID, ok := ctx.Value(orderIDKey{}).(string)
	if !ok {
		return "", cerror.New("OrderID is not found in context", cerror.WithNotFoundCode())
	}

	return orderID, nil
}

func setStoreID(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, storeIDKey{}, value)
}

func GetStoreID(ctx context.Context) (string, error) {
	storeID, ok := ctx.Value(storeIDKey{}).(string)
	if !ok {
		return "", cerror.New("StoreID is not found in context", cerror.WithNotFoundCode())
	}

	return storeID, nil
}
