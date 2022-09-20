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
	b.Handle(http.MethodPost, "/{id:string}/move-robber", "MoveRobber")
	b.Handle(http.MethodPost, "/{id:string}/end-turn", "EndTurn")
	b.Handle(http.MethodPost, "/{id:string}/build-settlement", "BuildSettlement")
	b.Handle(http.MethodPost, "/{id:string}/build-road", "BuildRoad")
	b.Handle(http.MethodPost, "/{id:string}/upgrade-city", "UpgradeCity")
	b.Handle(http.MethodPost, "/{id:string}/buy-development-card", "BuyDevelopmentCard")
	b.Handle(http.MethodPost, "/{id:string}/toggle-resource-cards", "ToggleResourceCards")
	b.Handle(http.MethodPost, "/{id:string}/maritime-trade", "MaritimeTrade")
	b.Handle(http.MethodPost, "/{id:string}/offer-trading", "OfferTrading")
	b.Handle(http.MethodPost, "/{id:string}/confirm-trading", "ConfirmTrading")
	b.Handle(http.MethodPost, "/{id:string}/cancel-trading", "CancelTrading")
	b.Handle(http.MethodPost, "/{id:string}/play-knight-card", "PlayKnightCard")
	b.Handle(http.MethodPost, "/{id:string}/play-road-building-card", "PlayRoadBuildingCard")
	b.Handle(http.MethodPost, "/{id:string}/play-year-of-plenty-card", "PlayYearOfPlentyCard")
	b.Handle(http.MethodPost, "/{id:string}/play-monopoly-card", "PlayMonopolyCard")
}

