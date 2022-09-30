package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/responses"
	pb_responses "github.com/vulpes-ferrilata/chat-service-proto/pb/responses"
)

func ToMessageHttpResponse(messagePbResponse *pb_responses.Message) *responses.Message {
	if messagePbResponse == nil {
		return nil
	}

	return &responses.Message{
		ID:     messagePbResponse.GetID(),
		RoomID: messagePbResponse.GetRoomID(),
		UserID: messagePbResponse.GetUserID(),
		Detail: messagePbResponse.GetDetail(),
	}
}
