package infrastructure

import (
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/clients"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/interceptors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/iris"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/iris/middlewares"
	"github.com/vulpes-ferrilata/api-gateway/presentation"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/user"
	"go.uber.org/dig"
)

func NewContainer() *dig.Container {
	container := dig.New()

	//Infrastructure layer
	container.Provide(NewConfig)
	container.Provide(NewUniversalTranslator)
	container.Provide(NewValidator)
	container.Provide(iris.NewServer)
	container.Provide(iris.NewWebsocketServer)
	//--GRPC Clients
	container.Provide(clients.NewUserClient)
	container.Provide(clients.NewAuthenticationClient)
	container.Provide(clients.NewCatanClient)
	container.Provide(clients.NewChatClient)
	//--Grpc interceptors
	container.Provide(interceptors.NewLocaleInterceptor)
	//--Iris middlewares
	container.Provide(middlewares.NewErrorHandlerMiddleware)
	container.Provide(middlewares.NewAuthenticationMiddleware)
	container.Provide(middlewares.NewLocaleMiddleware)

	//Presentation layer
	//--Router
	container.Provide(presentation.NewRouter)
	//--Controllers
	container.Provide(user.NewUserController)
	container.Provide(authentication.NewAuthenticationController)
	container.Provide(catan.NewCatanController)
	container.Provide(chat.NewChatController)

	return container
}
