package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryangel/ryangel-backend/internal/models"
)

type CreateOrderParams struct {
	ClientID    int64
	EbuyStoreID string
	Name        string
	Phone       string // Just for verification
	Email       string
	Instagram   string
	ProofPath   string
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
	err = tx.QueryRow(ctx, "SELECT cart_id FROM cart WHERE client_id = $1", params.ClientID).Scan(&cartID)
	if err != nil {
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

	var orderID int64
	var orderDate time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO orders (
			order_number, client_id, order_status, 
			subtotal_amount, discount_amount, shipping_amount, tax_amount, total_amount,
			ebuy_store_id, payment_method, payment_status, customer_notes, order_date
		) VALUES (
			$1, $2, 'pending', 
			$3, 0, 0, 0, $3, 
			$4, 'mpay', 'pending', $5, NOW()
		) RETURNING order_id, order_date`,
		orderNum, params.ClientID, subtotal, params.EbuyStoreID, customerNotes,
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
               o.cancelled_at, o.customer_notes, o.admin_notes,
			   (SELECT proof_path FROM payment_proofs WHERE order_id = o.order_id ORDER BY created_at DESC LIMIT 1) as payment_proof
		FROM orders o
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
            &o.CancelledAt, &o.CustomerNotes, &o.AdminNotes, &o.PaymentProof,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}
	return orders, rows.Err()
}
