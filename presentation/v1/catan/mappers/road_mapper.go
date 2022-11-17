package mappers

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan/responses"
	pb_responses "github.com/vulpes-ferrilata/catan-service-proto/pb/responses"
)

type roadMapper struct{}

func (r roadMapper) ToHttpResponse(roadPbResponse *pb_responses.Road) (*responses.Road, error) {
	if roadPbResponse == nil {
		return nil, nil
	}

	pathHttpResponse, err := pathMapper{}.ToHttpResponse(roadPbResponse.GetPath())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &responses.Road{
		ID:   roadPbResponse.GetID(),
		Path: pathHttpResponse,
	}, nil
}
