package authentication

import (
	"context"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/saga"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication/requests"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication/responses"
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

// @Summary Register user
// @Description Create new user
// @Accept  json
// @Produce  json
// @Param	displayName    body    requests.Register	true	"Display Name"
// @Param	email	   body    requests.Register	true	"Email"
// @Param	password	   body    requests.Register	true	"Password"
// @Success 201 {object} responses.User	"ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 422 {object} iris.Problem "email is already exists"
// @Failure 422 {object} iris.Problem "unable to encrypt password"
// @Router /auth/register [post]
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

	userResponse := &responses.User{
		ID: userID,
	}

	return &mvc.Response{
		Code:   iris.StatusCreated,
		Object: userResponse,
	}, nil
}

// @Summary Login user
// @Description Get tokens
// @Accept  json
// @Produce  json
// @Param	email	   body    requests.Login	true	"Email"
// @Param	password	   body    requests.Login	true	"Password"
// @Success 200 {object} responses.Token	"ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "user credential not found"
// @Failure 422 {object} iris.Problem "password is invalid"
// @Router /auth/login [post]
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

// @Summary Refresh token
// @Description Provide new access token
// @Accept  json
// @Produce  json
// @Param	refreshToken	   body    requests.Refresh	true	"Refresh Token"
// @Success 200 {object} responses.Token	"ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "claim not found"
// @Failure 422 {object} iris.Problem "token has been expired"
// @Failure 422 {object} iris.Problem "token has been revoked"
// @Router /auth/refresh [post]
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

// @Summary Revoke token
// @Description Invalidate tokens
// @Accept  json
// @Produce  json
// @Param	refreshToken	   body    requests.Revoke	true	"Refresh Token"
// @Success 200 {object} responses.Token	"ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "claim not found"
// @Failure 422 {object} iris.Problem "token has been expired"
// @Failure 422 {object} iris.Problem "token has been revoked"
// @Router /auth/revoke [post]
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
