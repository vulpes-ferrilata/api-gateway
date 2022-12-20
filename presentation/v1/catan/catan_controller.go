package catan

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/neffos"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/context_values"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/mappers"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/requests"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	"github.com/vulpes-ferrilata/catan-service-proto/pb"
	pb_requests "github.com/vulpes-ferrilata/catan-service-proto/pb/requests"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCatanController(catanClient pb.CatanClient,
	websocketServer *neffos.Server) *CatanController {
	return &CatanController{
		catanClient:     catanClient,
		websocketServer: websocketServer,
	}
}

type CatanController struct {
	catanClient     pb.CatanClient
	websocketServer *neffos.Server
}

func (c CatanController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle(http.MethodPost, "/{id:string}/join", "Join")
	b.Handle(http.MethodPost, "/{id:string}/start", "Start")
	b.Handle(http.MethodPost, "/{id:string}/build-settlement-and-road", "BuildSettlementAndRoad")
	b.Handle(http.MethodPost, "/{id:string}/roll-dices", "RollDices")
	b.Handle(http.MethodPost, "/{id:string}/discard-resource-cards", "DiscardResourceCards")
	b.Handle(http.MethodPost, "/{id:string}/move-robber", "MoveRobber")
	b.Handle(http.MethodPost, "/{id:string}/end-turn", "EndTurn")
	b.Handle(http.MethodPost, "/{id:string}/build-settlement", "BuildSettlement")
	b.Handle(http.MethodPost, "/{id:string}/build-road", "BuildRoad")
	b.Handle(http.MethodPost, "/{id:string}/upgrade-city", "UpgradeCity")
	b.Handle(http.MethodPost, "/{id:string}/buy-development-card", "BuyDevelopmentCard")
	b.Handle(http.MethodPost, "/{id:string}/toggle-resource-cards", "ToggleResourceCards")
	b.Handle(http.MethodPost, "/{id:string}/maritime-trade", "MaritimeTrade")
	b.Handle(http.MethodPost, "/{id:string}/send-trade-offer", "SendTradeOffer")
	b.Handle(http.MethodPost, "/{id:string}/confirm-trade-offer", "ConfirmTradeOffer")
	b.Handle(http.MethodPost, "/{id:string}/cancel-trade-offer", "CancelTradeOffer")
	b.Handle(http.MethodPost, "/{id:string}/play-knight-card", "PlayKnightCard")
	b.Handle(http.MethodPost, "/{id:string}/play-road-building-card", "PlayRoadBuildingCard")
	b.Handle(http.MethodPost, "/{id:string}/play-year-of-plenty-card", "PlayYearOfPlentyCard")
	b.Handle(http.MethodPost, "/{id:string}/play-monopoly-card", "PlayMonopolyCard")
	b.Handle(http.MethodPost, "/{id:string}/play-victory-point-card", "PlayVictoryPointCard")
}

