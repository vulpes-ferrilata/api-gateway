package main

import (
	"github.com/kataras/iris/v12"
	"github.com/pkg/errors"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure"
	"github.com/vulpes-ferrilata/api-gateway/infrastructure/config"
)

func main() {
	container := infrastructure.NewContainer()

	if err := container.Invoke(func(server *iris.Application, config config.Config) error {
		if err := server.Run(
			iris.Addr(config.Server.Address),
			iris.WithoutServerError(iris.ErrServerClosed),
		); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}
