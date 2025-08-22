package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad_DefaultValues(t *testing.T) {
	// Clean environment
	os.Clearenv()

	// Act
	config := Load()

	// Assert
	assert.Equal(t, "5000", config.Port)
	assert.Equal(t, "development", config.Environment)
	assert.Equal(t, "mongodb://localhost:27017/meal-planner", config.MongoURI)
	assert.Equal(t, "fallback-secret-key", config.JWTSecret)
	assert.Equal(t, 7*24*time.Hour, config.JWTExpiresIn) // 7 days
	assert.Equal(t, []string{"http://localhost:3000", "http://localhost:5173"}, config.AllowedOrigins)
	assert.Equal(t, 15*time.Minute, config.RateLimitWindowMS)
	assert.Equal(t, 100, config.RateLimitMaxRequests)
	assert.Equal(t, "info", config.LogLevel)
	assert.Equal(t, "json", config.LogFormat)
	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 120*time.Second, config.IdleTimeout)
	assert.Equal(t, 1048576, config.MaxHeaderBytes) // 1 MB
}

func TestLoad_CustomEnvironmentValues(t *testing.T) {
	// Arrange - Set custom environment variables
	os.Setenv("PORT", "8080")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("MONGODB_URI", "mongodb://prod-host:27017/prod-db")
	os.Setenv("JWT_SECRET", "custom-jwt-secret")
	os.Setenv("JWT_EXPIRES_IN", "48h")
	os.Setenv("ALLOWED_ORIGINS", "https://example.com,https://app.example.com")
	os.Setenv("RATE_LIMIT_WINDOW_MS", "5m")
	os.Setenv("RATE_LIMIT_MAX_REQUESTS", "200")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_FORMAT", "text")
	
	defer func() {
		// Clean up after test
		os.Clearenv()
	}()

	// Act
	config := Load()

	// Assert
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "production", config.Environment)
	assert.Equal(t, "mongodb://prod-host:27017/prod-db", config.MongoURI)
	assert.Equal(t, "custom-jwt-secret", config.JWTSecret)
	assert.Equal(t, 48*time.Hour, config.JWTExpiresIn)
	assert.Equal(t, []string{"https://example.com", "https://app.example.com"}, config.AllowedOrigins)
	assert.Equal(t, 5*time.Minute, config.RateLimitWindowMS)
	assert.Equal(t, 200, config.RateLimitMaxRequests)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, "text", config.LogFormat)
}

func TestLoad_DatabaseConfig(t *testing.T) {
	// Arrange
	os.Setenv("DB_MAX_POOL_SIZE", "50")
	os.Setenv("DB_MIN_POOL_SIZE", "5")
	os.Setenv("DB_MAX_IDLE_TIME", "30m")
	
	defer func() {
		os.Clearenv()
	}()

	// Act
	config := Load()

	// Assert
	assert.Equal(t, 50, config.DatabaseConfig.MaxPoolSize)
	assert.Equal(t, 5, config.DatabaseConfig.MinPoolSize)
	assert.Equal(t, 30*time.Minute, config.DatabaseConfig.MaxIdleTime)
}

func TestLoad_ServerTimeouts(t *testing.T) {
	// Arrange
	os.Setenv("READ_TIMEOUT", "15s")
	os.Setenv("WRITE_TIMEOUT", "20s")
	os.Setenv("IDLE_TIMEOUT", "180s")
	os.Setenv("MAX_HEADER_BYTES", "2097152") // 2 MB
	
	defer func() {
		os.Clearenv()
	}()

	// Act
	config := Load()

	// Assert
	assert.Equal(t, 15*time.Second, config.ReadTimeout)
	assert.Equal(t, 20*time.Second, config.WriteTimeout)
	assert.Equal(t, 180*time.Second, config.IdleTimeout)
	assert.Equal(t, 2097152, config.MaxHeaderBytes)
}

func TestLoad_InvalidDurationValues(t *testing.T) {
	// Arrange - Set invalid duration values
	os.Setenv("JWT_EXPIRES_IN", "invalid-duration")
	os.Setenv("RATE_LIMIT_WINDOW_MS", "not-a-duration")
	
	defer func() {
		os.Clearenv()
	}()

	// Act
	config := Load()

	// Assert - Should fallback to default values when parsing fails
	assert.Equal(t, 7*24*time.Hour, config.JWTExpiresIn) // 7 days default
	assert.Equal(t, 15*time.Minute, config.RateLimitWindowMS)
}

func TestLoad_InvalidIntegerValues(t *testing.T) {
	// Arrange - Set invalid integer values
	os.Setenv("RATE_LIMIT_MAX_REQUESTS", "not-a-number")
	os.Setenv("DB_MAX_POOL_SIZE", "invalid")
	os.Setenv("MAX_HEADER_BYTES", "not-an-int")
	
	defer func() {
		os.Clearenv()
	}()

	// Act
	config := Load()

	// Assert - Should fallback to default values when parsing fails
	assert.Equal(t, 100, config.RateLimitMaxRequests)
	assert.Equal(t, 10, config.DatabaseConfig.MaxPoolSize) // default value
	assert.Equal(t, 1048576, config.MaxHeaderBytes) // default 1 MB
}

func TestLoad_EmptyAllowedOrigins(t *testing.T) {
	// Arrange
	os.Setenv("ALLOWED_ORIGINS", "")
	
	defer func() {
		os.Clearenv()
	}()

	// Act
	config := Load()

	// Assert - Should use default origins when empty
	assert.Equal(t, []string{"http://localhost:3000", "http://localhost:5173"}, config.AllowedOrigins)
}

func TestLoad_SingleAllowedOrigin(t *testing.T) {
	// Arrange
	os.Setenv("ALLOWED_ORIGINS", "https://single-origin.com")
	
	defer func() {
		os.Clearenv()
	}()

	// Act
	config := Load()

	// Assert
	assert.Equal(t, []string{"https://single-origin.com"}, config.AllowedOrigins)
}

func TestConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{"production environment", "production", true},
		{"prod environment", "prod", true},
		{"development environment", "development", false},
		{"test environment", "test", false},
		{"empty environment", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Environment: tt.environment}
			assert.Equal(t, tt.expected, config.IsProduction())
		})
	}
}

func TestConfig_IsDevelopment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		expected    bool
	}{
		{"development environment", "development", true},
		{"dev environment", "dev", true},
		{"production environment", "production", false},
		{"test environment", "test", false},
		{"empty environment", "", true}, // defaults to development
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Environment: tt.environment}
			assert.Equal(t, tt.expected, config.IsDevelopment())
		})
	}
}