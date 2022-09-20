package presentation

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/iris/middlewares"
	"github.com/vulpes-ferrilata/api-gateway/presentation/v1/controllers"
)

type Router interface {
	Init(server *iris.Application)
}

func NewRouter(errorHandlerMiddleware *middlewares.ErrorHandlerMiddleware,
	authenticationMiddleware *middlewares.AuthenticationMiddleware,
	authenticationController *controllers.AuthenticationController,
	userController *controllers.UserController,
	catanController *controllers.CatanController,
	websocketServer *neffos.Server) Router {
	return &router{
		errorHandlerMiddleware:   errorHandlerMiddleware,
		authenticationMiddleware: authenticationMiddleware,
		authenticationController: authenticationController,
		userController:           userController,
		catanController:          catanController,
		websocketServer:          websocketServer,
	}
}

type router struct {
	errorHandlerMiddleware   *middlewares.ErrorHandlerMiddleware
	authenticationMiddleware *middlewares.AuthenticationMiddleware
	authenticationController *controllers.AuthenticationController
	userController           *controllers.UserController
	catanController          *controllers.CatanController
	websocketServer          *neffos.Server
}

func (r router) Init(server *iris.Application) {
	server.Get("/", websocket.Handler(r.websocketServer))

	api := server.Party("/api")
	v1 := api.Party("/v1")

	auth := mvc.New(v1.Party("/auth"))
	auth.HandleError(r.errorHandlerMiddleware.Serve())
	auth.Handle(r.authenticationController)

	user := mvc.New(v1.Party("/user"))
	user.Router.Use(r.authenticationMiddleware.Serve())
	user.HandleError(r.errorHandlerMiddleware.Serve())
	user.Handle(r.userController)

	catan := mvc.New(v1.Party("/catan"))
	catan.Router.Use(r.authenticationMiddleware.Serve())
	catan.HandleError(r.errorHandlerMiddleware.Serve())
	catan.Handle(r.catanController)
}
