package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type playerMapper struct{}

func (p playerMapper) ToHttpResponse(playerPbResponse *pb_responses.Player) (*responses.Player, error) {
	if playerPbResponse == nil {
		return nil, nil
	}

	achievementHttpResponses, err := slices.Map(func(achievementPbResponse *pb_responses.Achievement) (*responses.Achievement, error) {
		return achievementMapper{}.ToHttpResponse(achievementPbResponse)
	}, playerPbResponse.GetAchievements())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resourceCardHttpResponses, err := slices.Map(func(resourceCardPbResponse *pb_responses.ResourceCard) (*responses.ResourceCard, error) {
		return resourceCardMapper{}.ToHttpResponse(resourceCardPbResponse)
	}, playerPbResponse.GetResourceCards())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	developmentCardHttpResponses, err := slices.Map(func(developmentCardPbResponse *pb_responses.DevelopmentCard) (*responses.DevelopmentCard, error) {
		return developmentCardMapper{}.ToHttpResponse(developmentCardPbResponse)
	}, playerPbResponse.GetDevelopmentCards())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	constructionHttpResponses, err := slices.Map(func(constructionPbResponse *pb_responses.Construction) (*responses.Construction, error) {
		return constructionMapper{}.ToHttpResponse(constructionPbResponse)
	}, playerPbResponse.GetConstructions())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	roadHttpResponses, err := slices.Map(func(roadPbResponse *pb_responses.Road) (*responses.Road, error) {
		return roadMapper{}.ToHttpResponse(roadPbResponse)
	}, playerPbResponse.GetRoads())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &responses.Player{
		ID:                 playerPbResponse.GetID(),
		UserID:             playerPbResponse.GetUserID(),
		Color:              playerPbResponse.GetColor(),
		TurnOrder:          int(playerPbResponse.GetTurnOrder()),
		ReceivedOffer:      playerPbResponse.GetReceivedOffer(),
		DiscardedResources: playerPbResponse.GetDiscardedResources(),
		Score:              int(playerPbResponse.GetScore()),
		Achievements:       achievementHttpResponses,
		ResourceCards:      resourceCardHttpResponses,
		DevelopmentCards:   developmentCardHttpResponses,
		Constructions:      constructionHttpResponses,
		Roads:              roadHttpResponses,
	}, nil
}
