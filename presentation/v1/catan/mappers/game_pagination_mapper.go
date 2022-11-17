package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/utils/slices"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

var GamePaginationMapper gamePaginationMapper = gamePaginationMapper{}

type gamePaginationMapper struct{}

func (g gamePaginationMapper) ToHttpResponse(gamePaginationPbResponse *pb_responses.GamePagination) (*responses.GamePagination, error) {
	if gamePaginationPbResponse == nil {
		return nil, nil
	}

	gameHttpResponses, err := slices.Map(func(gamePbResponse *pb_responses.Game) (*responses.Game, error) {
		return gameMapper{}.ToHttpResponse(gamePbResponse)
	}, gamePaginationPbResponse.GetData())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &responses.GamePagination{
		Total: int(gamePaginationPbResponse.GetTotal()),
		Data:  gameHttpResponses,
	}, nil
}
