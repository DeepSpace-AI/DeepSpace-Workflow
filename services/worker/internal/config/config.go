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
	EmailEnabled     bool
	EmailFromName    string
	EmailFromAddress string
	SMTPHost         string
	SMTPPort         int
	SMTPUser         string
	SMTPPassword     string
	SMTPUseTLS       bool
	EmailTemplateDir string

	RedisURL      string
	RedisQueueKey string
	RedisDeadKey  string

	PollTimeout time.Duration
	RetryMax    int
}

func Load() *Config {
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load()

	return &Config{
		EmailEnabled:     getEnvBool("EMAIL_ENABLED", false),
		EmailFromName:    getEnv("EMAIL_FROM_NAME", "DeepSpace"),
		EmailFromAddress: getEnv("EMAIL_FROM_ADDRESS", ""),
		SMTPHost:         getEnv("SMTP_HOST", ""),
		SMTPPort:         getEnvInt("SMTP_PORT", 587),
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPassword:     getEnv("SMTP_PASSWORD", ""),
		SMTPUseTLS:       getEnvBool("SMTP_USE_TLS", true),
		EmailTemplateDir: getEnv("EMAIL_TEMPLATE_DIR", "../../templates"),

		RedisURL:      getEnv("REDIS_URL", ""),
		RedisQueueKey: getEnv("REDIS_QUEUE_KEY", "email:queue"),
		RedisDeadKey:  getEnv("REDIS_DEAD_KEY", "email:dead"),

		PollTimeout: time.Duration(getEnvInt("EMAIL_QUEUE_TIMEOUT", 10)) * time.Second,
		RetryMax:    getEnvInt("EMAIL_RETRY_MAX", 5),
	}
}

func (c *Config) Validate() error {
	if !c.EmailEnabled {
		return fmt.Errorf("EMAIL_ENABLED must be true")
	}
	if strings.TrimSpace(c.EmailFromAddress) == "" {
		return fmt.Errorf("EMAIL_FROM_ADDRESS is required")
	}
	if strings.TrimSpace(c.SMTPHost) == "" {
		return fmt.Errorf("SMTP_HOST is required")
	}
	if c.SMTPPort <= 0 {
		return fmt.Errorf("SMTP_PORT must be positive")
	}
	if strings.TrimSpace(c.RedisURL) == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if strings.TrimSpace(c.RedisQueueKey) == "" {
		return fmt.Errorf("REDIS_QUEUE_KEY is required")
	}
	if strings.TrimSpace(c.RedisDeadKey) == "" {
		return fmt.Errorf("REDIS_DEAD_KEY is required")
	}
	if c.PollTimeout <= 0 {
		return fmt.Errorf("EMAIL_QUEUE_TIMEOUT must be positive")
	}
	if c.RetryMax <= 0 {
		return fmt.Errorf("EMAIL_RETRY_MAX must be positive")
	}
	return nil
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
