package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GRPCPort    string
	DatabaseURL string
	RedisAddr   string
	SMTPHost    string
	SMTPPort    int
	SMTPUser    string
	SMTPPass    string
	LogLevel    string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{
		GRPCPort:    viper.GetString("grpc.port"),
		DatabaseURL: viper.GetString("database.url"),
		RedisAddr:   viper.GetString("redis.addr"),
		SMTPHost:    viper.GetString("smtp.host"),
		SMTPPort:    viper.GetInt("smtp.port"),
		SMTPUser:    viper.GetString("smtp.user"),
		SMTPPass:    viper.GetString("smtp.pass"),
		LogLevel:    viper.GetString("log.level"),
	}, nil
}
