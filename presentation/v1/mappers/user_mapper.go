package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/user"
)

func ToUserHttpResponse(userGrpcResponse *user.UserResponse) *responses.User {
	return &responses.User{
		ID:          userGrpcResponse.GetID(),
		DisplayName: userGrpcResponse.GetDisplayName(),
	}
}
