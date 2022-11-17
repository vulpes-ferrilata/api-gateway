package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type constructionMapper struct{}

func (c constructionMapper) ToHttpResponse(constructionPbResponse *pb_responses.Construction) (*responses.Construction, error) {
	if constructionPbResponse == nil {
		return nil, nil
	}

	landHttpResponse, err := landMapper{}.ToHttpResponse(constructionPbResponse.GetLand())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &responses.Construction{
		ID:   constructionPbResponse.GetID(),
		Type: constructionPbResponse.GetType(),
		Land: landHttpResponse,
	}, nil
}
