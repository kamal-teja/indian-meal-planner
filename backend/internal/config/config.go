package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Port        string
	Environment string
	MongoURI    string
	
	// JWT Configuration
	JWTSecret    string
	JWTExpiresIn time.Duration
	
	// CORS Configuration
	AllowedOrigins []string
	
	// Rate Limiting
	RateLimitWindowMS     time.Duration
	RateLimitMaxRequests  int
	
	// Logging
	LogLevel  string
	LogFormat string
	
	// Server Configuration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
	
	// Database Configuration
	DatabaseConfig DatabaseConfig
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	MaxPoolSize  int
	MinPoolSize  int
	MaxIdleTime  time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "5000"),
		Environment: getEnv("NODE_ENV", "development"),
		MongoURI:    getEnv("MONGODB_URI", "mongodb://localhost:27017/meal-planner"),
		
		JWTSecret:    getEnv("JWT_SECRET", "fallback-secret-key"),
		JWTExpiresIn: parseDuration(getEnv("JWT_EXPIRES_IN", "7d"), 7*24*time.Hour),
		
		AllowedOrigins: parseOrigins(getEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:5173")),
		
		RateLimitWindowMS:    parseDuration(getEnv("RATE_LIMIT_WINDOW_MS", "900000ms"), 15*time.Minute),
		RateLimitMaxRequests: parseInt(getEnv("RATE_LIMIT_MAX_REQUESTS", "100"), 100),
		
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),
		
		ReadTimeout:    parseDuration(getEnv("READ_TIMEOUT", "30s"), 30*time.Second),
		WriteTimeout:   parseDuration(getEnv("WRITE_TIMEOUT", "30s"), 30*time.Second),
		IdleTimeout:    parseDuration(getEnv("IDLE_TIMEOUT", "120s"), 120*time.Second),
		MaxHeaderBytes: parseInt(getEnv("MAX_HEADER_BYTES", "1048576"), 1048576),
		
		DatabaseConfig: DatabaseConfig{
			MaxPoolSize: parseInt(getEnv("DB_MAX_POOL_SIZE", "10"), 10),
			MinPoolSize: parseInt(getEnv("DB_MIN_POOL_SIZE", "1"), 1),
			MaxIdleTime: parseDuration(getEnv("DB_MAX_IDLE_TIME", "60s"), 60*time.Second),
		},
	}
}

// getEnv gets an environment variable with a fallback default
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseInt parses an integer from string with fallback
func parseInt(s string, defaultValue int) int {
	if value, err := strconv.Atoi(s); err == nil {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration from string with fallback
func parseDuration(s string, defaultValue time.Duration) time.Duration {
	// Handle simple day format like "7d"
	if strings.HasSuffix(s, "d") {
		if days, err := strconv.Atoi(strings.TrimSuffix(s, "d")); err == nil {
			return time.Duration(days) * 24 * time.Hour
		}
	}
	
	if duration, err := time.ParseDuration(s); err == nil {
		return duration
	}
	return defaultValue
}

// parseOrigins parses comma-separated origins
func parseOrigins(s string) []string {
	if s == "" {
		return []string{}
	}
	
	origins := strings.Split(s, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	return origins
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