func (c CatanController) Get(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	findGamesByUserIDGrpcRequest := &catan.FindGamesByUserIDRequest{
		UserID: userID,
	}

	gameListGrpcResponse, err := c.catanClient.FindGamesByUserID(ctx, findGamesByUserIDGrpcRequest)
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

	getGameByIDByUserIDGrpcRequest := &catan.GetGameByIDByUserIDRequest{
		UserID: userID,
		GameID: id,
	}

	gameGrpcResponse, err := c.catanClient.GetGameByIDByUserID(ctx, getGameByIDByUserIDGrpcRequest)
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

	rollDicesGrpcRequest := &catan.RollDicesRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.RollDices(ctx, rollDicesGrpcRequest); err != nil {
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

func (c CatanController) MoveRobber(ctx iris.Context, id string) (mvc.Result, error) {
	moveRobberRequest := &requests.MoveRobber{}

	if err := ctx.ReadJSON(moveRobberRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	moveRobberGrpcRequest := &catan.MoveRobberRequest{
		UserID:    userID,
		GameID:    id,
		TerrainID: moveRobberRequest.TerrainID,
		PlayerID:  moveRobberRequest.PlayerID,
	}

	if _, err := c.catanClient.MoveRobber(ctx, moveRobberGrpcRequest); err != nil {
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

func (c CatanController) EndTurn(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	endTurnGrpcRequest := &catan.EndTurnRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.EndTurn(ctx, endTurnGrpcRequest); err != nil {
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

func (c CatanController) BuildSettlement(ctx iris.Context, id string) (mvc.Result, error) {
	buildSettlementRequest := &requests.BuildSettlement{}

	if err := ctx.ReadJSON(buildSettlementRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	buildSettlementGrpcRequest := &catan.BuildSettlementRequest{
		UserID: userID,
		GameID: id,
		LandID: buildSettlementRequest.LandID,
	}

	if _, err := c.catanClient.BuildSettlement(ctx, buildSettlementGrpcRequest); err != nil {
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

func (c CatanController) BuildRoad(ctx iris.Context, id string) (mvc.Result, error) {
	buildRoadRequest := &requests.BuildRoad{}

	if err := ctx.ReadJSON(buildRoadRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	buildRoadGrpcRequest := &catan.BuildRoadRequest{
		UserID: userID,
		GameID: id,
		PathID: buildRoadRequest.PathID,
	}

	if _, err := c.catanClient.BuildRoad(ctx, buildRoadGrpcRequest); err != nil {
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

func (c CatanController) UpgradeCity(ctx iris.Context, id string) (mvc.Result, error) {
	upgradeCityRequest := &requests.UpgradeCity{}

	if err := ctx.ReadJSON(upgradeCityRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	upgradeCityGrpcRequest := &catan.UpgradeCityRequest{
		UserID:         userID,
		GameID:         id,
		ConstructionID: upgradeCityRequest.ConstructionID,
	}

	if _, err := c.catanClient.UpgradeCity(ctx, upgradeCityGrpcRequest); err != nil {
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

func (c CatanController) BuyDevelopmentCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	buyDevelopmentCardGrpcRequest := &catan.BuyDevelopmentCardRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.BuyDevelopmentCard(ctx, buyDevelopmentCardGrpcRequest); err != nil {
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

func (c CatanController) ToggleResourceCards(ctx iris.Context, id string) (mvc.Result, error) {
	toggleResourceCardsRequest := &requests.ToggleResourceCards{}

	if err := ctx.ReadJSON(toggleResourceCardsRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	toggleResourceCardsGrpcRequest := &catan.ToggleResourceCardsRequest{
		UserID:          userID,
		GameID:          id,
		ResourceCardIDs: toggleResourceCardsRequest.ResourceCardIDs,
	}

	if _, err := c.catanClient.ToggleResourceCards(ctx, toggleResourceCardsGrpcRequest); err != nil {
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

func (c CatanController) MaritimeTrade(ctx iris.Context, id string) (mvc.Result, error) {
	maritimeTradeRequest := &requests.MaritimeTrade{}

	if err := ctx.ReadJSON(maritimeTradeRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	maritimeTradeGrpcRequest := &catan.MaritimeTradeRequest{
		UserID:           userID,
		GameID:           id,
		ResourceCardType: maritimeTradeRequest.ResourceCardType,
	}

	if _, err := c.catanClient.MaritimeTrade(ctx, maritimeTradeGrpcRequest); err != nil {
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

func (c CatanController) OfferTrading(ctx iris.Context, id string) (mvc.Result, error) {
	offerTradingRequest := &requests.OfferTrading{}

	if err := ctx.ReadJSON(offerTradingRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	offerTradingGrpcRequest := &catan.OfferTradingRequest{
		UserID:   userID,
		GameID:   id,
		PlayerID: offerTradingRequest.PlayerID,
	}

	if _, err := c.catanClient.OfferTrading(ctx, offerTradingGrpcRequest); err != nil {
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

func (c CatanController) ConfirmTrading(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	confirmTradingGrpcRequest := &catan.ConfirmTradingRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.ConfirmTrading(ctx, confirmTradingGrpcRequest); err != nil {
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

func (c CatanController) CancelTrading(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx.Request().Context())

	cancelTradingGrpcRequest := &catan.CancelTradingRequest{
		UserID: userID,
		GameID: id,
	}

	if _, err := c.catanClient.CancelTrading(ctx, cancelTradingGrpcRequest); err != nil {
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

func (c CatanController) PlayKnightCard(ctx iris.Context, id string) (mvc.Result, error) {
	playKnightCardRequest := &requests.PlayKnightCard{}

	if err := ctx.ReadJSON(playKnightCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	playKnightCardGrpcRequest := &catan.PlayKnightCardRequest{
		UserID:    userID,
		GameID:    id,
		TerrainID: playKnightCardRequest.TerrainID,
		PlayerID:  playKnightCardRequest.PlayerID,
	}

	if _, err := c.catanClient.PlayKnightCard(ctx, playKnightCardGrpcRequest); err != nil {
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

func (c CatanController) PlayRoadBuildingCard(ctx iris.Context, id string) (mvc.Result, error) {
	playRoadBuildingCardRequest := &requests.PlayRoadBuildingCard{}

	if err := ctx.ReadJSON(playRoadBuildingCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	playRoadBuildingCardGrpcRequest := &catan.PlayRoadBuildingCardRequest{
		UserID:  userID,
		GameID:  id,
		PathIDs: playRoadBuildingCardRequest.PathIDs,
	}

	if _, err := c.catanClient.PlayRoadBuildingCard(ctx, playRoadBuildingCardGrpcRequest); err != nil {
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

func (c CatanController) PlayYearOfPlentyCard(ctx iris.Context, id string) (mvc.Result, error) {
	playYearOfPlentyCardRequest := &requests.PlayYearOfPlentyCard{}

	if err := ctx.ReadJSON(playYearOfPlentyCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	playYearOfPlentyCardGrpcRequest := &catan.PlayYearOfPlentyCardRequest{
		UserID:            userID,
		GameID:            id,
		ResourceCardTypes: playYearOfPlentyCardRequest.ResourceCardTypes,
	}

	if _, err := c.catanClient.PlayYearOfPlentyCard(ctx, playYearOfPlentyCardGrpcRequest); err != nil {
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

func (c CatanController) PlayMonopolyCard(ctx iris.Context, id string) (mvc.Result, error) {
	playMonopolyCardRequest := &requests.PlayMonopolyCard{}

	if err := ctx.ReadJSON(playMonopolyCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	userID := context_values.GetUserID(ctx.Request().Context())

	playMonopolyCardGrpcRequest := &catan.PlayMonopolyCardRequest{
		UserID:           userID,
		GameID:           id,
		ResourceCardType: playMonopolyCardRequest.ResourceCardType,
	}

	if _, err := c.catanClient.PlayMonopolyCard(ctx, playMonopolyCardGrpcRequest); err != nil {
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
