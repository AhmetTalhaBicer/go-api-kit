package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents the complete application configuration structure.
type Config struct {
	App          AppConfig
	Server       ServerConfig
	Log          LogConfig
	DB           DBConfig
	Redis        RedisConfig
	Auth         AuthConfig
	Integrations IntegrationsConfig
	Business     BusinessConfig
	Worker       WorkerConfig
}

// AppConfig holds core application metadata and environment settings.
type AppConfig struct {
	Name    string
	Env     string // dev  | prod
	Version string
	Port    string
	Debug   bool
	BaseURL string
}

// ServerConfig holds HTTP server timeouts and networking configurations.
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// LogConfig holds logger level and output format.
type LogConfig struct {
	Level  string // debug | info | warn | error
	Format string // json | text
}

// DBConfig holds PostgreSQL connection pool and migration DSN settings.
type DBConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	Schema          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	DatabaseURL     string
}

// DSN builds and returns the PostgreSQL connection string.
func (d DBConfig) DSN() string {
	if d.DatabaseURL != "" {
		return os.ExpandEnv(d.DatabaseURL)
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.Schema,
	)
}

// RedisConfig holds Redis cache, rate-limiter, and session configurations.
type RedisConfig struct {
	Host       string
	Port       string
	Password   string
	DB         int
	TTLDefault time.Duration
}

// Addr returns the host:port string formatted for the Redis client.
func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

// AuthConfig holds JWT and OTP security settings.
type AuthConfig struct {
	JWTAccessSecret            string
	JWTRefreshSecret           string
	JWTAccessExpirationMinutes int
	JWTRefreshExpirationDays   int
	OTPExpirationMinutes       int
	OTPRateLimitPerMinute      int
	OTPRateLimitPerHour        int
	OTPMasterTestCode          string
}

// IntegrationsConfig holds driver switches (mock/real) and third-party API credentials.
type IntegrationsConfig struct {
	SMSProvider     string // mock 
	StorageProvider string // mock | s3 | r2
	PaymentProvider string // mock 

	// AWS S3 / Cloudflare R2
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	AWSS3BucketName    string
	AWSS3Endpoint      string

}

// BusinessConfig holds domain constants, commission rates, and threshold limits.
type BusinessConfig struct {
	PlatformDefaultCommissionRate float64
	AutoReleaseDeliveryHours      int
	DefaultFireThresholdPercent   float64
}

// WorkerConfig holds Outbox poller and background task worker settings.
type WorkerConfig struct {
	OutboxPollerIntervalMs int
	OutboxBatchSize        int
}

// Validate checks critical security constraints before application boot.
func (c *Config) Validate() error {
	if c.App.Env == "production" {
		if c.Auth.JWTAccessSecret == "" || strings.Contains(c.Auth.JWTAccessSecret, "change_in_prod") {
			return errors.New("CRITICAL SECURITY RISK: JWT_ACCESS_SECRET must be set to a strong secret in production")
		}
		if c.Auth.JWTRefreshSecret == "" || strings.Contains(c.Auth.JWTRefreshSecret, "change_in_prod") {
			return errors.New("CRITICAL SECURITY RISK: JWT_REFRESH_SECRET must be set to a strong secret in production")
		}
	}
	return nil
}

