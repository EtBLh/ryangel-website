package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/ryangel/ryangel-backend/internal/config"
	"github.com/ryangel/ryangel-backend/internal/models"
	"github.com/ryangel/ryangel-backend/internal/repository"
)

// ProductHandler handles product-related HTTP requests.
type ProductHandler struct {
	Repo   *repository.ProductRepository
	Config *config.Config
}

// Register wires the product routes onto the router.
func (h ProductHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/products", h.ListProducts)
	rg.GET("/products/:product_id", h.GetProduct)
	rg.GET("/products/:product_id/images", h.GetProductImages)
}

// ListProducts handles GET /products with filtering, sorting, and pagination.
func (h ProductHandler) ListProducts(c *gin.Context) {
	// Parse query parameters
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	pageSize := 20
	if sizeStr := c.Query("page_size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil && s > 0 && s <= 100 {
			pageSize = s
		}
	}

	filters := repository.ProductFilters{
		Query: c.Query("q"),
		SKU:   c.Query("sku"),
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			filters.CategoryID = &categoryID
		}
	}

	if productTypeStr := c.Query("product_type"); productTypeStr != "" {
		if productType := models.ProductType(productTypeStr); productType == models.ProductTypeFaiachun || productType == models.ProductTypeBag {
			filters.ProductType = &productType
		}
	}

	if isFeaturedStr := c.Query("is_featured"); isFeaturedStr != "" {
		if isFeatured, err := strconv.ParseBool(isFeaturedStr); err == nil {
			filters.IsFeatured = &isFeatured
		}
	}

	if priceMinStr := c.Query("price_min"); priceMinStr != "" {
		if priceMin, err := strconv.ParseFloat(priceMinStr, 64); err == nil {
			filters.PriceMin = &priceMin
		}
	}

	if priceMaxStr := c.Query("price_max"); priceMaxStr != "" {
		if priceMax, err := strconv.ParseFloat(priceMaxStr, 64); err == nil {
			filters.PriceMax = &priceMax
		}
	}

	// Parse sort parameter (format: field or field|-field for desc)
	sort := repository.ProductSort{Field: "created_at", Order: "desc"}
	if sortStr := c.Query("sort"); sortStr != "" {
		if strings.HasPrefix(sortStr, "-") {
			sort.Field = strings.TrimPrefix(sortStr, "-")
			sort.Order = "desc"
		} else {
			sort.Field = sortStr
			sort.Order = "asc"
		}
	}

	// Fetch products
	products, total, err := h.Repo.ListProducts(c.Request.Context(), filters, sort, page, pageSize)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch products.", nil)
		return
	}

	// Calculate total pages
	totalPages := (total + pageSize - 1) / pageSize

	response := models.ProductListResponse{
		Data: products,
		Meta: models.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetProduct handles GET /products/{product_id}.
func (h ProductHandler) GetProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid product ID.", nil)
		return
	}

	product, err := h.Repo.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		if err == repository.ErrNotFound {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "Product not found.", nil)
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch product.", nil)
		return
	}

	c.JSON(http.StatusOK, product)
}

// ProductImageResponse represents the response for product images API.
type ProductImageResponse struct {
	URL       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
	AltText   string `json:"alt_text"`
	SortOrder int    `json:"sort_order"`
}

// GetProductImages handles GET /products/{product_id}/images.
func (h ProductHandler) GetProductImages(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		writeError(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid product ID.", nil)
		return
	}

	// Check if product exists
	_, err = h.Repo.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		if err == repository.ErrNotFound {
			writeError(c, http.StatusNotFound, "NOT_FOUND", "Product not found.", nil)
			return
		}
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch product.", nil)
		return
	}

	// Get product images
	images, err := h.Repo.GetProductImages(c.Request.Context(), productID)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch product images.", nil)
		return
	}

	// Convert to response format with URLs
	response := make([]ProductImageResponse, len(images))
	for i, img := range images {
		response[i] = ProductImageResponse{
			URL:       img.URL,
			IsPrimary: img.IsPrimary,
			AltText:   img.AltText,
			SortOrder: img.SortOrder,
		}
	}

	c.JSON(http.StatusOK, response)
}