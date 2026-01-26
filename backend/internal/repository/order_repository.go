package repository

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryangel/ryangel-backend/internal/models"
)

type CreateOrderParams struct {
	ClientID    int64
	EbuyStoreID string
	Name        string
	Phone       string // Contact phone for this order
	Email       string
	Instagram   string
	ProofPath   string
	CartID      string // Optional: explicit cart ID
}

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, params CreateOrderParams) (*models.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// 1. Get Cart
	var cartID string
	var errQuery error

	if params.CartID != "" {
		// Use provided Cart ID if verified
		errQuery = tx.QueryRow(ctx, "SELECT cart_id FROM cart WHERE cart_id = $1 AND client_id = $2", params.CartID, params.ClientID).Scan(&cartID)
	} else {
		// Fallback: pick the most recently updated cart
		errQuery = tx.QueryRow(ctx, "SELECT cart_id FROM cart WHERE client_id = $1 ORDER BY updated_at DESC LIMIT 1", params.ClientID).Scan(&cartID)
	}
	
	if errQuery != nil {
		return nil, fmt.Errorf("cart not found or empty")
	}

	// 2. Get Cart Items
	queryItems := `
		SELECT ci.product_id, ci.size_type, ci.quantity, p.price, p.product_name, p.product_type, p.sku
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.product_id
		WHERE ci.cart_id = $1
	`
	rows, err := tx.Query(ctx, queryItems, cartID)
	if err != nil {
		return nil, err
	}
	
	type cartItem struct {
		ProductID   int64
		SizeType    *string
		Quantity    int
		Price       float64
		ProductName string
		ProductType string
		SKU         string
	}
	var items []cartItem
	var subtotal float64
	
	for rows.Next() {
		var i cartItem
		if err := rows.Scan(&i.ProductID, &i.SizeType, &i.Quantity, &i.Price, &i.ProductName, &i.ProductType, &i.SKU); err != nil {
			rows.Close()
			return nil, err
		}
		items = append(items, i)
		subtotal += i.Price * float64(i.Quantity)
	}
	rows.Close()

	if len(items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculate Discounts & Shipping (Replicating logic from CartHandler)
	itemDiscountAmount := 0.0
	finalShippingFee := 5.0

	// Fetch active auto-apply discounts within the transaction
	discountQuery := `
		SELECT 
			discount_id, discount_code, discount_name, discount_type, discount_value,
			buy_quantity, get_quantity, product_type_restriction, is_auto_apply,
			start_date, end_date, applies_to, is_active
		FROM discounts 
		WHERE is_active = true AND is_auto_apply = true 
		  AND NOW() BETWEEN start_date AND end_date
	`
	dRows, err := tx.Query(ctx, discountQuery)
	if err != nil {
		// Log error but continue (fail safe: no discount)
		fmt.Printf("Error fetching discounts during order creation: %v\n", err)
	} else {
		defer dRows.Close()
		var discounts []models.Discount
		for dRows.Next() {
			var d models.Discount
			var typeStr string
			var restrictionStr *string
			var appliesToStr string
			
			err := dRows.Scan(
				&d.DiscountID, &d.DiscountCode, &d.DiscountName, &typeStr, &d.DiscountValue,
				&d.BuyQuantity, &d.GetQuantity, &restrictionStr, &d.IsAutoApply,
				&d.StartDate, &d.EndDate, &appliesToStr, &d.IsActive,
			)
			if err == nil {
				d.DiscountType = typeStr
				if restrictionStr != nil {
					d.ProductTypeRestriction = restrictionStr
				}
				// ignoring appliesTo enum conversion for now as it's string in model usually
				discounts = append(discounts, d)
			}
		}
		
		// Apply discounts logic
		for _, d := range discounts {
			if d.DiscountType == "bxgy" {
				if d.ProductTypeRestriction != nil {
					restriction := *d.ProductTypeRestriction
					var applicablePrices []float64
					
					for _, item := range items {
						// Match product type using string comparison
						if item.ProductType == string(restriction) {
							for k := 0; k < item.Quantity; k++ {
								applicablePrices = append(applicablePrices, item.Price)
							}
						}
					}
					
					if len(applicablePrices) > 0 && d.BuyQuantity != nil && d.GetQuantity != nil {
						buy := *d.BuyQuantity
						get := *d.GetQuantity
						groupSize := buy + get
						
						sort.Float64s(applicablePrices) // sort asc
						
						numGroups := len(applicablePrices) / groupSize
						numFree := numGroups * get
						
						for i := 0; i < numFree; i++ {
							itemDiscountAmount += applicablePrices[i]
						}
					}
				}
			} else if d.DiscountType == "free_shipping" {
				faiachunCount := 0
				for _, item := range items {
					if item.ProductType == "faiachun" {
						faiachunCount += item.Quantity
					}
				}
				if faiachunCount >= 4 {
					finalShippingFee = 0.0
				}
			}
		}
	}

	// 3. Update Client Info (Name, Email) if provided
	if params.Name != "" || params.Email != "" {
		_, err = tx.Exec(ctx, `
			UPDATE client 
			SET username = COALESCE(NULLIF($2, ''), username), 
			    email = COALESCE(NULLIF($3, ''), email) 
			WHERE client_id = $1`, params.ClientID, params.Name, params.Email)
		if err != nil {
			return nil, err
		}
	}

	// 4. Create Order
	// Generate Order Number: ORD-YYYYMMDD-Random
	orderNum := fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%100000)

	customerNotes := fmt.Sprintf("Store: %s\nContact: %s\nIG: %s\nEmail: %s", params.EbuyStoreID, params.Name, params.Instagram, params.Email)

	shippingAmount := finalShippingFee
	discountAmount := itemDiscountAmount
	// Ensure non-negative
	if discountAmount > subtotal {
		discountAmount = subtotal
	}
	
	totalAmount := subtotal - discountAmount + shippingAmount

	var orderID int64
	var orderDate time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (
			order_number, client_id, order_status, 
			subtotal_amount, discount_amount, shipping_amount, tax_amount, total_amount,
			ebuy_store_id, payment_method, payment_status, customer_notes, order_date, contact_phone
		) VALUES (
			$1, $2, 'pending', 
			$3, $4, $5, 0, $6, 
			$7, 'mpay', 'pending', $8, NOW(), $9
		) RETURNING order_id, order_date`,
		orderNum, params.ClientID, subtotal, discountAmount, shippingAmount, totalAmount,
		params.EbuyStoreID, customerNotes, params.Phone,
	).Scan(&orderID, &orderDate)
	if err != nil {
		return nil, err
	}

	// 5. Insert Order Items
	_, err = tx.Prepare(ctx, "insert_order_item", `
		INSERT INTO order_items (
			order_id, product_id, quantity, unit_price, discount_amount, total_price,
			product_name, product_type, product_sku, size_type
		) VALUES ($1, $2, $3, $4, 0, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		return nil, err
	}
	
	for _, i := range items {
		totalPrice := i.Price * float64(i.Quantity)
		_, err := tx.Exec(ctx, "insert_order_item", 
			orderID, i.ProductID, i.Quantity, i.Price, totalPrice,
			i.ProductName, i.ProductType, i.SKU, i.SizeType,
		)
		if err != nil {
			return nil, err
		}
	}

	// 6. Handle Payment Proof
	if params.ProofPath != "" {
		_, err = tx.Exec(ctx, `
			INSERT INTO payment_proofs (
				order_id, client_id, payment_method, amount, proof_path, status, created_at
			) VALUES ($1, $2, 'mpay', $3, $4, 'submitted', NOW())`,
			orderID, params.ClientID, subtotal, params.ProofPath,
		)
		if err != nil {
			return nil, err
		}
	}

	// 7. Clear Cart
	_, err = tx.Exec(ctx, "DELETE FROM cart_items WHERE cart_id = $1", cartID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &models.Order{
		OrderID:     orderID,
		OrderNumber: orderNum,
		TotalAmount: subtotal,
		OrderDate:   orderDate,
		OrderStatus: models.OrderStatusPending,
	}, nil
}

func (r *OrderRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models.Order, error) {
	const query = `
		SELECT o.order_id, o.order_number, o.client_id, o.order_status, o.subtotal_amount, 
               o.discount_amount, o.shipping_amount, o.tax_amount, o.total_amount, 
               o.discount_id, o.discount_code, o.shipping_address_id, o.ebuy_store_id,
               o.payment_method, o.payment_status, o.payment_reference, o.tracking_number,
               o.shipping_carrier, o.order_date, o.confirmed_at, o.shipped_at, o.delivered_at,
               o.cancelled_at, o.customer_notes, o.admin_notes, COALESCE(o.contact_phone, ''),
			   (SELECT proof_path FROM payment_proofs WHERE order_id = o.order_id ORDER BY created_at DESC LIMIT 1) as payment_proof,
			   s.store_name
		FROM orders o
        LEFT JOIN ebuy_store s ON o.ebuy_store_id = s.store_id
		WHERE o.client_id = $1
		ORDER BY o.order_date DESC`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.OrderID, &o.OrderNumber, &o.ClientID, &o.OrderStatus, &o.SubtotalAmount,
            &o.DiscountAmount, &o.ShippingAmount, &o.TaxAmount, &o.TotalAmount,
            &o.DiscountID, &o.DiscountCode, &o.ShippingAddressID, &o.EbuyStoreID,
            &o.PaymentMethod, &o.PaymentStatus, &o.PaymentReference, &o.TrackingNumber,
            &o.ShippingCarrier, &o.OrderDate, &o.ConfirmedAt, &o.ShippedAt, &o.DeliveredAt,
            &o.CancelledAt, &o.CustomerNotes, &o.AdminNotes, &o.ContactPhone, &o.PaymentProof,
            &o.EbuyStoreName,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, rows.Err()
}

type DashboardStats struct {
	TotalOrders   int64   `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	PendingOrders int64   `json:"pending_orders"`
}

func (r *OrderRepository) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	stats := &DashboardStats{}
	const query = `
		SELECT 
			COUNT(*), 
			COALESCE(SUM(total_amount), 0),
			COUNT(*) FILTER (WHERE order_status = 'pending')
		FROM orders
		WHERE order_status != 'cancelled'`

	err := r.db.QueryRow(ctx, query).Scan(&stats.TotalOrders, &stats.TotalRevenue, &stats.PendingOrders)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *OrderRepository) GetOrders(ctx context.Context, limit, offset int) ([]*models.Order, error) {
	const query = `
		SELECT o.order_id, o.order_number, o.client_id, o.order_status, o.subtotal_amount, 
               o.discount_amount, o.shipping_amount, o.tax_amount, o.total_amount, 
               o.discount_id, o.discount_code, o.shipping_address_id, o.ebuy_store_id,
               o.payment_method, o.payment_status, o.payment_reference, o.tracking_number,
               o.shipping_carrier, o.order_date, o.confirmed_at, o.shipped_at, o.delivered_at,
               o.cancelled_at, o.customer_notes, o.admin_notes, COALESCE(o.contact_phone, ''),
			   (SELECT proof_path FROM payment_proofs WHERE order_id = o.order_id ORDER BY created_at DESC LIMIT 1) as payment_proof,
			   c.username, c.phone, s.store_name
		FROM orders o
		JOIN client c ON o.client_id = c.client_id
		LEFT JOIN ebuy_store s ON o.ebuy_store_id = s.store_id
		ORDER BY o.order_date DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.OrderID, &o.OrderNumber, &o.ClientID, &o.OrderStatus, &o.SubtotalAmount,
            &o.DiscountAmount, &o.ShippingAmount, &o.TaxAmount, &o.TotalAmount,
            &o.DiscountID, &o.DiscountCode, &o.ShippingAddressID, &o.EbuyStoreID,
            &o.PaymentMethod, &o.PaymentStatus, &o.PaymentReference, &o.TrackingNumber,
            &o.ShippingCarrier, &o.OrderDate, &o.ConfirmedAt, &o.ShippedAt, &o.DeliveredAt,
            &o.CancelledAt, &o.CustomerNotes, &o.AdminNotes, &o.ContactPhone, &o.PaymentProof,
			&o.ClientName, &o.ClientPhone, &o.EbuyStoreName,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, rows.Err()
}

func (r *OrderRepository) GetOrderItems(ctx context.Context, orderID int64) ([]models.OrderItem, error) {
	const query = `
		SELECT oi.order_item_id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, 
               oi.discount_amount, oi.total_price, oi.product_name, oi.product_type, oi.product_sku,
               oi.size_type, oi.is_free_item, oi.parent_discount_id,
			   (SELECT image_path FROM product_images WHERE product_id = oi.product_id AND is_primary = true LIMIT 1) as product_image
		FROM order_items oi
		WHERE oi.order_id = $1`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var i models.OrderItem
		if err := rows.Scan(
			&i.OrderItemID, &i.OrderID, &i.ProductID, &i.Quantity, &i.UnitPrice,
            &i.DiscountAmount, &i.TotalPrice, &i.ProductName, &i.ProductType, &i.ProductSKU,
            &i.SizeType, &i.IsFreeItem, &i.ParentDiscountID, &i.ProductImage,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID int64, status models.OrderStatus) error {
	var query string
	switch status {
	case models.OrderStatusConfirmed:
		query = `UPDATE orders SET order_status = $1, confirmed_at = COALESCE(confirmed_at, NOW()) WHERE order_id = $2`
	case models.OrderStatusShipped:
		query = `UPDATE orders SET order_status = $1, shipped_at = COALESCE(shipped_at, NOW()) WHERE order_id = $2`
	case models.OrderStatusDelivered:
		query = `UPDATE orders SET order_status = $1, delivered_at = COALESCE(delivered_at, NOW()) WHERE order_id = $2`
	case models.OrderStatusCancelled:
		query = `UPDATE orders SET order_status = $1, cancelled_at = COALESCE(cancelled_at, NOW()) WHERE order_id = $2`
	default:
		query = `UPDATE orders SET order_status = $1 WHERE order_id = $2`
	}

	_, err := r.db.Exec(ctx, query, status, orderID)
	return err
}
