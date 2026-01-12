package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler verifies application dependencies are reachable.
type HealthHandler struct {
	DB *pgxpool.Pool
}

// Register wires the health routes onto the router.
func (h HealthHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/healthz", h.Get)
}

// Get responds with liveness/readiness information.
func (h HealthHandler) Get(c *gin.Context) {
	status := http.StatusOK
	dbState := "ok"

	if h.DB != nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := h.DB.Ping(ctx); err != nil {
			status = http.StatusServiceUnavailable
			dbState = "error"
			c.JSON(status, gin.H{
				"status":   "degraded",
				"database": dbState,
				"error":    err.Error(),
			})
			return
		}
	}

	c.JSON(status, gin.H{
		"status":    "ok",
		"database":  dbState,
		"timestamp": time.Now().UTC(),
	})
}
