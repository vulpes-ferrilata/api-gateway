package main

import (
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/config"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure"
)

// @title Tumbleweeds Api Gateway
// @version 1.0
// @description Backend Rest API for Tumbleweeds project
// @termsOfService http://swagger.io/terms/

// @contact.name Trung Hieu Nguyen
// @contact.email hieunguyen6694@gmail.com

// @BasePath /api/v1
func main() {
	container := infrastructure.NewContainer()

	if err := container.Invoke(func(server *iris.Application, config config.Config) error {
		if err := server.Run(
			iris.Addr(config.Server.Address),
			iris.WithoutServerError(iris.ErrServerClosed),
			iris.WithoutPathCorrectionRedirection,
		); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}
