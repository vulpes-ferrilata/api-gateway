package chat

import (
	"github.com/gocql/gocql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/neffos"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/requests"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/responses"
	"github.com/vulpes-ferrilata/chat-service-proto/pb"
	pb_requests "github.com/vulpes-ferrilata/chat-service-proto/pb/requests"
	pb_responses "github.com/vulpes-ferrilata/chat-service-proto/pb/responses"
)

func NewChatController(chatClient pb.ChatClient, websocketServer *neffos.Server) *ChatController {
	return &ChatController{
		chatClient:      chatClient,
		websocketServer: websocketServer,
	}
}

type ChatController struct {
	chatClient      pb.ChatClient
	websocketServer *neffos.Server
}

func (c ChatController) Get(ctx iris.Context) (mvc.Result, error) {
	roomID := ctx.URLParam("roomID")

	findMessagesByRoomIDPbRequest := &pb_requests.FindMessagesByRoomID{
		RoomID: roomID,
	}

	messageListPbResponse, err := c.chatClient.FindMessagesByRoomID(ctx, findMessagesByRoomIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponses, err := slices.Map(func(messagePbResponse *pb_responses.Message) (*responses.Message, error) {
		return mappers.ToMessageHttpResponse(messagePbResponse), nil
	}, messageListPbResponse.Messages)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: messageResponses,
	}, nil
}

func (c ChatController) Post(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	messageRequest := &requests.Message{}

	if err := ctx.ReadJSON(messageRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageID := gocql.TimeUUID()

	messageResponse := &responses.Message{
		ID:     messageID.String(),
		RoomID: messageRequest.RoomID,
		UserID: userID,
		Detail: messageRequest.Detail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "chat",
		Room:      messageRequest.RoomID,
		Event:     "message:created",
		Body:      neffos.Marshal(messageResponse),
	})

	//eventually consistence
	createMessagePbRequest := &pb_requests.CreateMessage{
		MessageID: messageID.String(),
		RoomID:    messageRequest.RoomID,
		UserID:    userID,
		Detail:    messageRequest.Detail,
	}

	if _, err := c.chatClient.CreateMessage(ctx, createMessagePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}
