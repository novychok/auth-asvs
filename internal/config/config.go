package config

import (
	"os"
	"time"

	"github.com/novychok/authasvs/internal/handler/authapiv1"
	"github.com/novychok/authasvs/internal/pkg/postgres"
	"github.com/spf13/viper"
)

type JWT struct {
	ExpiresIn        time.Duration `mapstructure:"JWT_EXPIRES_IN"`
	RefreshExpiresIn time.Duration `mapstructure:"JWT_REFRESH_EXPIRES_IN"`
}

type Config struct {
	Postgres  postgres.Config  `mapstructure:",squash"`
	AuthApiV1 authapiv1.Config `mapstructure:",squash"`
	JWT       JWT              `mapstructure:",squash"`
}

func New() (*Config, error) {
	viper.AutomaticEnv()

	configPath := "auth.env"

	_, err := os.Stat(configPath)
	if err == nil {
		viper.AddConfigPath("env")
		viper.SetConfigType("env")
		viper.SetConfigName("auth")
		viper.AddConfigPath(".")

		err := viper.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	cfg := &Config{}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetPostgres(cfg *Config) *postgres.Config {
	return &cfg.Postgres
}

func GetPlatfromAPIV1(cfg *Config) *authapiv1.Config {
	return &cfg.AuthApiV1
}

func GetJWT(cfg *Config) *JWT {
	return &cfg.JWT
}
