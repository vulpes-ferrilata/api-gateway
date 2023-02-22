package chat

import (
	"github.com/gocql/gocql"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/neffos"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat/models"
	"github.com/vulpes-ferrilata/chat-service-proto/pb"
	pb_models "github.com/vulpes-ferrilata/chat-service-proto/pb/models"
	"github.com/vulpes-ferrilata/slices"
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

// @Summary Find messages
// @Description Find messages by room id
// @Accept  json
// @Produce  json
// @Success 200 {array} responses.Message	"ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Router /chat/messages [get]
func (c ChatController) Get(ctx iris.Context) (mvc.Result, error) {
	roomID := ctx.URLParam("roomID")

	findMessagesByRoomIDPbRequest := &pb_models.FindMessagesByRoomIDRequest{
		RoomID: roomID,
	}

	messageListPbResponse, err := c.chatClient.FindMessagesByRoomID(ctx, findMessagesByRoomIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponses, err := slices.Map(func(messagePbResponse *pb_models.Message) (*models.Message, error) {
		return mappers.MessageMapper{}.ToHttpResponse(messagePbResponse)
	}, messageListPbResponse.Messages...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: messageResponses,
	}, nil
}

// @Summary Create message
// @Description Create new message
// @Accept  json
// @Produce  json
// @Success 201 {object} responses.Message "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Router /chat/messages [post]
func (c ChatController) Post(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	createMessageRequest := &models.CreateMessageRequest{}

	if err := ctx.ReadJSON(createMessageRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageID := gocql.TimeUUID()

	createMessagePbRequest := &pb_models.CreateMessageRequest{
		MessageID: messageID.String(),
		RoomID:    createMessageRequest.RoomID,
		UserID:    userID,
		Detail:    createMessageRequest.Detail,
	}

	if _, err := c.chatClient.CreateMessage(ctx, createMessagePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &models.Message{
		ID:     messageID.String(),
		UserID: userID,
		Detail: createMessageRequest.Detail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Chat",
		Room:      createMessageRequest.RoomID,
		Event:     "MessageCreated",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: messageResponse,
	}, nil
}
