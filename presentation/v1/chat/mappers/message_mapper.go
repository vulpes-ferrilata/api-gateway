package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/responses"
	pb_responses "github.com/vulpes-ferrilata/chat-service-proto/pb/responses"
)

var MessageMapper messageMapper = messageMapper{}

type messageMapper struct{}

func (m messageMapper) ToHttpResponse(messagePbResponse *pb_responses.Message) (*responses.Message, error) {
	if messagePbResponse == nil {
		return nil, nil
	}

	return &responses.Message{
		ID:     messagePbResponse.GetID(),
		UserID: messagePbResponse.GetUserID(),
		Detail: messagePbResponse.GetDetail(),
	}, nil
}
