package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/models"
	pb_models "github.com/vulpes-ferrilata/chat-service-proto/pb/models"
)

type MessageMapper struct{}

func (m MessageMapper) ToHttpResponse(messagePbResponse *pb_models.Message) (*models.Message, error) {
	if messagePbResponse == nil {
		return nil, nil
	}

	return &models.Message{
		ID:     messagePbResponse.GetID(),
		UserID: messagePbResponse.GetUserID(),
		Detail: messagePbResponse.GetDetail(),
	}, nil
}
