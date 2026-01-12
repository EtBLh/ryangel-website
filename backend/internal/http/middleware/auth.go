package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ryangel/ryangel-backend/internal/models"
	authsvc "github.com/ryangel/ryangel-backend/internal/services/auth"
)

const (
	adminContextKey  = "auth.admin"
	clientContextKey = "auth.client"
)

// AdminAuth ensures the request has a valid admin bearer token.
func AdminAuth(service *authsvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c.GetHeader("Authorization"))
		if token == "" {
			abortUnauthorized(c, "AUTH_INVALID_CREDENTIALS", "Missing bearer token.")
			return
		}

		admin, err := service.ValidateAdminToken(c.Request.Context(), token)
		if err != nil {
			if errors.Is(err, authsvc.ErrInvalidToken) {
				abortUnauthorized(c, "AUTH_TOKEN_EXPIRED", "Token expired. Refresh login.")
			} else {
				abortUnauthorized(c, "AUTH_INVALID_CREDENTIALS", "Unable to verify bearer token.")
			}
			return
		}

		c.Set(adminContextKey, admin)
		c.Next()
	}
}

// ClientAuth ensures the request has a valid client bearer token.
func ClientAuth(service *authsvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractBearer(c.GetHeader("Authorization"))
		if token == "" {
			abortUnauthorized(c, "AUTH_INVALID_CREDENTIALS", "Missing bearer token.")
			return
		}

		client, err := service.ValidateClientToken(c.Request.Context(), token)
		if err != nil {
			if errors.Is(err, authsvc.ErrInvalidToken) {
				abortUnauthorized(c, "AUTH_TOKEN_EXPIRED", "Token expired. Refresh login.")
			} else {
				abortUnauthorized(c, "AUTH_INVALID_CREDENTIALS", "Unable to verify bearer token.")
			}
			return
		}

		c.Set(clientContextKey, client)
		c.Next()
	}
}

// AdminFromContext fetches the authenticated admin from the Gin context.
func AdminFromContext(c *gin.Context) (*models.Admin, bool) {
	value, ok := c.Get(adminContextKey)
	if !ok {
		return nil, false
	}
	admin, ok := value.(*models.Admin)
	return admin, ok
}

// ClientFromContext fetches the authenticated client from the Gin context.
func ClientFromContext(c *gin.Context) (*models.Client, bool) {
	value, ok := c.Get(clientContextKey)
	if !ok {
		return nil, false
	}
	client, ok := value.(*models.Client)
	return client, ok
}

func extractBearer(header string) string {
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func abortUnauthorized(c *gin.Context, code, message string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
