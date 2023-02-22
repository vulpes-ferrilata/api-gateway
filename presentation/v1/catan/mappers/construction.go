package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type constructionMapper struct{}

func (c constructionMapper) ToHttpResponse(constructionPbResponse *pb_models.Construction) (*models.Construction, error) {
	if constructionPbResponse == nil {
		return nil, nil
	}

	landHttpResponse, err := landMapper{}.ToHttpResponse(constructionPbResponse.GetLand())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.Construction{
		ID:   constructionPbResponse.GetID(),
		Type: constructionPbResponse.GetType(),
		Land: landHttpResponse,
	}, nil
}
