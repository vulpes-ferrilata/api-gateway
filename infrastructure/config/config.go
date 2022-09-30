package config

type Config struct {
	Server                ServerConfig  `mapstructure:"server"`
	UserService           ServiceConfig `mapstructure:"user_service"`
	AuthenticationService ServiceConfig `mapstructure:"authentication_service"`
	CatanService          ServiceConfig `mapstructure:"catan_service"`
	ChatService           ServiceConfig `mapstructure:"chat_service"`
}
