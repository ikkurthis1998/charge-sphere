package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Temporal TemporalConfig `mapstructure:"temporal"`
	OCPI     OCPIConfig     `mapstructure:"ocpi"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type MongoDBConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Timeout  int    `mapstructure:"timeout"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type TemporalConfig struct {
	Host      string `mapstructure:"host"`
	Namespace string `mapstructure:"namespace"`
}

type OCPIConfig struct {
	Version     string `mapstructure:"version"`
	CountryCode string `mapstructure:"country_code"`
	PartyID     string `mapstructure:"party_id"`
}

// Load reads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Set up environment variable binding
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	// Map environment variables to config keys
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "SERVER_MODE")

	viper.BindEnv("mongodb.uri", "MONGODB_URI")
	viper.BindEnv("mongodb.database", "MONGODB_DATABASE")
	viper.BindEnv("mongodb.timeout", "MONGODB_TIMEOUT")

	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	viper.BindEnv("temporal.host", "TEMPORAL_HOST")
	viper.BindEnv("temporal.namespace", "TEMPORAL_NAMESPACE")

	viper.BindEnv("ocpi.version", "OCPI_VERSION")
	viper.BindEnv("ocpi.country_code", "OCPI_COUNTRY_CODE")
	viper.BindEnv("ocpi.party_id", "OCPI_PARTY_ID")

	// Try to read config file (optional in production)
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		// Config file is optional if all env vars are set
		fmt.Printf("Warning: config file not found, using environment variables\n")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetServerAddr returns the server address in host:port format
func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
