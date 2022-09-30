package presentation

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/iris/middlewares"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/authentication"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/catan"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/chat"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/user"
)

type Router interface {
	Init(server *iris.Application)
}

func NewRouter(errorHandlerMiddleware *middlewares.ErrorHandlerMiddleware,
	authenticationMiddleware *middlewares.AuthenticationMiddleware,
	userController *user.UserController,
	authenticationController *authentication.AuthenticationController,
	catanController *catan.CatanController,
	websocketServer *neffos.Server,
	chatController *chat.ChatController) Router {
	return &router{
		errorHandlerMiddleware:   errorHandlerMiddleware,
		authenticationMiddleware: authenticationMiddleware,
		userController:           userController,
		authenticationController: authenticationController,
		catanController:          catanController,
		websocketServer:          websocketServer,
		chatController:           chatController,
	}
}

type router struct {
	errorHandlerMiddleware   *middlewares.ErrorHandlerMiddleware
	authenticationMiddleware *middlewares.AuthenticationMiddleware
	userController           *user.UserController
	authenticationController *authentication.AuthenticationController
	catanController          *catan.CatanController
	websocketServer          *neffos.Server
	chatController           *chat.ChatController
}

func (r router) Init(server *iris.Application) {
	server.Get("/", websocket.Handler(r.websocketServer))

	api := server.Party("/api")
	v1 := api.Party("/v1")

	user := mvc.New(v1.Party("/user"))
	user.Router.Use(r.authenticationMiddleware.Serve())
	user.HandleError(r.errorHandlerMiddleware.Serve())
	user.Handle(r.userController)

	auth := mvc.New(v1.Party("/auth"))
	auth.HandleError(r.errorHandlerMiddleware.Serve())
	auth.Handle(r.authenticationController)

	catan := mvc.New(v1.Party("/catan"))
	catan.Router.Use(r.authenticationMiddleware.Serve())
	catan.HandleError(r.errorHandlerMiddleware.Serve())
	catan.Handle(r.catanController)

	chat := mvc.New(v1.Party("/chat"))
	chat.Router.Use(r.authenticationMiddleware.Serve())
	chat.HandleError(r.errorHandlerMiddleware.Serve())
	chat.Handle(r.chatController)
}
