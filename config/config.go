package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP    `yaml:"http"`
	MongoDB `yaml:"mongodb"`
}

type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

type MongoDB struct {
	URI string `env-required:"true" yaml:"uri" env:"MONGODB_URI"`
}

// New возвращает конфиг приложения
func New() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
