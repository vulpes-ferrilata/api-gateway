package clients

import (
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/config"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/interceptors"
	"github.com/vulpes-ferrilata/chat-service-proto/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewChatClient(config config.Config,
	localeInterceptor *interceptors.LocaleInterceptor) (pb.ChatClient, error) {
	conn, err := grpc.Dial(config.ChatService.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(localeInterceptor.ClientUnaryInterceptor))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	catanClient := pb.NewChatClient(conn)

	return catanClient, nil
}
