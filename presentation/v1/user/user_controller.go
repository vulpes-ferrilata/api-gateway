package user

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/user/mappers"
	"github.com/vulpes-ferrilata/user-service-proto/pb"
	pb_requests "github.com/vulpes-ferrilata/user-service-proto/pb/requests"
)

func NewUserController(userClient pb.UserClient) *UserController {
	return &UserController{
		userClient: userClient,
	}
}

type UserController struct {
	userClient pb.UserClient
}

func (u UserController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodGet, "/me", "Me")
}

func (u UserController) GetBy(ctx iris.Context, id string) (mvc.Result, error) {
	getUserByIDPbRequest := &pb_requests.GetUserByID{
		UserID: id,
	}

	userPbResponse, err := u.userClient.GetUserByID(ctx, getUserByIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResponse := mappers.ToUserHttpResponse(userPbResponse)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: userResponse,
	}, nil
}

func (u UserController) Me(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	getUserByIDPbRequest := &pb_requests.GetUserByID{
		UserID: userID,
	}

	userPbResponse, err := u.userClient.GetUserByID(ctx, getUserByIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResponse := mappers.ToUserHttpResponse(userPbResponse)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: userResponse,
	}, nil
}
