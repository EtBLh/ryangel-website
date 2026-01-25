package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration for the API server.
type Config struct {
	AppHost string
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	TokenTTLMinutes int
	LogLevel        string

	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioPhoneNumber string
	SkipSMSSending   bool
	MediaStoragePath string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

// FromEnv constructs Config from environment variables with sensible defaults.
func FromEnv() (*Config, error) {
	cfg := &Config{
		AppHost:    getEnv("APP_HOST", "0.0.0.0"),
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     getEnv("DB_NAME", "postgres"),
		SSLMode:    getEnv("DB_SSLMODE", "disable"),
		TokenTTLMinutes: getEnvAsInt("TOKEN_TTL_MINUTES", 1440),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		TwilioAccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioPhoneNumber: os.Getenv("TWILIO_PHONE_NUMBER"),
		SkipSMSSending:   getEnvAsBool("SKIP_SMS_SENDING", false),
		MediaStoragePath: getEnv("MEDIA_STORAGE_PATH", "./media"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "https://ryangel.com/api/auth/google/callback"),
	}

	if cfg.DBPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required")
	}

	return cfg, nil
}

// HTTPAddr renders the listen address for the HTTP server.
func (c *Config) HTTPAddr() string {
	return fmt.Sprintf("%s:%s", c.AppHost, c.AppPort)
}

// DatabaseURL builds a DSN suitable for pgxpool.
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(c.DBUser),
		url.QueryEscape(c.DBPassword),
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}

// TokenTTL returns the duration a bearer token remains valid.
func (c *Config) TokenTTL() time.Duration {
	return time.Duration(c.TokenTTLMinutes) * time.Minute
}
