package infrastructure

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/vulpes-ferrilata/api-gateway/config"
)

func NewConfig() (config.Config, error) {
	config := config.Config{}

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return config, errors.WithStack(err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, errors.WithStack(err)
	}

	return config, nil
}
