package middlewares

import (
	"fmt"
	"strings"

	"github.com/go-playground/pure"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/hero"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/authentication-service-proto/pb"
	"github.com/vulpes-ferrilata/authentication-service-proto/pb/requests"
)

type TokenExtractor func(iris.Context) (string, error)

func FromAuthHeader(ctx iris.Context) (string, error) {
	authHeader := ctx.GetHeader(pure.Authorization)
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}

func FromParameter(param string) TokenExtractor {
	return func(ctx iris.Context) (string, error) {
		return ctx.URLParam(param), nil
	}
}

func FromFirst(extractors ...TokenExtractor) TokenExtractor {
	return func(ctx iris.Context) (string, error) {
		for _, ex := range extractors {
			token, err := ex(ctx)
			if err != nil {
				return "", err
			}
			if token != "" {
				return token, nil
			}
		}
		return "", nil
	}
}

func NewAuthenticationMiddleware(authenticationClient pb.AuthenticationClient,
	errorHandlerMiddleware *ErrorHandlerMiddleware) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authenticationClient: authenticationClient,
		errorHandlerFunc:     errorHandlerMiddleware.Serve(),
	}
}

type AuthenticationMiddleware struct {
	authenticationClient pb.AuthenticationClient
	errorHandlerFunc     hero.ErrorHandlerFunc
}

func (a AuthenticationMiddleware) Serve() iris.Handler {
	return func(ctx iris.Context) {
		accessToken, err := FromFirst(
			FromAuthHeader,
			FromParameter("token"),
		)(ctx)
		if err != nil {
			a.errorHandlerFunc(ctx, err)
			return
		}

		getClaimByAccessTokenPbRequest := &requests.GetClaimByAccessToken{
			AccessToken: accessToken,
		}

		claimPbResponse, err := a.authenticationClient.GetClaimByAccessToken(ctx.Request().Context(), getClaimByAccessTokenPbRequest)
		if err != nil {
			a.errorHandlerFunc(ctx, err)
			return
		}

		request := ctx.Request()
		requestCtx := request.Context()
		requestCtx = context_values.WithUserID(requestCtx, claimPbResponse.UserID)
		request = request.WithContext(requestCtx)
		ctx.ResetRequest(request)

		ctx.Next()
	}
}
