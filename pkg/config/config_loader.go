package config

import (
	"os"

	"github.com/spf13/viper"
)

type Configs interface {
	AppConfig |
		MigrationConfig |
		WorkerConfig
}

func LoadConfigFromEnv[K Configs](c *K) error {
	v := viper.New()

	if os.Getenv("ENV") != "production" {
		v.AddConfigPath(".")
		v.SetConfigName(".env")
		v.SetConfigType("env")
		v.AutomaticEnv()
	}

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.Unmarshal(&c); err != nil {
		return err
	}

	return nil
}
