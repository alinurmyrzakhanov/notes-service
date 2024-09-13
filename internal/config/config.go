package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config содержит конфигурацию приложения
type Config struct {
	ServerAddress         string `envconfig:"SERVER_ADDRESS" default:":8080"`
	DatabaseURL           string `envconfig:"DATABASE_URL" required:"true"`
	YandexSpellcheckerURL string `envconfig:"YANDEX_SPELLCHECKER_URL" default:"https://speller.yandex.net/services/spellservice.json/checkText"`
	JWTSecret             string `envconfig:"JWT_SECRET" required:"true"`
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
