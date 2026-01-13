package models

import (
	"time"
)

type Discount struct {
	DiscountID             int       `json:"discount_id"`
	DiscountCode           *string   `json:"discount_code"`
	DiscountName           string    `json:"discount_name"`
	DiscountType           string    `json:"discount_type"` // enum
	DiscountValue          *float64  `json:"discount_value"`
	BuyQuantity            *int      `json:"buy_quantity"`
	GetQuantity            *int      `json:"get_quantity"`
	FreeProductID          *int      `json:"free_product_id"`
	AppliesToSameProduct   *bool     `json:"applies_to_same_product"`
	MinimumOrderAmount     *float64  `json:"minimum_order_amount"`
	MaximumDiscountAmount  *float64  `json:"maximum_discount_amount"`
	ProductTypeRestriction *string   `json:"product_type_restriction"` // enum
	StartDate              time.Time `json:"start_date"`
	EndDate                time.Time `json:"end_date"`
	UsageLimit             *int      `json:"usage_limit"`
	UsedCount              int       `json:"used_count"`
	IsActive               bool      `json:"is_active"`
	AppliesTo              string    `json:"applies_to"` // enum
	IsAutoApply            bool      `json:"is_auto_apply"`
}
