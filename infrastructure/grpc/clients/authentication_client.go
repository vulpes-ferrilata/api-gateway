package clients

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/config"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/interceptors"
	"github.com/vulpes-ferrilata/authentication-service-proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewAuthenticationClient(config config.Config,
	localeInterceptor *interceptors.LocaleInterceptor) (pb.AuthenticationClient, error) {
	conn, err := grpc.Dial(config.AuthenticationService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(localeInterceptor.ClientUnaryInterceptor))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	authenticationClient := pb.NewAuthenticationClient(conn)

	return authenticationClient, nil
}
