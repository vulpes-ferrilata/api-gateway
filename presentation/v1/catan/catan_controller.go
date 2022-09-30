package catan

import (
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/neffos"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/requests"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	"github.com/vulpes-ferrilata/catan-service-proto/pb"
	pb_requests "github.com/vulpes-ferrilata/catan-service-proto/pb/requests"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCatanController(catanClient pb.CatanClient,
	universalTranslator *ut.UniversalTranslator,
	websocketServer *neffos.Server) *CatanController {
	return &CatanController{
		catanClient:         catanClient,
		universalTranslator: universalTranslator,
		websocketServer:     websocketServer,
	}
}

type CatanController struct {
	catanClient         pb.CatanClient
	universalTranslator *ut.UniversalTranslator
	websocketServer     *neffos.Server
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

	findGamesByUserIDPbRequest := &pb_requests.FindGamesByUserID{
		UserID: userID,
	}

	gameListPbResponse, err := c.catanClient.FindGamesByUserID(ctx, findGamesByUserIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gameResponses, _ := slices.Map(func(gamePbResponse *pb_responses.Game) (*responses.Game, error) {
		return mappers.ToGameHttpResponse(gamePbResponse), nil
	}, gameListPbResponse.Games)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: gameResponses,
	}, nil
}

func (c CatanController) GetBy(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	getGameByIDByUserIDPbRequest := &pb_requests.GetGameByIDByUserID{
		GameID: id,
		UserID: userID,
	}

	gamePbResponse, err := c.catanClient.GetGameByIDByUserID(ctx, getGameByIDByUserIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gameResponse := mappers.ToGameHttpResponse(gamePbResponse)

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: gameResponse,
	}, nil
}

func (c CatanController) Post(ctx iris.Context) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	gameID := primitive.NewObjectID()

	createGamePbRequest := &pb_requests.CreateGame{
		GameID: gameID.Hex(),
		UserID: userID,
	}

	if _, err := c.catanClient.CreateGame(ctx, createGamePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-created-game")
	if err != nil {
		translatedDetail = "i-created-game"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Event:     "game:created",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusCreated,
		Object: &struct {
			ID string `json:"id"`
		}{
			ID: gameID.Hex(),
		},
	}, nil
}

func (c CatanController) Join(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	joinGamePbRequest := &pb_requests.JoinGame{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.JoinGame(ctx, joinGamePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-joined-game")
	if err != nil {
		translatedDetail = "i-joined-game"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) Start(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	startGamePbRequest := &pb_requests.StartGame{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.StartGame(ctx, startGamePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-started-game")
	if err != nil {
		translatedDetail = "i-started-game"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) BuildSettlementAndRoad(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	buildSettlementAndRoadRequest := &requests.BuildSettlementAndRoad{}

	if err := ctx.ReadJSON(buildSettlementAndRoadRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	buildSettlementAndRoadPbRequest := &pb_requests.BuildSettlementAndRoad{
		GameID: id,
		UserID: userID,
		LandID: buildSettlementAndRoadRequest.LandID,
		PathID: buildSettlementAndRoadRequest.PathID,
	}

	if _, err := c.catanClient.BuildSettlementAndRoad(ctx, buildSettlementAndRoadPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-built-settlement-and-road")
	if err != nil {
		translatedDetail = "i-built-settlement-and-road"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) RollDices(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	rollDicesPbRequest := &pb_requests.RollDices{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.RollDices(ctx, rollDicesPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-rolled-dices")
	if err != nil {
		translatedDetail = "i-rolled-dices"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) MoveRobber(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	moveRobberRequest := &requests.MoveRobber{}

	if err := ctx.ReadJSON(moveRobberRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	moveRobberPbRequest := &pb_requests.MoveRobber{
		GameID:    id,
		UserID:    userID,
		TerrainID: moveRobberRequest.TerrainID,
		PlayerID:  moveRobberRequest.PlayerID,
	}

	if _, err := c.catanClient.MoveRobber(ctx, moveRobberPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-moved-robber")
	if err != nil {
		translatedDetail = "i-moved-robber"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) EndTurn(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	endTurnPbRequest := &pb_requests.EndTurn{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.EndTurn(ctx, endTurnPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-ended-turn")
	if err != nil {
		translatedDetail = "i-ended-turn"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) BuildSettlement(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	buildSettlementRequest := &requests.BuildSettlement{}

	if err := ctx.ReadJSON(buildSettlementRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	buildSettlementPbRequest := &pb_requests.BuildSettlement{
		GameID: id,
		UserID: userID,
		LandID: buildSettlementRequest.LandID,
	}

	if _, err := c.catanClient.BuildSettlement(ctx, buildSettlementPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-built-settlement")
	if err != nil {
		translatedDetail = "i-built-settlement"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) BuildRoad(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	buildRoadRequest := &requests.BuildRoad{}

	if err := ctx.ReadJSON(buildRoadRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	buildRoadPbRequest := &pb_requests.BuildRoad{
		GameID: id,
		UserID: userID,
		PathID: buildRoadRequest.PathID,
	}

	if _, err := c.catanClient.BuildRoad(ctx, buildRoadPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-built-road")
	if err != nil {
		translatedDetail = "i-built-road"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) UpgradeCity(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	upgradeCityRequest := &requests.UpgradeCity{}

	if err := ctx.ReadJSON(upgradeCityRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	upgradeCityPbRequest := &pb_requests.UpgradeCity{
		GameID:         id,
		UserID:         userID,
		ConstructionID: upgradeCityRequest.ConstructionID,
	}

	if _, err := c.catanClient.UpgradeCity(ctx, upgradeCityPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-upgraded-city")
	if err != nil {
		translatedDetail = "i-upgraded-city"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) BuyDevelopmentCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	buyDevelopmentCardPbRequest := &pb_requests.BuyDevelopmentCard{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.BuyDevelopmentCard(ctx, buyDevelopmentCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-bought-development-card")
	if err != nil {
		translatedDetail = "i-bought-development-card"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) ToggleResourceCards(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	toggleResourceCardsRequest := &requests.ToggleResourceCards{}

	if err := ctx.ReadJSON(toggleResourceCardsRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	toggleResourceCardsPbRequest := &pb_requests.ToggleResourceCards{
		GameID:          id,
		UserID:          userID,
		ResourceCardIDs: toggleResourceCardsRequest.ResourceCardIDs,
	}

	if _, err := c.catanClient.ToggleResourceCards(ctx, toggleResourceCardsPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-toggled-resource-cards")
	if err != nil {
		translatedDetail = "i-toggled-resource-cards"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) MaritimeTrade(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	maritimeTradeRequest := &requests.MaritimeTrade{}

	if err := ctx.ReadJSON(maritimeTradeRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	maritimeTradePbRequest := &pb_requests.MaritimeTrade{
		GameID:           id,
		UserID:           userID,
		ResourceCardType: maritimeTradeRequest.ResourceCardType,
	}

	if _, err := c.catanClient.MaritimeTrade(ctx, maritimeTradePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-traded-with-maritime")
	if err != nil {
		translatedDetail = "i-traded-with-maritime"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) OfferTrading(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	offerTradingRequest := &requests.OfferTrading{}

	if err := ctx.ReadJSON(offerTradingRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	offerTradingPbRequest := &pb_requests.OfferTrading{
		GameID:   id,
		UserID:   userID,
		PlayerID: offerTradingRequest.PlayerID,
	}

	if _, err := c.catanClient.OfferTrading(ctx, offerTradingPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-sent-a-trade-offer")
	if err != nil {
		translatedDetail = "i-sent-a-trade-offer"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) ConfirmTrading(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	confirmTradingPbRequest := &pb_requests.ConfirmTrading{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.ConfirmTrading(ctx, confirmTradingPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-accepted-trade-offer")
	if err != nil {
		translatedDetail = "i-accepted-trade-offer"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) CancelTrading(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	cancelTradingPbRequest := &pb_requests.CancelTrading{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.CancelTrading(ctx, cancelTradingPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-cancelled-trade-offer")
	if err != nil {
		translatedDetail = "i-cancelled-trade-offer"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) PlayKnightCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playKnightCardRequest := &requests.PlayKnightCard{}

	if err := ctx.ReadJSON(playKnightCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playKnightCardPbRequest := &pb_requests.PlayKnightCard{
		GameID:    id,
		UserID:    userID,
		TerrainID: playKnightCardRequest.TerrainID,
		PlayerID:  playKnightCardRequest.PlayerID,
	}

	if _, err := c.catanClient.PlayKnightCard(ctx, playKnightCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-played-knight-card")
	if err != nil {
		translatedDetail = "i-played-knight-card"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) PlayRoadBuildingCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playRoadBuildingCardRequest := &requests.PlayRoadBuildingCard{}

	if err := ctx.ReadJSON(playRoadBuildingCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playRoadBuildingCardPbRequest := &pb_requests.PlayRoadBuildingCard{
		GameID:  id,
		UserID:  userID,
		PathIDs: playRoadBuildingCardRequest.PathIDs,
	}

	if _, err := c.catanClient.PlayRoadBuildingCard(ctx, playRoadBuildingCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-played-road-building-card")
	if err != nil {
		translatedDetail = "i-played-road-building-card"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) PlayYearOfPlentyCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playYearOfPlentyCardRequest := &requests.PlayYearOfPlentyCard{}

	if err := ctx.ReadJSON(playYearOfPlentyCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playYearOfPlentyCardPbRequest := &pb_requests.PlayYearOfPlentyCard{
		GameID:            id,
		UserID:            userID,
		ResourceCardTypes: playYearOfPlentyCardRequest.ResourceCardTypes,
	}

	if _, err := c.catanClient.PlayYearOfPlentyCard(ctx, playYearOfPlentyCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-played-year-of-plenty-card")
	if err != nil {
		translatedDetail = "i-played-year-of-plenty-card"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

func (c CatanController) PlayMonopolyCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playMonopolyCardRequest := &requests.PlayMonopolyCard{}

	if err := ctx.ReadJSON(playMonopolyCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playMonopolyCardPbRequest := &pb_requests.PlayMonopolyCard{
		GameID:           id,
		UserID:           userID,
		ResourceCardType: playMonopolyCardRequest.ResourceCardType,
	}

	if _, err := c.catanClient.PlayMonopolyCard(ctx, playMonopolyCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	locales := context_values.GetLocales(ctx)
	translator, _ := c.universalTranslator.FindTranslator(locales...)
	translatedDetail, err := translator.T("i-played-monopoly-card")
	if err != nil {
		translatedDetail = "i-played-monopoly-card"
	}

	notification := &responses.Notification{
		UserID: userID,
		Detail: translatedDetail,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "catan",
		Room:      id,
		Event:     "game:updated",
		Body:      neffos.Marshal(notification),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}
