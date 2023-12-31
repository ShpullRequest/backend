package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var Config NodeConfig

type NodeConfig struct {
	Address string `env:"ADDRESS"`

	AppAccessToken string `env:"APP_ACCESS_TOKEN"`
	AppSecretToken string `env:"APP_SECRET_TOKEN"`
	VkMapsAPIKey   string `env:"VK_MAPS_API_KEY"`

	MasterDSN      string `env:"MASTER_DSN"`
	MasterMaxOpen  int    `env:"MASTER_MAX_OPEN"`
	ReplicaDSN     string `env:"REPLICA_DSN"`
	ReplicaMaxOpen int    `env:"REPLICA_MAX_OPEN"`
	MigrationsFlag bool   `env:"MIGRATIONS_FLAG"`

	ProdFlag bool `env:"PROD_FLAG"`
}

func Load() {
	flag.StringVar(&Config.Address, "address", ":8080", "api address")

	flag.StringVar(&Config.AppAccessToken, "app-access-token", "", "app access token")
	flag.StringVar(&Config.AppSecretToken, "app-secret-token", "", "app secret token")
	flag.StringVar(&Config.VkMapsAPIKey, "vk-maps-api-key", "", "vk maps api key")

	flag.StringVar(&Config.MasterDSN, "master-dsn", "", "postgres master dsn")
	flag.IntVar(&Config.MasterMaxOpen, "master-max-open", 6, "maximum opened pools for master")
	flag.StringVar(&Config.ReplicaDSN, "replica-dsn", "", "postgres replica dsn")
	flag.IntVar(&Config.ReplicaMaxOpen, "replica-max-open", 6, "maximum opened pools for replica")
	flag.BoolVar(&Config.MigrationsFlag, "migrations-flag", false, "database flag migrations")

	flag.BoolVar(&Config.ProdFlag, "prod-flag", false, "flag for production server")
}

func Parse() error {
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		return err
	}

	return env.Parse(&Config)
}
