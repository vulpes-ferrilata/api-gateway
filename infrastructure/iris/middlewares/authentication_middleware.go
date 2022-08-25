package middlewares

import (
	"fmt"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/shared/proto/v1/authentication"
)

type TokenExtractor func(iris.Context) (string, error)

func FromAuthHeader(ctx iris.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
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

func NewAuthenticationMiddleware(authenticationClient authentication.AuthenticationClient,
	errorHandlerMiddleware *ErrorHandlerMiddleware) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authenticationClient:   authenticationClient,
		errorHandlerMiddleware: errorHandlerMiddleware,
	}
}

type AuthenticationMiddleware struct {
	authenticationClient   authentication.AuthenticationClient
	errorHandlerMiddleware *ErrorHandlerMiddleware
}

func (a AuthenticationMiddleware) Handle(ctx iris.Context) {
	accessToken, err := FromFirst(
		FromAuthHeader,
		FromParameter("token"),
	)(ctx)
	if err != nil {
		a.errorHandlerMiddleware.Handle(ctx, err)
		return
	}

	getClaimByAccessTokenGrpcRequest := &authentication.GetClaimByAccessTokenRequest{
		AccessToken: accessToken,
	}

	claimGrpcResponse, err := a.authenticationClient.GetClaimByAccessToken(ctx.Request().Context(), getClaimByAccessTokenGrpcRequest)
	if err != nil {
		a.errorHandlerMiddleware.Handle(ctx, err)
		return
	}

	request := ctx.Request()
	requestCtx := request.Context()
	requestCtx = context_values.WithUserID(requestCtx, claimGrpcResponse.UserID)
	request = request.WithContext(requestCtx)
	ctx.ResetRequest(request)

	ctx.Next()
}
