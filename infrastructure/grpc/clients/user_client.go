package clients

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/config"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/interceptors"
	"github.com/vulpes-ferrilata/user-service-proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewUserClient(config config.Config,
	localeInterceptor *interceptors.LocaleInterceptor) (pb.UserClient, error) {
	conn, err := grpc.Dial(config.UserService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(localeInterceptor.ClientUnaryInterceptor))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	userClient := pb.NewUserClient(conn)

	return userClient, nil
}
