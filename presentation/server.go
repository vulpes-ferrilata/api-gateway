package presentation

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/iris/middlewares"
)

func NewServer(validator *validator.Validate, localeMiddleware *middlewares.LocaleMiddleware,
	router Router) *iris.Application {
	server := iris.New()

	server.Validator = validator

	server.UseRouter(logger.New())
	server.UseRouter(recover.New())
	server.UseRouter(cors.New().
		ExtractOriginFunc(cors.DefaultOriginExtractor).
		ReferrerPolicy(cors.NoReferrerWhenDowngrade).
		AllowOriginFunc(cors.AllowAnyOrigin).
		Handler())

	server.Use(localeMiddleware.Serve())

	router.Init(server)

	return server
}
