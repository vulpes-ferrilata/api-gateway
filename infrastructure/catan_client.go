package infrastructure

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/config"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/interceptors"
	"github.com/vulpes-ferrilata/shared/proto/v1/catan"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCatanClient(config config.Config,
	localeInterceptor *interceptors.LocaleInterceptor) (catan.CatanClient, error) {
	conn, err := grpc.Dial(config.CatanService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(localeInterceptor.ClientUnaryInterceptor))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	catanClient := catan.NewCatanClient(conn)

	return catanClient, nil
}
