package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/user/responses"
	pb_responses "github.com/vulpes-ferrilata/user-service-proto/pb/responses"
)

var UserMapper userMapper = userMapper{}

type userMapper struct{}

func (u userMapper) ToHttpResponse(userPbResponse *pb_responses.User) (*responses.User, error) {
	if userPbResponse == nil {
		return nil, nil
	}

	return &responses.User{
		ID:          userPbResponse.GetID(),
		DisplayName: userPbResponse.GetDisplayName(),
	}, nil
}
