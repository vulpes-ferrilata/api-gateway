package controllers

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/saga"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/requests"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/authentication"
	"github.com/vulpes-ferrilata/shared/proto/v1/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewAuthenticationController(userClient user.UserClient,
	authenticationClient authentication.AuthenticationClient) *AuthenticationController {
	return &AuthenticationController{
		userClient:           userClient,
		authenticationClient: authenticationClient,
	}
}

type AuthenticationController struct {
	userClient           user.UserClient
	authenticationClient authentication.AuthenticationClient
}

func (a AuthenticationController) PostRegister(ctx iris.Context) (mvc.Result, error) {
	registerRequest := &requests.Register{}

	if err := ctx.ReadJSON(registerRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userCredentialID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	coordinator := saga.NewCoordinator()

	if err := coordinator.Execute(ctx.Request().Context(),
		&saga.Step{
			Handle: func(ctx context.Context) error {
				createUserCredentialGrpcRequest := &authentication.CreateUserCredentialRequest{
					ID:       userCredentialID,
					UserID:   userID,
					Email:    registerRequest.Email,
					Password: registerRequest.Password,
				}

				if _, err := a.authenticationClient.CreateUserCredential(ctx, createUserCredentialGrpcRequest); err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
			Compensate: func(ctx context.Context) error {
				deleteUserCredentialGrpcRequest := &authentication.DeleteUserCredentialRequest{
					ID: userCredentialID,
				}

				if _, err := a.authenticationClient.DeleteUserCredential(ctx, deleteUserCredentialGrpcRequest); err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
		},
		&saga.Step{
			Handle: func(ctx context.Context) error {
				createUserGrpcRequest := &user.CreateUserRequest{
					ID:          userID,
					DisplayName: registerRequest.DisplayName,
				}

				if _, err := a.userClient.CreateUser(ctx, createUserGrpcRequest); err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
		}); err != nil {
		return nil, errors.WithStack(err)
	}

	userCreatedResponse := responses.UserCreated{
		ID: userID,
	}

	return &mvc.Response{
		Code:   iris.StatusCreated,
		Object: userCreatedResponse,
	}, nil
}

func (a AuthenticationController) PostLogin(ctx iris.Context) (mvc.Result, error) {
	loginRequest := &requests.Login{}

	if err := ctx.ReadJSON(loginRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	claimID := primitive.NewObjectID().Hex()

	loginGrpcRequest := &authentication.LoginRequest{
		ClaimID:  claimID,
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	if _, err := a.authenticationClient.Login(ctx.Request().Context(), loginGrpcRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	getTokenByClaimIDGrpcRequest := &authentication.GetTokenByClaimIDRequest{
		ClaimID: claimID,
	}

	tokenGrpcResponse, err := a.authenticationClient.GetTokenByClaimID(ctx.Request().Context(), getTokenByClaimIDGrpcRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenResponse := mappers.ToTokenHttpResponse(tokenGrpcResponse)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: tokenResponse,
	}, nil
}

func (a AuthenticationController) PostRefresh(ctx iris.Context) (mvc.Result, error) {
	refreshRequest := &requests.Refresh{}

	if err := ctx.ReadJSON(refreshRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	getTokenByRefreshTokenGrpcRequest := &authentication.GetTokenByRefreshTokenRequest{
		RefreshToken: refreshRequest.RefreshToken,
	}

	tokenGrpcResponse, err := a.authenticationClient.GetTokenByRefreshToken(ctx.Request().Context(), getTokenByRefreshTokenGrpcRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenResponse := mappers.ToTokenHttpResponse(tokenGrpcResponse)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: tokenResponse,
	}, nil
}

func (a AuthenticationController) PostRevoke(ctx iris.Context) (mvc.Result, error) {
	revokeRequest := &requests.Revoke{}

	if err := ctx.ReadJSON(revokeRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	revokeTokenGrpcRequest := &authentication.RevokeTokenRequest{
		RefreshToken: revokeRequest.RefreshToken,
	}

	if _, err := a.authenticationClient.RevokeToken(ctx.Request().Context(), revokeTokenGrpcRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}
