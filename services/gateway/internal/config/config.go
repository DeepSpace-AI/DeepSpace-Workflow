package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	NewAPIBaseURL string
	NewAPIKey     string

	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	DBAutoMigrate  bool
	DBResetOnStart bool

	JWTSecret       string
	JWTIssuer       string
	JWTExpiresIn    time.Duration
	JWTCookieName   string
	JWTCookieSecure bool

	KBStoragePath string
	KBMaxUploadMB int
	KBAllowedMIME []string

	EmailEnabled     bool
	EmailFromName    string
	EmailFromAddress string
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPassword     string
	SMTPUseTLS       bool
	EmailTemplateDir string
	RedisEnabled     bool
	RedisURL         string
	RedisQueueKey    string
}

func Load() *Config {
	// Try to load .env file from various locations
	// We ignore errors because the file might not exist (e.g. in production using real env vars)

	// 1. Try root directory (assuming running from services/gateway)
	_ = godotenv.Load("../../.env")
	// 2. Try current directory (assuming running from root or .env exists locally)
	_ = godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		NewAPIBaseURL: getEnv("NEWAPI_BASE_URL", "http://localhost:3000"),
		NewAPIKey:     getEnv("NEWAPI_API_KEY", ""),

		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "deepspace"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		DBAutoMigrate:  getEnvBool("DB_AUTO_MIGRATE", true),
		DBResetOnStart: getEnvBool("DB_RESET_ON_START", false),

		JWTSecret:       getEnv("JWT_SECRET", ""),
		JWTIssuer:       getEnv("JWT_ISSUER", "deepspace"),
		JWTExpiresIn:    time.Duration(getEnvInt("JWT_EXPIRES_IN_SECONDS", 86400)) * time.Second,
		JWTCookieName:   getEnv("JWT_COOKIE_NAME", "dsp_session"),
		JWTCookieSecure: getEnvBool("JWT_COOKIE_SECURE", false),

		KBStoragePath: getEnv("KB_STORAGE_PATH", "./data/kb"),
		KBMaxUploadMB: getEnvInt("KB_MAX_UPLOAD_MB", 25),
		KBAllowedMIME: parseCommaList(getEnv("KB_ALLOWED_MIME", "")),

		EmailEnabled:     getEnvBool("EMAIL_ENABLED", false),
		EmailFromName:    getEnv("EMAIL_FROM_NAME", "DeepSpace"),
		EmailFromAddress: getEnv("EMAIL_FROM_ADDRESS", ""),
		SMTPHost:         getEnv("SMTP_HOST", ""),
		SMTPPort:         getEnvInt("SMTP_PORT", 587),
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPassword:     getEnv("SMTP_PASSWORD", ""),
		SMTPUseTLS:       getEnvBool("SMTP_USE_TLS", true),
		EmailTemplateDir: getEnv("EMAIL_TEMPLATE_DIR", "../../templates"),

		RedisEnabled:  getEnvBool("REDIS_ENABLED", false),
		RedisURL:      getEnv("REDIS_URL", ""),
		RedisQueueKey: getEnv("REDIS_QUEUE_KEY", "email:queue"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return fallback
	}
}

func (c *Config) PostgresDSN() string {
	// Compose DSN from discrete DB_* envs to avoid a single DATABASE_URL.
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func (c *Config) Validate() error {
	if strings.TrimSpace(c.Port) == "" {
		return fmt.Errorf("PORT is required")
	}
	if strings.TrimSpace(c.NewAPIBaseURL) == "" {
		return fmt.Errorf("NEWAPI_BASE_URL is required")
	}
	if strings.TrimSpace(c.DBHost) == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if strings.TrimSpace(c.DBUser) == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if strings.TrimSpace(c.DBName) == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	if strings.TrimSpace(c.JWTSecret) == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if strings.TrimSpace(c.KBStoragePath) == "" {
		return fmt.Errorf("KB_STORAGE_PATH is required")
	}
	if c.KBMaxUploadMB <= 0 {
		return fmt.Errorf("KB_MAX_UPLOAD_MB must be positive")
	}
	if c.EmailEnabled {
		if strings.TrimSpace(c.EmailFromAddress) == "" {
			return fmt.Errorf("EMAIL_FROM_ADDRESS is required")
		}
		if strings.TrimSpace(c.SMTPHost) == "" {
			return fmt.Errorf("SMTP_HOST is required")
		}
		if c.SMTPPort <= 0 {
			return fmt.Errorf("SMTP_PORT must be positive")
		}
	}
	if c.RedisEnabled {
		if strings.TrimSpace(c.RedisURL) == "" {
			return fmt.Errorf("REDIS_URL is required")
		}
		if strings.TrimSpace(c.RedisQueueKey) == "" {
			return fmt.Errorf("REDIS_QUEUE_KEY is required")
		}
	}
	return nil
}

func getEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return fallback
	}
	return parsed
}

func parseCommaList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.ToLower(strings.TrimSpace(part))
		if item == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}

func (c *Config) KBMaxUploadBytes() int64 {
	return int64(c.KBMaxUploadMB) * 1024 * 1024
}
