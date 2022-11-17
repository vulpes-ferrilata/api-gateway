package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

var GameDetailMapper gameDetailMapper = gameDetailMapper{}

type gameDetailMapper struct{}

func (g gameDetailMapper) ToHttpResponse(gameDetailPbResponse *pb_responses.GameDetail) (*responses.GameDetail, error) {
	if gameDetailPbResponse == nil {
		return nil, nil
	}

	activePlayerHttpResponse, err := playerMapper{}.ToHttpResponse(gameDetailPbResponse.GetActivePlayer())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	playerHttpResponses, err := slices.Map(func(playerPbResponse *pb_responses.Player) (*responses.Player, error) {
		return playerMapper{}.ToHttpResponse(playerPbResponse)
	}, gameDetailPbResponse.GetPlayers())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	diceHttpResponses, err := slices.Map(func(dicePbResponse *pb_responses.Dice) (*responses.Dice, error) {
		return diceMapper{}.ToHttpResponse(dicePbResponse)
	}, gameDetailPbResponse.GetDices())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	achievementHttpResponses, err := slices.Map(func(achievementPbResponse *pb_responses.Achievement) (*responses.Achievement, error) {
		return achievementMapper{}.ToHttpResponse(achievementPbResponse)
	}, gameDetailPbResponse.GetAchievements())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resourceCardHttpResponses, err := slices.Map(func(resourceCardPbResponse *pb_responses.ResourceCard) (*responses.ResourceCard, error) {
		return resourceCardMapper{}.ToHttpResponse(resourceCardPbResponse)
	}, gameDetailPbResponse.GetResourceCards())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	developmentCardHttpResponses, err := slices.Map(func(developmentCardPbResponse *pb_responses.DevelopmentCard) (*responses.DevelopmentCard, error) {
		return developmentCardMapper{}.ToHttpResponse(developmentCardPbResponse)
	}, gameDetailPbResponse.GetDevelopmentCards())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	terrainHttpResponses, err := slices.Map(func(terrainPbResponse *pb_responses.Terrain) (*responses.Terrain, error) {
		return terrainMapper{}.ToHttpResponse(terrainPbResponse)
	}, gameDetailPbResponse.GetTerrains())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	landHttpResponses, err := slices.Map(func(landPbResponse *pb_responses.Land) (*responses.Land, error) {
		return landMapper{}.ToHttpResponse(landPbResponse)
	}, gameDetailPbResponse.GetLands())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pathHttpResponses, err := slices.Map(func(pathPbResponse *pb_responses.Path) (*responses.Path, error) {
		return pathMapper{}.ToHttpResponse(pathPbResponse)
	}, gameDetailPbResponse.GetPaths())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &responses.GameDetail{
		ID:               gameDetailPbResponse.GetID(),
		Status:           gameDetailPbResponse.GetStatus(),
		Phase:            gameDetailPbResponse.GetPhase(),
		Turn:             int(gameDetailPbResponse.GetTurn()),
		ActivePlayer:     activePlayerHttpResponse,
		Players:          playerHttpResponses,
		Dices:            diceHttpResponses,
		Achievements:     achievementHttpResponses,
		ResourceCards:    resourceCardHttpResponses,
		DevelopmentCards: developmentCardHttpResponses,
		Terrains:         terrainHttpResponses,
		Lands:            landHttpResponses,
		Paths:            pathHttpResponses,
	}, nil
}
