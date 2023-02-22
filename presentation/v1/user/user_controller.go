package user

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/user/mappers"
	"github.com/vulpes-ferrilata/user-service-proto/pb"
	pb_models "github.com/vulpes-ferrilata/user-service-proto/pb/models"
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

// @Summary Get user
// @Description Get user by id
// @Accept  json
// @Produce  json
// @Param   id     path    string     true        "User ID"
// @Success 200 {object} responses.User	"ok"
// @Failure 422 {object} iris.Problem "ID must be a valid ObjectID"
// @Failure 404 {object} iris.Problem "User not found"
// @Router /users/{id} [get]
func (u UserController) GetBy(ctx iris.Context, id string) (mvc.Result, error) {
	getUserByIDPbRequest := &pb_models.GetUserByIDRequest{
		UserID: id,
	}

	userPbResponse, err := u.userClient.GetUserByID(ctx, getUserByIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResponse, err := mappers.UserMapper.ToHttpResponse(userPbResponse)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: userResponse,
	}, nil
}

// @Summary Get authorized user
// @Description Get user by user credential
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.User	"ok"
// @Failure 404 {object} iris.Problem "User not found"
// @Router /users/me [get]
func (u UserController) Me(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	getUserByIDPbRequest := &pb_models.GetUserByIDRequest{
		UserID: userID,
	}

	userPbResponse, err := u.userClient.GetUserByID(ctx, getUserByIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userResponse, err := mappers.UserMapper.ToHttpResponse(userPbResponse)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: userResponse,
	}, nil
}
