package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConfigHandler struct {
	DB *pgxpool.Pool
}

func (h ConfigHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/config/public", h.GetPublicConfig)
}

func (h ConfigHandler) GetPublicConfig(c *gin.Context) {
	// Only fetch safe-to-expose keys
	keys := []string{"banner_message", "shipping_fee_ebuy", "shipping_fee_sf"}
	
	result := make(map[string]string)
	
	for _, key := range keys {
		var val string
		err := h.DB.QueryRow(context.Background(), "SELECT config_value FROM sys_config WHERE config_key = $1", key).Scan(&val)
		if err == nil {
			result[key] = val
		}
	}
	
	c.JSON(http.StatusOK, result)
}
