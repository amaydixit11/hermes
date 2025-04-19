package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	// Server configuration
	// In your config struct
	Server struct {
		Address         string `mapstructure:"address"`
		Port            int    `mapstructure:"port"`
		TimeoutRead     int    `mapstructure:"timeout_read"`
		TimeoutWrite    int    `mapstructure:"timeout_write"`
		TimeoutIdle     int    `mapstructure:"timeout_idle"`
		TimeoutShutdown int    `mapstructure:"timeout_shutdown"`
	} `mapstructure:"server"`

	// Database configuration
	Database struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	// JWT configuration
	JWT struct {
		Secret     string `mapstructure:"secret"`
		ExpiryHour int    `mapstructure:"expiry_hour"`
	} `mapstructure:"jwt"`

	// HealthCheck configuration
	HealthCheck struct {
		Interval int `mapstructure:"interval"` // in seconds
	} `mapstructure:"health_check"`

	// Log level (debug, info, warn, error)
	LogLevel string `mapstructure:"log_level"`

	// Environment (development, staging, production)
	Environment string `mapstructure:"environment"`
}

// Load reads configuration from file or environment variables
func Load() (*Config, error) {
	// Set default config name and path
	configName := "config"
	configPath := "configs"

	// Check for environment-specific config from ENV var
	env := os.Getenv("HERMES_ENV")
	if env != "" {
		configName = fmt.Sprintf("config.%s", strings.ToLower(env))
	}

	// Initialize Viper
	v := viper.New()
	v.SetConfigName(configName)
	v.AddConfigPath(configPath)
	v.AddConfigPath(".")
	v.SetConfigType("yaml")

	// Enable environment variable override
	v.AutomaticEnv()
	v.SetEnvPrefix("HERMES")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the config file
	if err := v.ReadInConfig(); err != nil {
		// Config file not found is a common error that we can handle
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found: %s", filepath.Join(configPath, configName+".yaml"))
		}
		return nil, fmt.Errorf("error reading config file: %s", err)
	}

	// Parse config into struct
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %s", err)
	}

	// Set environment from file if not set by ENV var
	if env == "" {
		config.Environment = v.GetString("environment")
	} else {
		config.Environment = env
	}

	// Use default log level if not set
	if config.LogLevel == "" {
		config.LogLevel = "info"
	}

	return &config, nil
}

// DSN returns the PostgreSQL DSN string based on the config values
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.Username,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
