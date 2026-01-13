package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ryangel/ryangel-backend/internal/repository"
)

// EbuyStoreHandler handles ebuy store-related HTTP requests.
type EbuyStoreHandler struct {
	Repo *repository.EbuyStoreRepository
}

// Register wires the ebuy store routes onto the router.
func (h EbuyStoreHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/ebuystores", h.GetEbuyStores)
	rg.GET("/ebuystores/:store_id", h.GetEbuyStore)
}

// GetEbuyStores handles GET /stores.
func (h EbuyStoreHandler) GetEbuyStores(c *gin.Context) {
	stores, err := h.Repo.GetEbuyStores(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ebuy stores"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stores})
}

// GetEbuyStore handles GET /stores/{store_id}.
func (h EbuyStoreHandler) GetEbuyStore(c *gin.Context) {
	storeID := c.Param("store_id")

	store, err := h.Repo.GetEbuyStoreByID(c.Request.Context(), storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get ebuy store"})
		return
	}
	if store == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Store not found"})
		return
	}

	c.JSON(http.StatusOK, store)
}