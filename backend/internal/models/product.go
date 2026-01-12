package models

import (
	"time"
)

// ProductType represents the type of product.
type ProductType string

const (
	ProductTypeFaiachun ProductType = "faiachun"
	ProductTypeBag      ProductType = "bag"
)

// SizeType represents the available sizes for products.
type SizeType string

const (
	SizeTypeVRect    SizeType = "v-rect"
	SizeTypeSquare   SizeType = "square"
	SizeTypeFatVRect SizeType = "fat-v-rect"
)

// Product represents a product in the catalog.
type Product struct {
	ID             int64       `json:"product_id"`
	Name           string      `json:"product_name"`
	Description    *string     `json:"product_description"`
	Type           ProductType `json:"product_type"`
	Hashtag        *string     `json:"hashtag"`
	SKU            string      `json:"sku"`
	Price          float64     `json:"price"`
	CompareAtPrice *float64    `json:"compare_at_price"`
	CostPrice      *float64    `json:"cost_price"`
	Quantity       int         `json:"quantity"`
	Weight         *float64    `json:"weight"`
	Dimensions     *string     `json:"dimensions"` // JSONB stored as string
	AvailableSizes []SizeType  `json:"available_sizes"`
	IsFeatured     bool        `json:"is_featured"`
	IsActive       bool        `json:"is_active"`
	SEOTitle       *string     `json:"seo_title"`
	SEODescription *string     `json:"seo_description"`
	Tags           *string     `json:"tags"` // JSONB stored as string
	CreatedBy      *int64      `json:"created_by"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// ProductImage represents an image associated with a product.
type ProductImage struct {
	ID        int64     `json:"image_id"`
	ProductID int64     `json:"product_id"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	SizeType  *SizeType `json:"size_type"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Category represents a product category.
type Category struct {
	ID          int64  `json:"category_id"`
	Name        string `json:"category_name"`
	Description *string `json:"category_description"`
	ParentID    *int64 `json:"parent_category_id"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// ProductWithDetails represents a product with its images and categories.
type ProductWithDetails struct {
	Product
	Images     []ProductImage `json:"images"`
	Categories []Category     `json:"categories"`
}

// ProductListResponse represents the paginated response for product listing.
type ProductListResponse struct {
	Data []ProductWithDetails `json:"data"`
	Meta PaginationMeta       `json:"meta"`
}

// PaginationMeta contains pagination information.
type PaginationMeta struct {
	Page      int `json:"page"`
	PageSize  int `json:"page_size"`
	Total     int `json:"total"`
	TotalPages int `json:"total_pages,omitempty"`
}