package models

import (
	"time"
)

// Cart represents a shopping cart.
type Cart struct {
	CartID    string     `json:"cart_id"`
	ClientID  *int64     `json:"client_id"`
	DiscountID *int64    `json:"discount_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem represents an item in the cart.
type CartItem struct {
	CartItemID    int64      `json:"cart_item_id"`
	CartID        string     `json:"cart_id"`
	ProductID     int64      `json:"product_id"`
	SizeType      *SizeType  `json:"size_type"`
	Quantity      int        `json:"quantity"`
	AddedAt       time.Time  `json:"added_at"`
	ProductName   string     `json:"product_name"`
	UnitPrice     float64    `json:"unit_price"`
	StockQuantity int        `json:"stock_quantity"`
	ThumbnailURL  string     `json:"thumbnail_url"`
}

// CartItemResponse represents a simplified cart item for responses.
type CartItemResponse struct {
	CartItemID    int64     `json:"cart_item_id"`
	ProductID     int64     `json:"product_id"`
	SizeType      *SizeType `json:"size_type"`
	Quantity      int       `json:"quantity"`
	AddedAt       time.Time `json:"added_at"`
	ProductName   string    `json:"product_name"`
	ProductType   string    `json:"product_type"`
	UnitPrice     float64   `json:"unit_price"`
	StockQuantity int       `json:"stock_quantity"`
	ThumbnailURL  string    `json:"thumbnail_url"`
}

// CartWithItems represents a cart with its items.
type CartWithItems struct {
	Cart        Cart       `json:"cart"`
	Items       []CartItem `json:"items"`
	Subtotal    float64    `json:"subtotal"`
	Discount    float64    `json:"discount"`
	Total       float64    `json:"total"`
}