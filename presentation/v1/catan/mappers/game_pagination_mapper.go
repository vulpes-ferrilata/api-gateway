package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
	"github.com/vulpes-ferrilata/slices"
)

var GamePaginationMapper gamePaginationMapper = gamePaginationMapper{}

type gamePaginationMapper struct{}

func (g gamePaginationMapper) ToHttpResponse(gamePaginationPbResponse *pb_models.GamePagination) (*models.GamePagination, error) {
	if gamePaginationPbResponse == nil {
		return nil, nil
	}

	gameHttpResponses, err := slices.Map(func(gamePbResponse *pb_models.Game) (*models.Game, error) {
		return gameMapper{}.ToHttpResponse(gamePbResponse)
	}, gamePaginationPbResponse.GetData()...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.GamePagination{
		Total: int(gamePaginationPbResponse.GetTotal()),
		Data:  gameHttpResponses,
	}, nil
}