// Load reads configuration from environment variables, automatically loading .env.dev or .env.prod.
func Load() *Config {
	env := os.Getenv("APP_ENV")
	if env == "production" || env == "prod" {
		loadDotEnv(".env.prod")
	} else {
		loadDotEnv(".env.dev")
	}

	port := getEnv("APP_PORT", "8080")

	cfg := &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "go-api-kit"),
			Env:     getEnv("APP_ENV", "development"),
			Version: getEnv("APP_VERSION", "2.0.0"),
			Port:    port,
			Debug:   getBoolEnv("APP_DEBUG", true),
			BaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
		},
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         port,
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "app_user"),
			Password:        getEnv("DB_PASSWORD", "app_pass"),
			Name:            getEnv("DB_NAME", "app_db"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			Schema:          getEnv("DB_SCHEMA", "public"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 15*time.Minute),
			DatabaseURL:     os.Getenv("DATABASE_URL"),
		},
		Redis: RedisConfig{
			Host:       getEnv("REDIS_HOST", "localhost"),
			Port:       getEnv("REDIS_PORT", "6379"),
			Password:   getEnv("REDIS_PASSWORD", ""),
			DB:         getIntEnv("REDIS_DB", 0),
			TTLDefault: getDurationEnv("REDIS_TTL_DEFAULT", 24*time.Hour),
		},
		Auth: AuthConfig{
			JWTAccessSecret:            getEnv("JWT_ACCESS_SECRET", "your_super_secret_access_key_change_in_prod_min_32_chars"),
			JWTRefreshSecret:           getEnv("JWT_REFRESH_SECRET", "your_super_secret_refresh_key_change_in_prod_min_32_chars"),
			JWTAccessExpirationMinutes: getIntEnv("JWT_ACCESS_EXPIRATION_MINUTES", 15),
			JWTRefreshExpirationDays:   getIntEnv("JWT_REFRESH_EXPIRATION_DAYS", 30),
			OTPExpirationMinutes:       getIntEnv("OTP_EXPIRATION_MINUTES", 5),
			OTPRateLimitPerMinute:      getIntEnv("OTP_RATE_LIMIT_PER_MINUTE", 1),
			OTPRateLimitPerHour:        getIntEnv("OTP_RATE_LIMIT_PER_HOUR", 5),
			OTPMasterTestCode:          getEnv("OTP_MASTER_TEST_CODE", "123456"),
		},
		Integrations: IntegrationsConfig{
			SMSProvider:        getEnv("SMS_PROVIDER", "mock"),
			StorageProvider:    getEnv("STORAGE_PROVIDER", "mock"),
			PaymentProvider:    getEnv("PAYMENT_PROVIDER", "mock"),
		
			AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			AWSRegion:          getEnv("AWS_REGION", "eu-central-1"),
			AWSS3BucketName:    getEnv("AWS_S3_BUCKET_NAME", "app-media-dev"),
			AWSS3Endpoint:      getEnv("AWS_S3_ENDPOINT", ""),
		},
		Business: BusinessConfig{
			PlatformDefaultCommissionRate: getFloatEnv("PLATFORM_DEFAULT_COMMISSION_RATE", 0.0300),
			AutoReleaseDeliveryHours:      getIntEnv("AUTO_RELEASE_DELIVERY_HOURS", 24),
			DefaultFireThresholdPercent:   getFloatEnv("DEFAULT_FIRE_THRESHOLD_PERCENT", 5.00),
		},
		Worker: WorkerConfig{
			OutboxPollerIntervalMs: getIntEnv("OUTBOX_POLLER_INTERVAL_MS", 1000),
			OutboxBatchSize:        getIntEnv("OUTBOX_BATCH_SIZE", 50),
		},
	}

	return cfg
}

// loadDotEnv parses a key=value formatted .env file and sets OS env vars if not already defined.
func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])

			// Strip inline comments if not enclosed in quotes
			if !strings.HasPrefix(val, `"`) && !strings.HasPrefix(val, `'`) {
				if commentIdx := strings.Index(val, "#"); commentIdx != -1 {
					val = strings.TrimSpace(val[:commentIdx])
				}
			}
			val = strings.Trim(val, `"'`)
			val = os.ExpandEnv(val)

			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
	_ = scanner.Err()
}

// getEnv fetches string env variable or returns fallback.
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// getIntEnv fetches integer env variable or returns fallback.
func getIntEnv(key string, fallback int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

// getFloatEnv fetches float64 env variable or returns fallback.
func getFloatEnv(key string, fallback float64) float64 {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return fallback
	}
	return val
}

// getBoolEnv fetches boolean env variable or returns fallback.
func getBoolEnv(key string, fallback bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}

// getDurationEnv fetches time.Duration env variable or returns fallback.
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return fallback
	}
	d, err := time.ParseDuration(valStr)
	if err != nil {
		return fallback
	}
	return d
}
