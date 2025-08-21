package middleware

import (
	"net/http"

	"nourish-backend/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSMiddleware creates CORS middleware
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	config := cors.Config{
		AllowAllOrigins: false,
		AllowOrigins:    cfg.AllowedOrigins,
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
		},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	// Allow all origins in development
	if cfg.IsDevelopment() {
		config.AllowAllOrigins = true
		config.AllowOrigins = nil
	}

	return cors.New(config)
}
