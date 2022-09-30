package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

func ToGameHttpResponse(gamePbResponse *pb_responses.Game) *responses.Game {
	if gamePbResponse == nil {
		return nil
	}

	activePlayerHttpResponse := toPlayerHttpResponse(gamePbResponse.GetActivePlayer())

	playerHttpResponses, _ := slices.Map(func(playerPbResponse *pb_responses.Player) (*responses.Player, error) {
		return toPlayerHttpResponse(playerPbResponse), nil
	}, gamePbResponse.GetPlayers())

	diceHttpResponses, _ := slices.Map(func(dicePbResponse *pb_responses.Dice) (*responses.Dice, error) {
		return toDiceHttpResponse(dicePbResponse), nil
	}, gamePbResponse.GetDices())

	achievementHttpResponses, _ := slices.Map(func(achievementPbResponse *pb_responses.Achievement) (*responses.Achievement, error) {
		return toAchievementHttpResponse(achievementPbResponse), nil
	}, gamePbResponse.GetAchievements())

	resourceCardHttpResponses, _ := slices.Map(func(resourceCardPbResponse *pb_responses.ResourceCard) (*responses.ResourceCard, error) {
		return toResourceCardHttpResponse(resourceCardPbResponse), nil
	}, gamePbResponse.GetResourceCards())

	developmentCardHttpResponses, _ := slices.Map(func(developmentCardPbResponse *pb_responses.DevelopmentCard) (*responses.DevelopmentCard, error) {
		return toDevelopmentCardHttpResponse(developmentCardPbResponse), nil
	}, gamePbResponse.GetDevelopmentCards())

	terrainHttpResponses, _ := slices.Map(func(terrainPbResponse *pb_responses.Terrain) (*responses.Terrain, error) {
		return toTerrainHttpResponse(terrainPbResponse), nil
	}, gamePbResponse.GetTerrains())

	landHttpResponses, _ := slices.Map(func(landPbResponse *pb_responses.Land) (*responses.Land, error) {
		return toLandHttpResponse(landPbResponse), nil
	}, gamePbResponse.GetLands())

	pathHttpResponses, _ := slices.Map(func(pathPbResponse *pb_responses.Path) (*responses.Path, error) {
		return toPathHttpResponse(pathPbResponse), nil
	}, gamePbResponse.GetPaths())

	return &responses.Game{
		ID:               gamePbResponse.GetID(),
		Status:           gamePbResponse.GetStatus(),
		Phase:            gamePbResponse.GetPhase(),
		Turn:             int(gamePbResponse.GetTurn()),
		ActivePlayer:     activePlayerHttpResponse,
		Players:          playerHttpResponses,
		Dices:            diceHttpResponses,
		Achievements:     achievementHttpResponses,
		ResourceCards:    resourceCardHttpResponses,
		DevelopmentCards: developmentCardHttpResponses,
		Terrains:         terrainHttpResponses,
		Lands:            landHttpResponses,
		Paths:            pathHttpResponses,
	}
}