// @Summary Get game pagination
// @Description Find games by limit by offset
// @Accept  json
// @Produce  json
// @Param	limit	   query    int    false	"Limit"
// @Param	offset	   query    int	   false	"Offset"
// @Success 200 {object} responses.GamePagination "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Router /catan/games [get]
func (c CatanController) Get(ctx iris.Context) (mvc.Result, error) {
	limit := ctx.URLParamIntDefault("limit", 0)
	offset := ctx.URLParamIntDefault("offset", 0)

	findGamePaginationByLimitByOffsetPbRequest := &pb_requests.FindGamePaginationByLimitByOffset{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	gamePaginationPbResponse, err := c.catanClient.FindGamePaginationByLimitByOffset(ctx, findGamePaginationByLimitByOffsetPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gameResponses, err := mappers.GamePaginationMapper.ToHttpResponse(gamePaginationPbResponse)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: gameResponses,
	}, nil
}

// @Summary Get game
// @Description Get game by id
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {object} responses.GameDetail "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "game not found"
// @Router /catan/games/{id} [get]
func (c CatanController) GetBy(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	getGameDetailByIDByUserIDPbRequest := &pb_requests.GetGameDetailByIDByUserID{
		GameID: id,
		UserID: userID,
	}

	gameDetailPbResponse, err := c.catanClient.GetGameDetailByIDByUserID(ctx, getGameDetailByIDByUserIDPbRequest)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	gameDetailResponse, err := mappers.GameDetailMapper.ToHttpResponse(gameDetailPbResponse)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &mvc.Response{
		Code:   iris.StatusOK,
		Object: gameDetailResponse,
	}, nil
}

// @Summary Create game
// @Description Create new game
// @Accept  json
// @Produce  json
// @Success 200 {object} responses.GameDetail "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Router /catan/games/ [post]
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

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Event:     "GameCreated",
		Body:      neffos.Marshal(messageResponse),
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

// @Summary Join game
// @Description Join game at waiting state
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 422 {object} iris.Problem "game has full players"
// @Failure 422 {object} iris.Problem "game already started"
// @Failure 422 {object} iris.Problem "game already finished"
// @Router /catan/games/{id}/join [post]
func (c CatanController) Join(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	joinGamePbRequest := &pb_requests.JoinGame{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.JoinGame(ctx, joinGamePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "GameJoined",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Start game
// @Description Start game at waiting state
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "game must have at least 2 players"
// @Failure 422 {object} iris.Problem "game already started"
// @Failure 422 {object} iris.Problem "game already finished"
// @Router /catan/games/{id}/start [post]
func (c CatanController) Start(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	startGamePbRequest := &pb_requests.StartGame{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.StartGame(ctx, startGamePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "GameStarted",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Build settlement and road
// @Description Build settlement and road at setup phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	landID	   body    requests.BuildSettlementAndRoad	true	"Land ID"
// @Param	pathID	   body    requests.BuildSettlementAndRoad	true	"Path ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "land not found"
// @Failure 404 {object} iris.Problem "path not found"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "nearby lands must be vacant"
// @Failure 422 {object} iris.Problem "selected land and path must be adjacent"
// @Failure 422 {object} iris.Problem "you run out of settlements"
// @Failure 422 {object} iris.Problem "you run out of roads"
// @Router /catan/games/{id}/build-settlement-and-road [post]
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

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "SettlementAndRoadBuilt",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Roll dices
// @Description Roll dices at resource production phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource consumption phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Router /catan/games/{id}/roll-dices [post]
func (c CatanController) RollDices(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	rollDicesPbRequest := &pb_requests.RollDices{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.RollDices(ctx, rollDicesPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "DicesRolled",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Discard resource cards
// @Description Discard resource cards by half when handling more than 7 resource cards at resource discard phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	resourceCardIDs	   body    requests.BuildSettlementAndRoad	true	"List of Resource Card ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource consumption phase"
// @Failure 422 {object} iris.Problem "you already discarded resource cards"
// @Failure 422 {object} iris.Problem "you have no need to discard resource cards"
// @Failure 422 {object} iris.Problem "selected resource cards must be equals to half of your currently resource cards"
// @Router /catan/games/{id}/discard-resource-cards [post]
func (c CatanController) DiscardResourceCards(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	discardResourceCardsRequest := &requests.DiscardResourceCards{}

	if err := ctx.ReadJSON(discardResourceCardsRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	discardResourceCardsPbRequest := &pb_requests.DiscardResourceCards{
		GameID:          id,
		UserID:          userID,
		ResourceCardIDs: discardResourceCardsRequest.ResourceCardIDs,
	}

	if _, err := c.catanClient.DiscardResourceCards(ctx, discardResourceCardsPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "ResourceCardsDiscarded",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Move robber
// @Description Move robber and steal resource card if robber placed on terrain which has enemy construction nearby at robbing phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	terrainID	   body    requests.MoveRobber	true	"Terrain ID"
// @Param	playerID	   body    requests.MoveRobber	false	"Player ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "terrain not found"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource consumption phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "robber must be moved to other terrain"
// @Failure 422 {object} iris.Problem "you must rob player who has construction next to robber"
// @Failure 422 {object} iris.Problem "selected player must have construction next to robber"
// @Router /catan/games/{id}/move-robber [post]
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

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "RobberMoved",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary End turn
// @Description End current turn at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Router /catan/games/{id}/end-turn [post]
func (c CatanController) EndTurn(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	endTurnPbRequest := &pb_requests.EndTurn{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.EndTurn(ctx, endTurnPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "TurnEnded",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Build settlement
// @Description build settlement by using resource cards at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	landID	   body    requests.BuildSettlement	true	"Land ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "land not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "you have insufficient resource cards"
// @Failure 422 {object} iris.Problem "nearby lands must be vacant"
// @Failure 422 {object} iris.Problem "selected land must be adjacent to your road"
// @Failure 422 {object} iris.Problem "you run out of settlements"
// @Router /catan/games/{id}/build-settlement [post]
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

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "SettlementBuilt",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Build road
// @Description build road by using resource cards at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	pathID	   body    requests.BuildRoad	true	"Path ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "path not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "you have insufficient resource cards"
// @Failure 422 {object} iris.Problem "selected path must be adjacent to your construction or road"
// @Failure 422 {object} iris.Problem "selected path pass through construction of other player"
// @Failure 422 {object} iris.Problem "you run out of roads"
// @Router /catan/games/{id}/build-road [post]
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

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "RoadBuilt",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Upgrade city
// @Description upgrade your settlement to city by using resource cards at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	constructionID	   body    requests.UpgradeCity	true	"Construction ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "construction not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "you have insufficient resource cards"
// @Failure 422 {object} iris.Problem "selected construction already upgraded"
// @Failure 422 {object} iris.Problem "selected construction does not belong to any land"
// @Failure 422 {object} iris.Problem "you run out of cities"
// @Router /catan/games/{id}/upgrade-city [post]
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

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "CityUpgraded",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Buy development
// @Description Buy development by using resource cards at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "game run out of development cards"
// @Router /catan/games/{id}/buy-development-card [post]
func (c CatanController) BuyDevelopmentCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	buyDevelopmentCardPbRequest := &pb_requests.BuyDevelopmentCard{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.BuyDevelopmentCard(ctx, buyDevelopmentCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "DevelopmentCardBought",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Toggle resource cards
// @Description Turn selected resource card into offer/unoffer at resource consumption phase, the offering resource card will be showed up to other players and it will be used to trade with maritime or other players
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	resourceCardIDs	   body    requests.ToggleResourceCards	true	"List of Resource Card ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 404 {object} iris.Problem "resource card not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Router /catan/games/{id}/toggle-resource-cards [post]
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

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "ResourceCardsToggled",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Maritime trade
// @Description Exchange your offering resource cards with selected one on the table at resource consumption phase, all of the offering resource cards will be exchange with the lowest ratio in case of owning harbors
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	resourceCardType	   body    requests.MaritimeTrade	true	"Resource Card Type"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "game has insufficient resource cards"
// @Router /catan/games/{id}/maritime-trade [post]
func (c CatanController) MaritimeTrade(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	maritimeTradeRequest := &requests.MaritimeTrade{}

	if err := ctx.ReadJSON(maritimeTradeRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	maritimeTradePbRequest := &pb_requests.MaritimeTrade{
		GameID:                    id,
		UserID:                    userID,
		ResourceCardType:          maritimeTradeRequest.ResourceCardType,
		DemandingResourceCardType: maritimeTradeRequest.DemandingResourceCardType,
	}

	if _, err := c.catanClient.MaritimeTrade(ctx, maritimeTradePbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "MaritimeTraded",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Send trade offer
// @Description Offer other player to exchange their offering resource cards at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	playerID	   body    requests.SendTradeOffer	true	"Player ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "you already offered this player"
// @Failure 422 {object} iris.Problem "you must offer at least one resource card"
// @Failure 422 {object} iris.Problem "selected player must offer at least one resource card"
// @Router /catan/games/{id}/send-trade-offer [post]
func (c CatanController) SendTradeOffer(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	sendTradeOfferRequest := &requests.SendTradeOffer{}

	if err := ctx.ReadJSON(sendTradeOfferRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	sendTradeOfferPbRequest := &pb_requests.SendTradeOffer{
		GameID:   id,
		UserID:   userID,
		PlayerID: sendTradeOfferRequest.PlayerID,
	}

	if _, err := c.catanClient.SendTradeOffer(ctx, sendTradeOfferPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "TradeOfferSent",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Confirm trade offer
// @Description Confirm exchanging offering resource cards with active player at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you have not received any offer"
// @Failure 422 {object} iris.Problem "you must offer at least one resource card"
// @Failure 422 {object} iris.Problem "active player must offer at least one resource card"
// @Router /catan/games/{id}/confirm-trade-offer [post]
func (c CatanController) ConfirmTradeOffer(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	confirmTradeOfferPbRequest := &pb_requests.ConfirmTradeOffer{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.ConfirmTradeOffer(ctx, confirmTradeOfferPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "TradeOfferConfirmed",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Cancel trade offer
// @Description Cancel trade offer of active player at resource consumption phase
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in setup phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource production phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in resource discard phase"
// @Failure 422 {object} iris.Problem "you are unable to perform this action in robbing phase"
// @Failure 422 {object} iris.Problem "you have not received any offer"
// @Router /catan/games/{id}/cancel-trade-offer [post]
func (c CatanController) CancelTradeOffer(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)

	cancelTradeOfferPbRequest := &pb_requests.CancelTradeOffer{
		GameID: id,
		UserID: userID,
	}

	if _, err := c.catanClient.CancelTradeOffer(ctx, cancelTradeOfferPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "TradeOfferCancelled",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Play knight card
// @Description Play knight development card from your stack at any phase of started state
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	developmentCardID	   body    requests.PlayKnightCard	true	"Development Card ID"
// @Param	terrainID	   body    requests.PlayKnightCard	true	"Terrain ID"
// @Param	playerID	   body    requests.PlayKnightCard	false	"Player ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "development card not found"
// @Failure 404 {object} iris.Problem "player not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "selected development card must be knight card"
// @Failure 422 {object} iris.Problem "selected development card is unavailable to use"
// @Failure 422 {object} iris.Problem "robber must be moved to other terrain"
// @Failure 422 {object} iris.Problem "you must rob player who has construction next to robber"
// @Failure 422 {object} iris.Problem "selected player must have construction next to robber"
// @Router /catan/games/{id}/play-knight-card [post]
func (c CatanController) PlayKnightCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playKnightCardRequest := &requests.PlayKnightCard{}

	if err := ctx.ReadJSON(playKnightCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playKnightCardPbRequest := &pb_requests.PlayKnightCard{
		GameID:            id,
		UserID:            userID,
		DevelopmentCardID: playKnightCardRequest.DevelopmentCardID,
		TerrainID:         playKnightCardRequest.TerrainID,
		PlayerID:          playKnightCardRequest.PlayerID,
	}

	if _, err := c.catanClient.PlayKnightCard(ctx, playKnightCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "KnightCardPlayed",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Play road building card
// @Description Play road building development card from your stack at any phase of started state
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	developmentCardID	   body    requests.PlayRoadBuildingCard	true	"Development Card ID"
// @Param	pathIDs	   body    requests.PlayRoadBuildingCard	true	"List of Path ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "development card not found"
// @Failure 404 {object} iris.Problem "path not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "selected development card must be road building card"
// @Failure 422 {object} iris.Problem "selected development card is unavailable to use"
// @Failure 422 {object} iris.Problem "selected path must be adjacent to your construction or road"
// @Failure 422 {object} iris.Problem "selected path pass through construction of other player"
// @Failure 422 {object} iris.Problem "you run out of roads"
// @Router /catan/games/{id}/play-road-building-card [post]
func (c CatanController) PlayRoadBuildingCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playRoadBuildingCardRequest := &requests.PlayRoadBuildingCard{}

	if err := ctx.ReadJSON(playRoadBuildingCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playRoadBuildingCardPbRequest := &pb_requests.PlayRoadBuildingCard{
		GameID:            id,
		UserID:            userID,
		DevelopmentCardID: playRoadBuildingCardRequest.DevelopmentCardID,
		PathIDs:           playRoadBuildingCardRequest.PathIDs,
	}

	if _, err := c.catanClient.PlayRoadBuildingCard(ctx, playRoadBuildingCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "RoadBuildingCardPlayed",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Play year of plenty card
// @Description Play year of plenty development card from your stack at any phase of started state
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	developmentCardID	   body    requests.PlayYearOfPlentyCard	true	"Development Card ID"
// @Param	resourceCardTypes	   body    requests.PlayYearOfPlentyCard	true	"List of Resource Card Type"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "development card not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "selected development card must be year of plenty card"
// @Failure 422 {object} iris.Problem "selected development card is unavailable to use"
// @Failure 422 {object} iris.Problem "selected path must be adjacent to your construction or road"
// @Failure 422 {object} iris.Problem "selected path pass through construction of other player"
// @Failure 422 {object} iris.Problem "game has insufficient resource cards"
// @Router /catan/games/{id}/play-year-of-plenty-card [post]
func (c CatanController) PlayYearOfPlentyCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playYearOfPlentyCardRequest := &requests.PlayYearOfPlentyCard{}

	if err := ctx.ReadJSON(playYearOfPlentyCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playYearOfPlentyCardPbRequest := &pb_requests.PlayYearOfPlentyCard{
		GameID:                     id,
		UserID:                     userID,
		DevelopmentCardID:          playYearOfPlentyCardRequest.DevelopmentCardID,
		DemandingResourceCardTypes: playYearOfPlentyCardRequest.DemandingResourceCardTypes,
	}

	if _, err := c.catanClient.PlayYearOfPlentyCard(ctx, playYearOfPlentyCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "YearOfPlentyCardPlayed",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Play monopoly card
// @Description Play monopoly development card from your stack at any phase of started state
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	developmentCardID	   body    requests.PlayMonopolyCard	true	"Development Card ID"
// @Param	resourceCardType	   body    requests.PlayMonopolyCard	true	"Resource Card Type"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "development card not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "selected development card must be monopoly card"
// @Failure 422 {object} iris.Problem "selected development card is unavailable to use"
// @Failure 422 {object} iris.Problem "robber must be moved to other terrain"
// @Failure 422 {object} iris.Problem "you must rob player who has construction next to robber"
// @Failure 422 {object} iris.Problem "selected player must have construction next to robber"
// @Router /catan/games/{id}/play-monopoly-card [post]
func (c CatanController) PlayMonopolyCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playMonopolyCardRequest := &requests.PlayMonopolyCard{}

	if err := ctx.ReadJSON(playMonopolyCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playMonopolyCardPbRequest := &pb_requests.PlayMonopolyCard{
		GameID:                    id,
		UserID:                    userID,
		DevelopmentCardID:         playMonopolyCardRequest.DevelopmentCardID,
		DemandingResourceCardType: playMonopolyCardRequest.DemandingResourceCardType,
	}

	if _, err := c.catanClient.PlayMonopolyCard(ctx, playMonopolyCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "MonopolyCardPlayed",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}

// @Summary Play victory point card
// @Description Play victory point development card from your stack at any phase of started state
// @Accept  json
// @Produce  json
// @Param	id	   path    string	true	"Game ID"
// @Param	developmentCardID	   body    requests.PlayVictoryPointCard	true	"Development Card ID"
// @Success 200 {nil} nil "ok"
// @Failure 400 {object} iris.Problem "the request contains invalid parameters"
// @Failure 404 {object} iris.Problem "development card not found"
// @Failure 422 {object} iris.Problem "game has not started yet"
// @Failure 422 {object} iris.Problem "game already finished"
// @Failure 422 {object} iris.Problem "you are not in turn"
// @Failure 422 {object} iris.Problem "selected development card must be victory point card"
// @Failure 422 {object} iris.Problem "selected development card is unavailable to use"
// @Failure 422 {object} iris.Problem "robber must be moved to other terrain"
// @Failure 422 {object} iris.Problem "you must rob player who has construction next to robber"
// @Failure 422 {object} iris.Problem "selected player must have construction next to robber"
// @Router /catan/games/{id}/play-victory-point-card [post]
func (c CatanController) PlayVictoryPointCard(ctx iris.Context, id string) (mvc.Result, error) {
	userID := context_values.GetUserID(ctx)
	playVictoryPointCardRequest := &requests.PlayVictoryPointCard{}

	if err := ctx.ReadJSON(playVictoryPointCardRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	playVictoryPointCardPbRequest := &pb_requests.PlayVictoryPointCard{
		GameID:            id,
		UserID:            userID,
		DevelopmentCardID: playVictoryPointCardRequest.DevelopmentCardID,
	}

	if _, err := c.catanClient.PlayVictoryPointCard(ctx, playVictoryPointCardPbRequest); err != nil {
		return nil, errors.WithStack(err)
	}

	messageResponse := &responses.Message{
		UserID: userID,
	}

	c.websocketServer.Broadcast(nil, neffos.Message{
		Namespace: "Catan",
		Room:      id,
		Event:     "VictoryPointCardPlayed",
		Body:      neffos.Marshal(messageResponse),
	})

	return &mvc.Response{
		Code: iris.StatusOK,
	}, nil
}
