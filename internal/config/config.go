package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
	JWT         JWT `yaml:"jwt"`
	TLS         TLS `yaml:"tls"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8443"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type JWT struct {
	AccessSecret    string        `yaml:"access_secret" env-required:"true"`
	RefreshSecret   string        `yaml:"refresh_secret" env-required:"true"`
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-default:"15m"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-default:"7d"`
}

type TLS struct {
	PathToCert string `yaml:"path_to_cert" env-required:"true"`
	PathToKey  string `yaml:"path_to_key" env-required:"true"`
	Port       string `yaml:"port" env-default:":8443"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("config path is not set")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", configPath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
