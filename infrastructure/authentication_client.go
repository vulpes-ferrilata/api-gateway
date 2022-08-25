package infrastructure

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/config"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/interceptors"
	"github.com/vulpes-ferrilata/shared/proto/v1/authentication"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthenticationClient(config config.Config,
	localeInterceptor *interceptors.LocaleInterceptor) (authentication.AuthenticationClient, error) {
	conn, err := grpc.Dial(config.AuthenticationService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(localeInterceptor.ClientUnaryInterceptor))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	authenticationClient := authentication.NewAuthenticationClient(conn)

	return authenticationClient, nil
}
