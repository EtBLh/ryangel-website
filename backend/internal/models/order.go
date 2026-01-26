package models

import (
	"time"
)

// OrderStatus represents the status of an order.
type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

// PaymentMethod represents the payment method used.
type PaymentMethod string

const (
	PaymentMethodMPay       PaymentMethod = "mpay"
	PaymentMethodBOC        PaymentMethod = "boc"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

// PaymentStatus represents the status of payment.
type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Order represents an order in the system.
type Order struct {
	OrderID          int64          `json:"order_id"`
	OrderNumber      string         `json:"order_number"`
	ClientID         int64          `json:"client_id"`
	OrderStatus      OrderStatus    `json:"order_status"`
	SubtotalAmount   float64        `json:"subtotal_amount"`
	DiscountAmount   float64        `json:"discount_amount"`
	ShippingAmount   float64        `json:"shipping_amount"`
	TaxAmount        float64        `json:"tax_amount"`
	TotalAmount      float64        `json:"total_amount"`
	DiscountID       *int64         `json:"discount_id"`
	DiscountCode     *string        `json:"discount_code"`
	ShippingAddressID *int64        `json:"shipping_address_id"`
	EbuyStoreID      *string        `json:"ebuy_store_id"`
	ContactPhone     string         `json:"contact_phone"`
	PaymentMethod    PaymentMethod  `json:"payment_method"`
	PaymentStatus    PaymentStatus  `json:"payment_status"`
	PaymentReference *string        `json:"payment_reference"`
	TrackingNumber   *string        `json:"tracking_number"`
	ShippingCarrier  *string        `json:"shipping_carrier"`
	OrderDate        time.Time      `json:"order_date"`
	ConfirmedAt      *time.Time     `json:"confirmed_at"`
	ShippedAt        *time.Time     `json:"shipped_at"`
	DeliveredAt      *time.Time     `json:"delivered_at"`
	CancelledAt      *time.Time     `json:"cancelled_at"`
	CustomerNotes    *string        `json:"customer_notes"`
	AdminNotes       *string        `json:"admin_notes"`
	PaymentProof     *string        `json:"payment_proof"`
	ClientName       string         `json:"client_name"`
	ClientPhone      string         `json:"client_phone"`
	EbuyStoreName    *string        `json:"ebuy_store_name"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	OrderItemID     int64         `json:"order_item_id"`
	OrderID         int64         `json:"order_id"`
	ProductID       int64         `json:"product_id"`
	Quantity        int           `json:"quantity"`
	UnitPrice       float64       `json:"unit_price"`
	DiscountAmount  float64       `json:"discount_amount"`
	TotalPrice      float64       `json:"total_price"`
	ProductName     string        `json:"product_name"`
	ProductType     ProductType   `json:"product_type"`
	ProductSKU      string        `json:"product_sku"`
	SizeType        *SizeType     `json:"size_type"`
	IsFreeItem      bool          `json:"is_free_item"`
	ParentDiscountID *int64       `json:"parent_discount_id"`
	ProductImage    *string       `json:"product_image"`
}

// OrderWithItems represents an order with its items.
type OrderWithItems struct {
	Order Order       `json:"order"`
	Items []OrderItem `json:"items"`
}