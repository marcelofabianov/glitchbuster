package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	General   GeneralConfig   `mapstructure:"general"`
	Server    ServerConfig    `mapstructure:"server"`
	KurrentDB KurrentDBConfig `mapstructure:"kurrentdb"`
}

type GeneralConfig struct {
	Env string `mapstructure:"env"`
}

type ServerConfig struct {
	API APIConfig `mapstructure:"api"`
}

type APIConfig struct {
	Host         string        `mapstructure:"host"`
	Port         int           `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type KurrentDBConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
}

func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// --- Defaults ---
	v.SetDefault("general.env", "development")
	v.SetDefault("server.api.host", "0.0.0.0")
	v.SetDefault("server.api.port", 8080)
	v.SetDefault("server.api.read_timeout", "5s")
	v.SetDefault("server.api.write_timeout", "10s")
	v.SetDefault("kurrentdb.connection_string", "esdb://localhost:2113?tls=false")

	// --- Config File ---
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(path)

	// --- Environment Variables ---
	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// --- Reading Config ---
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	// --- Unmarshaling ---
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
