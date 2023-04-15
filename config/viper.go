package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	DATABASE_URI  string `mapstructure:"DATABASE_URI"`
	DATABASE_HOST string `mapstructure:"DATABASE_HOST"`
	PORT         string `mapstructure:"PORT"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			DATABASE_URI:  os.Getenv("DATABASE_URI"),
			DATABASE_HOST: os.Getenv("DATABASE_HOST"),
			PORT:         os.Getenv("PORT"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here
	if config.DATABASE_URI == "" {
		err = errors.New("DATABASE_URI is required")
		return
	}

	if config.DATABASE_HOST == "" {
		err = errors.New("DATABASE_HOST is required")
		return
	}

	return
}
