package main

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Hardcoded for the script
const DB_URL = "postgresql://ryangel:RyangelPa33word@localhost:5432/ryangel?sslmode=disable"

type Order struct {
	OrderID    int64
	Subtotal   float64
	Discount   float64
	Shipping   float64
	Total      float64
}

type OrderItem struct {
	OrderID     int64
	ProductType string
	UnitPrice   float64
	Quantity    int
}

func main() {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(DB_URL)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer pool.Close()

	fmt.Println("Connected to database...")

	// 1. Fetch all orders
	rows, err := pool.Query(ctx, "SELECT order_id FROM orders")
	if err != nil {
		log.Fatalf("Failed to query orders: %v", err)
	}

	var orderIDs []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("Failed to scan order id: %v", err)
		}
		orderIDs = append(orderIDs, id)
	}
	rows.Close()

	fmt.Printf("Found %d orders to process.\n", len(orderIDs))

	for _, orderID := range orderIDs {
		processOrder(ctx, pool, orderID)
	}
}

func processOrder(ctx context.Context, db *pgxpool.Pool, orderID int64) {
	// Fetch items
	rows, err := db.Query(ctx, `
		SELECT order_id, product_type, unit_price, quantity 
		FROM order_items 
		WHERE order_id = $1`, orderID)
	if err != nil {
		log.Printf("[Order %d] Failed to fetch items: %v", orderID, err)
		return
	}
	defer rows.Close()

	var items []OrderItem
	var subtotal float64
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(&i.OrderID, &i.ProductType, &i.UnitPrice, &i.Quantity); err != nil {
			log.Printf("[Order %d] Failed to scan item: %v", orderID, err)
			return
		}
		items = append(items, i)
		subtotal += i.UnitPrice * float64(i.Quantity)
	}
	rows.Close() // Ensure closed before next query

	// --- Calculate Discount (B3G1 for Faichun) ---
	var faiChunPrices []float64
	for _, item := range items {
		if item.ProductType == "faiachun" {
			for k := 0; k < item.Quantity; k++ {
				faiChunPrices = append(faiChunPrices, item.UnitPrice)
			}
		}
	}

	discountAmount := 0.0
	// Sort prices ascending
	sort.Float64s(faiChunPrices)

	// Buy 3 Get 1 Free = Groups of 4
	// Cheapest 1 in every group of 4 is free
	numFaiChuns := len(faiChunPrices)
	numFree := numFaiChuns / 4 // e.g. 5/4 = 1, 8/4 = 2

	for i := 0; i < numFree; i++ {
		discountAmount += faiChunPrices[i]
	}

	if discountAmount > subtotal {
		discountAmount = subtotal
	}

	discountedSubtotal := subtotal - discountAmount
	if discountedSubtotal < 0 {
		discountedSubtotal = 0
	}

	// --- Calculate Shipping (Free if > $60) ---
	// Assumption: Threshold applies to Discounted Subtotal
	shippingFee := 5.0
	if discountedSubtotal >= 60.0 {
		shippingFee = 0.0
	}

	total := discountedSubtotal + shippingFee

	// Update Order
	// subtotal_amount might already be correct, but let's ensure consistency
	// The problem described was discount and total were wrong. 
	// The user log shows 100 subtotal, 0 discount, 0 shipping, 100 total.
	// We'll update subtotal (just in case), discount, shipping, total.

	_, err = db.Exec(ctx, `
		UPDATE orders 
		SET subtotal_amount = $1, discount_amount = $2, shipping_amount = $3, total_amount = $4 
		WHERE order_id = $5`, 
		subtotal, discountAmount, shippingFee, total, orderID)
	
	if err != nil {
		log.Printf("[Order %d] Failed to update: %v", orderID, err)
	} else {
		fmt.Printf("[Order %d] Updated: Sub=%.2f, Disc=%.2f, Ship=%.2f, Total=%.2f\n", 
			orderID, subtotal, discountAmount, shippingFee, total)
	}
}
