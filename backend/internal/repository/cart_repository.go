package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ryangel/ryangel-backend/internal/models"
)

// CartRepository handles database operations for shopping cart.
type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

// GetCartByID retrieves a cart by ID.
func (r *CartRepository) GetCartByID(ctx context.Context, cartID int64) (*models.Cart, error) {
	query := `
		SELECT cart_id, client_id, discount_id, created_at, updated_at
		FROM cart
		WHERE cart_id = $1
	`
	var cart models.Cart
	var pgClientID, pgDiscountID pgtype.Int8
	err := r.db.QueryRow(ctx, query, cartID).Scan(&cart.CartID, &pgClientID, &pgDiscountID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("cart not found")
		}
		return nil, err
	}
	if pgClientID.Valid {
		cart.ClientID = &pgClientID.Int64
	}
	if pgDiscountID.Valid {
		cart.DiscountID = &pgDiscountID.Int64
	}
	return &cart, nil
}

// GetCartByClientID retrieves a cart by client ID.
func (r *CartRepository) GetCartByClientID(ctx context.Context, clientID int64) (*models.Cart, error) {
	query := `
		SELECT cart_id, client_id, discount_id, created_at, updated_at
		FROM cart
		WHERE client_id = $1
	`
	var cart models.Cart
	var pgClientID, pgDiscountID pgtype.Int8
	err := r.db.QueryRow(ctx, query, clientID).Scan(&cart.CartID, &pgClientID, &pgDiscountID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No cart for this client
		}
		return nil, err
	}
	if pgClientID.Valid {
		cart.ClientID = &pgClientID.Int64
	}
	if pgDiscountID.Valid {
		cart.DiscountID = &pgDiscountID.Int64
	}
	return &cart, nil
}

// CreateCart creates a new cart.
func (r *CartRepository) CreateCart(ctx context.Context, clientID *int64) (*models.Cart, error) {
	var pgClientID pgtype.Int8
	if clientID != nil {
		pgClientID.Int64 = *clientID
		pgClientID.Valid = true
	} else {
		pgClientID.Valid = false
	}
	query := `
		INSERT INTO cart (client_id)
		VALUES ($1)
		RETURNING cart_id, client_id, discount_id, created_at, updated_at
	`
	var cart models.Cart
	var pgClientIDResult, pgDiscountIDResult pgtype.Int8
	err := r.db.QueryRow(ctx, query, pgClientID).Scan(&cart.CartID, &pgClientIDResult, &pgDiscountIDResult, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if pgClientIDResult.Valid {
		cart.ClientID = &pgClientIDResult.Int64
	}
	if pgDiscountIDResult.Valid {
		cart.DiscountID = &pgDiscountIDResult.Int64
	}
	return &cart, nil
}

// UpdateCartClientID associates a cart with a client (for login).
func (r *CartRepository) UpdateCartClientID(ctx context.Context, cartID, clientID int64) error {
	query := `
		UPDATE cart
		SET client_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE cart_id = $1
	`
	_, err := r.db.Exec(ctx, query, cartID, clientID)
	return err
}

// ApplyDiscountToCart applies a discount to the cart.
func (r *CartRepository) ApplyDiscountToCart(ctx context.Context, cartID, discountID int64) error {
	query := `
		UPDATE cart
		SET discount_id = $2, updated_at = CURRENT_TIMESTAMP
		WHERE cart_id = $1
	`
	_, err := r.db.Exec(ctx, query, cartID, discountID)
	return err
}

// RemoveDiscountFromCart removes the discount from the cart.
func (r *CartRepository) RemoveDiscountFromCart(ctx context.Context, cartID int64) error {
	query := `
		UPDATE cart
		SET discount_id = NULL, updated_at = CURRENT_TIMESTAMP
		WHERE cart_id = $1
	`
	_, err := r.db.Exec(ctx, query, cartID)
	return err
}

// GetCartItems retrieves items in a cart.
func (r *CartRepository) GetCartItems(ctx context.Context, cartID int64) ([]models.CartItemResponse, error) {
	query := `
		SELECT ci.product_id, ci.size_type, ci.quantity, ci.added_at,
		       p.product_name, p.price, p.quantity as stock_quantity
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.product_id
		WHERE ci.cart_id = $1
		ORDER BY ci.added_at
	`
	rows, err := r.db.Query(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItemResponse
	for rows.Next() {
		var item models.CartItemResponse
		var sizeType *models.SizeType
		err := rows.Scan(&item.ProductID, &sizeType, &item.Quantity, &item.AddedAt,
			&item.ProductName, &item.UnitPrice, &item.StockQuantity)
		if err != nil {
			return nil, err
		}
		item.SizeType = sizeType
		items = append(items, item)
	}
	return items, rows.Err()
}

// AddItemToCart adds or updates an item in the cart.
func (r *CartRepository) AddItemToCart(ctx context.Context, cartID, productID int64, sizeType *models.SizeType, quantity int) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, size_type, quantity)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (cart_id, product_id, size_type)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, added_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(ctx, query, cartID, productID, sizeType, quantity)
	return err
}

// UpdateCartItem updates the quantity of an item.
func (r *CartRepository) UpdateCartItem(ctx context.Context, cartItemID int64, quantity int) error {
	if quantity <= 0 {
		return r.RemoveCartItem(ctx, cartItemID)
	}
	query := `
		UPDATE cart_items
		SET quantity = $2, added_at = CURRENT_TIMESTAMP
		WHERE cart_item_id = $1
	`
	_, err := r.db.Exec(ctx, query, cartItemID, quantity)
	return err
}

// RemoveCartItem removes an item from the cart.
func (r *CartRepository) RemoveCartItem(ctx context.Context, cartItemID int64) error {
	query := `DELETE FROM cart_items WHERE cart_item_id = $1`
	_, err := r.db.Exec(ctx, query, cartItemID)
	return err
}

// ClearCart removes all items from a cart.
func (r *CartRepository) ClearCart(ctx context.Context, cartID int64) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Exec(ctx, query, cartID)
	return err
}

// DeleteCart deletes a cart.
func (r *CartRepository) DeleteCart(ctx context.Context, cartID int64) error {
	query := `DELETE FROM cart WHERE cart_id = $1`
	_, err := r.db.Exec(ctx, query, cartID)
	return err
}