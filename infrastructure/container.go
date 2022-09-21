package infrastructure

import (
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/grpc/interceptors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/iris/middlewares"
	"github.com/vulpes-ferrilata/api-gateway/presentation"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/controllers"
	"go.uber.org/dig"
)

func NewContainer() *dig.Container {
	container := dig.New()

	//Infrastructure layer
	container.Provide(NewConfig)
	container.Provide(NewUniversalTranslator)
	container.Provide(NewValidator)
	//--GRPC Clients
	container.Provide(NewUserClient)
	container.Provide(NewAuthenticationClient)
	container.Provide(NewCatanClient)
	//--Grpc interceptors
	container.Provide(interceptors.NewLocaleInterceptor)
	//--Iris middlewares
	container.Provide(middlewares.NewErrorHandlerMiddleware)
	container.Provide(middlewares.NewAuthenticationMiddleware)
	container.Provide(middlewares.NewLocaleMiddleware)

	//Presentation layer
	//--Server
	container.Provide(presentation.NewServer)
	//--Websocket server
	container.Provide(presentation.NewWebsocketServer)
	//--Router
	container.Provide(presentation.NewRouter)
	//--Controllers
	container.Provide(controllers.NewAuthenticationController)
	container.Provide(controllers.NewUserController)
	container.Provide(controllers.NewCatanController)

	return container
}
