package middleware

import (
	"net/http"

	"nourish-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerMiddleware handles panics and errors
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Internal server error",
				Details: err,
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Success: false,
				Error:   "Internal server error",
			})
		}
		c.Abort()
	})
}
