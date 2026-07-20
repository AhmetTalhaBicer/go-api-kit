package config

import (
	"fmt"
	"os"
	"time"
)

// Config holds all application-wide configuration.
// Values are read from environment variables.
// Sensible defaults are provided — override via a .env file in production.
type Config struct {
	App    AppConfig
	Server ServerConfig
	Log    LogConfig
	DB     DBConfig
}

// AppConfig holds general application metadata.
type AppConfig struct {
	Name    string
	Env     string // development | staging | production
	Version string
}

// ServerConfig holds HTTP server settings.
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// LogConfig holds logging settings.
type LogConfig struct {
	Level  string // debug | info | warn | error
	Format string // json | text
}

// DBConfig holds PostgreSQL connection settings.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN returns a PostgreSQL connection string built from the individual fields.
func (d DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// Load reads Config from environment variables.
// Falls back to the provided default value when a variable is not set.
func Load() *Config {
	return &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "go-api-kit"),
			Env:     getEnv("APP_ENV", "development"),
			Version: getEnv("APP_VERSION", "1.0.0"),
		},
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "mydb"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// getEnv returns the value of the environment variable named by key.
// If the variable is not set or empty, fallback is returned.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getDurationEnv reads an environment variable as a time.Duration.
// Falls back to the provided default when the variable is absent or unparseable.
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}
	return d
}
