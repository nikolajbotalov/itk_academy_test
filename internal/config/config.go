package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"sync"
)

type Config struct {
	Listen        Listen
	DB            DB
	RetryAttempts uint `env:"RETRY_ATTEMPTS" envDefault:"10"`
}

type Listen struct {
	BindIP string `env:"BIND_IP" envDefault:"0.0.0.0"`
	Port   string `env:"PORT" envDefault:"8080"`
}

type DB struct {
	Host     string `env:"PG_HOST" envDefault:"localhost"`
	Port     string `env:"PG_PORT" envDefault:"5432"`
	User     string `env:"PG_USER" envDefault:"postgres"`
	Password string `env:"PG_PASSWORD" envDefault:"admin"`
	Name     string `env:"PG_NAME" envDefault:"wallet"`
}

var instance *Config
var once sync.Once

func LoadConfig(logger *zap.Logger) *Config {
	once.Do(func() {
		logger.Info("loading config")
		instance = &Config{}

		if err := cleanenv.ReadConfig("config.env", instance); err != nil {
			logger.Error("failed to read config", zap.Error(err))
		}
	})

	return instance
}
