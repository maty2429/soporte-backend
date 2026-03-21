package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Docs     DocsConfig
	Security SecurityConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name    string
	Env     string
	Version string
}

type HTTPConfig struct {
	Host            string
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	MaxHeaderBytes  int
}

type DocsConfig struct {
	Enabled bool
}

type SecurityConfig struct {
	CORSAllowAll         bool
	CORSAllowedOrigins   []string
	TrustedProxies       []string
	RateLimitRequests    int
	RateLimitWindow      time.Duration
	RequestSizeLimitByte int64
}

type DatabaseConfig struct {
	Enabled            bool
	FailFast           bool
	Driver             string
	Host               string
	Port               string
	User               string
	Password           string
	Name               string
	SSLMode            string
	TimeZone           string
	DSN                string // construido a partir de los campos anteriores
	MaxIdleConns       int
	MaxOpenConns       int
	ConnMaxLifetime    time.Duration
	ConnMaxIdleTime    time.Duration
	QueryTimeout       time.Duration
	SlowQueryThreshold time.Duration
}

func Load() (Config, error) {
	dbEnabled, err := getBool("DB_ENABLED", false)
	if err != nil {
		return Config{}, err
	}

	dbFailFast, err := getBool("DB_FAIL_FAST", false)
	if err != nil {
		return Config{}, err
	}

	corsAllowAll, err := getBool("SECURITY_CORS_ALLOW_ALL", true)
	if err != nil {
		return Config{}, err
	}

	docsEnabled, err := getBool("DOCS_ENABLED", true)
	if err != nil {
		return Config{}, err
	}

	readTimeout, err := getDuration("HTTP_READ_TIMEOUT", 15*time.Second)
	if err != nil {
		return Config{}, err
	}

	writeTimeout, err := getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second)
	if err != nil {
		return Config{}, err
	}

	idleTimeout, err := getDuration("HTTP_IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return Config{}, err
	}

	shutdownTimeout, err := getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}

	rateLimitWindow, err := getDuration("SECURITY_RATE_LIMIT_WINDOW", time.Minute)
	if err != nil {
		return Config{}, err
	}

	connMaxLifetime, err := getDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	if err != nil {
		return Config{}, err
	}

	connMaxIdleTime, err := getDuration("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)
	if err != nil {
		return Config{}, err
	}

	queryTimeout, err := getDuration("DB_QUERY_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}

	slowQueryThreshold, err := getDuration("DB_SLOW_QUERY_THRESHOLD", 200*time.Millisecond)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		App: AppConfig{
			Name:    getString("APP_NAME", "soporte"),
			Env:     getString("APP_ENV", defaultEnv()),
			Version: getString("APP_VERSION", "0.1.0"),
		},
		Docs: DocsConfig{
			Enabled: docsEnabled,
		},
		HTTP: HTTPConfig{
			Host:            getString("HTTP_HOST", "0.0.0.0"),
			Port:            getString("HTTP_PORT", "8080"),
			ReadTimeout:     readTimeout,
			WriteTimeout:    writeTimeout,
			IdleTimeout:     idleTimeout,
			ShutdownTimeout: shutdownTimeout,
			MaxHeaderBytes:  getInt("HTTP_MAX_HEADER_BYTES", 1<<20),
		},
		Security: SecurityConfig{
			CORSAllowAll:         corsAllowAll,
			CORSAllowedOrigins:   getStringSlice("SECURITY_CORS_ALLOWED_ORIGINS"),
			TrustedProxies:       getStringSlice("SECURITY_TRUSTED_PROXIES"),
			RateLimitRequests:    getInt("SECURITY_RATE_LIMIT_REQUESTS", 60),
			RateLimitWindow:      rateLimitWindow,
			RequestSizeLimitByte: getInt64("SECURITY_REQUEST_SIZE_LIMIT_BYTES", 1<<20),
		},
		Database: buildDatabaseConfig(dbEnabled, dbFailFast, connMaxLifetime, connMaxIdleTime, queryTimeout, slowQueryThreshold),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func buildDatabaseConfig(enabled, failFast bool, connMaxLifetime, connMaxIdleTime, queryTimeout, slowQueryThreshold time.Duration) DatabaseConfig {
	host := getString("DB_HOST", "localhost")
	port := getString("DB_PORT", "5432")
	user := getString("DB_USER", "")
	password := getString("DB_PASSWORD", "")
	name := getString("DB_NAME", "")
	sslMode := getString("DB_SSLMODE", "disable")
	timeZone := getString("DB_TIMEZONE", "America/Santiago")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, name, sslMode, timeZone,
	)

	return DatabaseConfig{
		Enabled:            enabled,
		FailFast:           failFast,
		Driver:             getString("DB_DRIVER", "postgres"),
		Host:               host,
		Port:               port,
		User:               user,
		Password:           password,
		Name:               name,
		SSLMode:            sslMode,
		TimeZone:           timeZone,
		DSN:                dsn,
		MaxIdleConns:       getInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:       getInt("DB_MAX_OPEN_CONNS", 25),
		ConnMaxLifetime:    connMaxLifetime,
		ConnMaxIdleTime:    connMaxIdleTime,
		QueryTimeout:       queryTimeout,
		SlowQueryThreshold: slowQueryThreshold,
	}
}

func (c Config) Validate() error {
	if c.HTTP.Port == "" {
		return fmt.Errorf("HTTP_PORT is required")
	}

	if c.HTTP.MaxHeaderBytes <= 0 {
		return fmt.Errorf("HTTP_MAX_HEADER_BYTES must be greater than 0")
	}

	if c.Security.RateLimitRequests <= 0 {
		return fmt.Errorf("SECURITY_RATE_LIMIT_REQUESTS must be greater than 0")
	}

	if c.Security.RateLimitWindow <= 0 {
		return fmt.Errorf("SECURITY_RATE_LIMIT_WINDOW must be greater than 0")
	}

	if c.Security.RequestSizeLimitByte <= 0 {
		return fmt.Errorf("SECURITY_REQUEST_SIZE_LIMIT_BYTES must be greater than 0")
	}

	if !c.Database.Enabled {
		return nil
	}

	switch c.Database.Driver {
	case "postgres", "sqlite":
	default:
		return fmt.Errorf("DB_DRIVER must be postgres or sqlite")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required when DB_ENABLED=true")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required when DB_ENABLED=true")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required when DB_ENABLED=true")
	}

	return nil
}

func getString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}

	return parsed
}

func getBool(key string, fallback bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("%s must be a boolean: %w", key, err)
	}

	return parsed, nil
}

func getDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", key, err)
	}

	return parsed, nil
}

// getStringSlice parses a comma-separated env variable into a trimmed slice.
// Returns nil (not an empty slice) when the variable is unset or blank so that
// callers can distinguish "not configured" from "empty list".
func getStringSlice(key string) []string {
	value := os.Getenv(key)
	if value == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}

	return out
}

func defaultEnv() string {
	if IsProduction {
		return "production"
	}
	return "development"
}
