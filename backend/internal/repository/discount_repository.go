package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ryangel/ryangel-backend/internal/models"
)

type DiscountRepository struct {
	pool *pgxpool.Pool
}

func NewDiscountRepository(pool *pgxpool.Pool) *DiscountRepository {
	return &DiscountRepository{pool: pool}
}

func (r *DiscountRepository) GetAutoApplyDiscounts(ctx context.Context) ([]models.Discount, error) {
	query := `
        SELECT 
            discount_id, discount_code, discount_name, discount_type, discount_value,
            buy_quantity, get_quantity, free_product_id, applies_to_same_product,
            minimum_order_amount, maximum_discount_amount, product_type_restriction,
            start_date, end_date, usage_limit, used_count, is_active, applies_to, is_auto_apply
        FROM discounts
        WHERE is_active = TRUE 
          AND is_auto_apply = TRUE
          AND start_date <= NOW() 
          AND end_date >= NOW()
    `
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var discounts []models.Discount
	for rows.Next() {
		var d models.Discount
		err := rows.Scan(
			&d.DiscountID, &d.DiscountCode, &d.DiscountName, &d.DiscountType, &d.DiscountValue,
			&d.BuyQuantity, &d.GetQuantity, &d.FreeProductID, &d.AppliesToSameProduct,
			&d.MinimumOrderAmount, &d.MaximumDiscountAmount, &d.ProductTypeRestriction,
			&d.StartDate, &d.EndDate, &d.UsageLimit, &d.UsedCount, &d.IsActive, &d.AppliesTo, &d.IsAutoApply,
		)
		if err != nil {
			return nil, err
		}
		discounts = append(discounts, d)
	}
	return discounts, nil
}
