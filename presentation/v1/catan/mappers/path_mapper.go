package mappers

import (
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type pathMapper struct{}

func (p pathMapper) ToHttpResponse(pathPbResponse *pb_models.Path) (*models.Path, error) {
	if pathPbResponse == nil {
		return nil, nil
	}

	return &models.Path{
		ID:       pathPbResponse.GetID(),
		Q:        int(pathPbResponse.GetQ()),
		R:        int(pathPbResponse.GetR()),
		Location: pathPbResponse.GetLocation(),
	}, nil
}
