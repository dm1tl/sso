package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	Env         string        `yaml:"env" env-required:"true"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    string        `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}

// loading a config from CONFIG_PATH env variable
func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config_path is empty")
	}
	if _, err := os.Stat(path); err != nil {
		panic("no config file on this path " + path)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("couldn't read config")
	}
	return &cfg
}

// loading config path from env variables
func fetchConfigPath() string {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("couldn't load env variables %v", err)
	}
	res := os.Getenv("CONFIG_PATH")
	return res
}
