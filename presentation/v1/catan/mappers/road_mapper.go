package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/models"
	pb_models "github.com/vulpes-ferrilata/catan-service-proto/pb/models"
)

type roadMapper struct{}

func (r roadMapper) ToHttpResponse(roadPbResponse *pb_models.Road) (*models.Road, error) {
	if roadPbResponse == nil {
		return nil, nil
	}

	pathHttpResponse, err := pathMapper{}.ToHttpResponse(roadPbResponse.GetPath())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.Road{
		ID:   roadPbResponse.GetID(),
		Path: pathHttpResponse,
	}, nil
}
