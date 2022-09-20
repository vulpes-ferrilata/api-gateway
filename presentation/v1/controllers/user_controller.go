package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/mappers"
	"github.com/vulpes-ferrilata/shared/proto/v1/user"
)

func NewUserController(userClient user.UserClient) *UserController {
	return &UserController{
		userClient: userClient,
	}
}

type UserController struct {
	userClient user.UserClient
}

func (u UserController) GetBy(ctx iris.Context, id string) (mvc.Result, error) {
	getUserByIDGrpcRequest := &user.GetUserByIDRequest{
		ID: id,
	}

	userGrpcResponse, err := u.userClient.GetUserByID(ctx, getUserByIDGrpcRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResponse := mappers.ToUserHttpResponse(userGrpcResponse)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: userResponse,
	}, nil
}
