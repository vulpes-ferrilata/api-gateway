package authentication

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/saga"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication/requests"
	authentication_pb "github.com/vulpes-ferrilata/authentication-service-proto/pb"
	authentication_pb_requests "github.com/vulpes-ferrilata/authentication-service-proto/pb/requests"
	user_pb "github.com/vulpes-ferrilata/user-service-proto/pb"
	user_pb_requests "github.com/vulpes-ferrilata/user-service-proto/pb/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewAuthenticationController(userClient user_pb.UserClient,
	authenticationClient authentication_pb.AuthenticationClient) *AuthenticationController {
	return &AuthenticationController{
		userClient:           userClient,
		authenticationClient: authenticationClient,
	}
}

type AuthenticationController struct {
	userClient           user_pb.UserClient
	authenticationClient authentication_pb.AuthenticationClient
}

func (a AuthenticationController) PostRegister(ctx iris.Context) (mvc.Result, error) {
	registerRequest := &requests.Register{}

	if err := ctx.ReadJSON(registerRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userCredentialID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()

	coordinator := saga.NewCoordinator()

	if err := coordinator.Execute(ctx,
		&saga.Step{
			Handle: func(ctx context.Context) error {
				createUserCredentialPbRequest := &authentication_pb_requests.CreateUserCredential{
					UserCredentialID: userCredentialID,
					UserID:           userID,
					Email:            registerRequest.Email,
					Password:         registerRequest.Password,
				}

				if _, err := a.authenticationClient.CreateUserCredential(ctx, createUserCredentialPbRequest); err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
			Compensate: func(ctx context.Context) error {
				deleteUserCredentialPbRequest := &authentication_pb_requests.DeleteUserCredential{
					UserCredentialID: userCredentialID,
				}

				if _, err := a.authenticationClient.DeleteUserCredential(ctx, deleteUserCredentialPbRequest); err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
		},
		&saga.Step{
			Handle: func(ctx context.Context) error {
				createUserPbRequest := &user_pb_requests.CreateUser{
					UserID:      userID,
					DisplayName: registerRequest.DisplayName,
				}

				if _, err := a.userClient.CreateUser(ctx, createUserPbRequest); err != nil {
					return errors.WithStack(err)
				}

				return nil
			},
		}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code: iris.StatusCreated,
		Object: &struct {
			ID string `json:"id"`
		}{
			ID: userID,
		},
	}, nil
}

func (a AuthenticationController) PostLogin(ctx iris.Context) (mvc.Result, error) {
	loginRequest := &requests.Login{}

	if err := ctx.ReadJSON(loginRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	claimID := primitive.NewObjectID().Hex()

	loginPbRequest := &authentication_pb_requests.Login{
		ClaimID:  claimID,
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	if _, err := a.authenticationClient.Login(ctx, loginPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	getTokenByClaimIDPbRequest := &authentication_pb_requests.GetTokenByClaimID{
		ClaimID: claimID,
	}

	tokenPbResponse, err := a.authenticationClient.GetTokenByClaimID(ctx, getTokenByClaimIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenResponse := mappers.ToTokenHttpResponse(tokenPbResponse)

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

	getTokenByRefreshTokenPbRequest := &authentication_pb_requests.GetTokenByRefreshToken{
		RefreshToken: refreshRequest.RefreshToken,
	}

	tokenPbResponse, err := a.authenticationClient.GetTokenByRefreshToken(ctx, getTokenByRefreshTokenPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	tokenResponse := mappers.ToTokenHttpResponse(tokenPbResponse)

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

	revokeTokenPbRequest := &authentication_pb_requests.RevokeToken{
		RefreshToken: revokeRequest.RefreshToken,
	}

	if _, err := a.authenticationClient.RevokeToken(ctx, revokeTokenPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}
