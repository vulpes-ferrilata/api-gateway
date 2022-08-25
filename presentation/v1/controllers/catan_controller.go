package controllers

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/neffos"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/requests"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/responses"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCatanController(catanClient catan.CatanClient,
	websocketServer *neffos.Server) *CatanController {
	return &CatanController{
		catanClient:     catanClient,
		websocketServer: websocketServer,
	}
}

type CatanController struct {
	catanClient     catan.CatanClient
	websocketServer *neffos.Server
}

func (c CatanController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodPost, "/{id:string}/join", "Join")
	b.Handle(http.MethodPost, "/{id:string}/start", "Start")
	b.Handle(http.MethodPost, "/{id:string}/build-settlement-and-road", "BuildSettlementAndRoad")
	b.Handle(http.MethodPost, "/{id:string}/roll-dices", "RollDices")
}

func (c CatanController) Get(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	findGamesByUserIDRequest := &catan.FindGamesByUserIDRequest{
		UserID: userID,
	}

	gameListGrpcResponse, err := c.catanClient.FindGamesByUserID(ctx, findGamesByUserIDRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gameResponses, _ := slices.Map(func(gameGrpcResponse *catan.GameResponse) (*responses.Game, error) {
		return mappers.ToGameHttpResponse(gameGrpcResponse), nil
	}, gameListGrpcResponse.Games)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: gameResponses,
	}, nil
}

func (c CatanController) GetBy(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	getGameByIDByUserIDRequest := &catan.GetGameByIDByUserIDRequest{
		UserID: userID,
		GameID: id,
	}

	gameGrpcResponse, err := c.catanClient.GetGameByIDByUserID(ctx, getGameByIDByUserIDRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gameResponse := mappers.ToGameHttpResponse(gameGrpcResponse)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: gameResponse,
	}, nil
}

func (c CatanController) Post(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	gameID := primitive.NewObjectID()

	createGameGrpcRequest := &catan.CreateGameRequest{
		UserID: userID,
		GameID: gameID.Hex(),
	}

	if _, err := c.catanClient.CreateGame(ctx, createGameGrpcRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Event:     "game:created",
	})

	gameCreatedResponse := &responses.GameCreated{
		ID: gameID.Hex(),
	}

	return &mvc.Response{
		Code:   iris.StatusCreated,
		Object: gameCreatedResponse,
	}, nil
}

func (c CatanController) Join(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	joinGameGrpcRequest := &catan.JoinGameRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.JoinGame(ctx, joinGameGrpcRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) Start(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	startGameGrpcRequest := &catan.StartGameRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.StartGame(ctx, startGameGrpcRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) BuildSettlementAndRoad(ctx iris.Context, id string) (mvc.Result, error) {
	buildSettlementAndRoadRequest := &requests.BuildSettlementAndRoad{}

	if err := ctx.ReadJSON(buildSettlementAndRoadRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	buildSettlementAndRoadGrpcRequest := &catan.BuildSettlementAndRoadRequest{
		UserID: userID,
		GameID: id,
		LandID: buildSettlementAndRoadRequest.LandID,
		PathID: buildSettlementAndRoadRequest.PathID,
	}

	if _, err := c.catanClient.BuildSettlementAndRoad(ctx, buildSettlementAndRoadGrpcRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) RollDices(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	rollDicesRequest := &catan.RollDicesRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.RollDices(ctx, rollDicesRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}
