package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/user/responses"
	pb_responses "github.com/vulpes-ferrilata/user-service-proto/pb/responses"
)

func ToUserHttpResponse(userPbResponse *pb_responses.User) *responses.User {
	if userPbResponse == nil {
		return nil
	}

	return &responses.User{
		ID:          userPbResponse.GetID(),
		DisplayName: userPbResponse.GetDisplayName(),
	}
}
