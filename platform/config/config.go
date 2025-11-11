package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Health   HealthConfig   `mapstructure:"health"`
	Swagger  SwaggerConfig  `mapstructure:"swagger"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port             string        `mapstructure:"port"`
	ReadTimeout      time.Duration `mapstructure:"read_timeout"`
	WriteTimeout     time.Duration `mapstructure:"write_timeout"`
	IdleTimeout      time.Duration `mapstructure:"idle_timeout"`
	AllowedOrigins   []string      `mapstructure:"allowed_origins"`
	AllowedMethods   []string      `mapstructure:"allowed_methods"`
	AllowedHeaders   []string      `mapstructure:"allowed_headers"`
	AllowCredentials bool          `mapstructure:"allow_credentials"`
	MaxAge           int           `mapstructure:"max_age"`
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	SSLMode         string        `mapstructure:"sslmode"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

// HealthConfig holds health check configuration
type HealthConfig struct {
	DatabaseTimeout time.Duration `mapstructure:"database_timeout"`
}

// SwaggerConfig holds Swagger/API documentation configuration
type SwaggerConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	OpenApiFilePath string `mapstructure:"openapi_file_path"`
	SwaggerFilePath string `mapstructure:"swagger_file_path"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	setDefaults()

	viper.SetEnvPrefix("WALLET")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./resources")

	// Reading config file is optional
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// setDefaults sets default values for all configuration sections
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.read_timeout", "10s")
	viper.SetDefault("server.write_timeout", "10s")
	viper.SetDefault("server.idle_timeout", "30s")
	viper.SetDefault("server.allowed_origins", []string{"*"})
	viper.SetDefault("server.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("server.allowed_headers", []string{"Content-Type", "Authorization", "X-Requested-With"})
	viper.SetDefault("server.allow_credentials", true)
	viper.SetDefault("server.max_age", 86400)

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.dbname", "go_clean_db")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.conn_max_lifetime", "5m")

	// Logging defaults
	viper.SetDefault("logging.level", "info")

	// Health check defaults
	viper.SetDefault("health.database_timeout", "5s")

	// Swagger defaults
	viper.SetDefault("swagger.enabled", true)
	viper.SetDefault("swagger.swagger_file_path", "./resources/swagger.html")
	viper.SetDefault("swagger.openapi_file_path", "./resources/openapi.yaml")
}
