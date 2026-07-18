package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER" envDefault:"admin"`
	Password string `env:"POSTGRES_PASS" env:",required"`
}

type RabbitConfig struct {
	Port     int    `env:"RABBIT_HOST" envDefault:"5672"`
	User     string `env:"RABBIT_USER" envDefault:"guest"`
	Password string `env:"RABBIT_PASS" env:"required"`
}

type RedisConfig struct {
	Port int `env:"REDIS_PORT" envDefault:"6379"`
}

type ClickHouseConfig struct {
	Port int `env:"CLICKHOUSE_PORT" envDefault:"8123"`
}

func Load[T any]() (*T, error) {

	_ = godotenv.Load()

	var cfg T
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil

}
