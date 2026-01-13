package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ryangel/ryangel-backend/internal/models"
)

// ProductRepository handles database operations for products.
type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{db: db}
}

// ProductFilters represents filters for product listing.
type ProductFilters struct {
	Query        string
	CategoryID   *int64
	ProductType  *models.ProductType
	IsFeatured   *bool
	PriceMin     *float64
	PriceMax     *float64
	SKU          string
	IsActive     *bool
}

// ProductSort represents sorting options.
type ProductSort struct {
	Field string
	Order string // "asc" or "desc"
}

// ListProducts retrieves products with pagination and filters.
func (r *ProductRepository) ListProducts(ctx context.Context, filters ProductFilters, sort ProductSort, page, pageSize int) ([]models.ProductWithDetails, int, error) {
	offset := (page - 1) * pageSize

	// Build WHERE clause
	whereParts := []string{"p.is_active = true"} // Default to active products
	args := []interface{}{}
	argCount := 1

	if filters.Query != "" {
		whereParts = append(whereParts, fmt.Sprintf("(p.product_name ILIKE $%d OR p.product_description ILIKE $%d)", argCount, argCount+1))
		likeQuery := "%" + filters.Query + "%"
		args = append(args, likeQuery, likeQuery)
		argCount += 2
	}

	if filters.ProductType != nil {
		whereParts = append(whereParts, fmt.Sprintf("p.product_type = $%d", argCount))
		args = append(args, *filters.ProductType)
		argCount++
	}

	if filters.IsFeatured != nil {
		whereParts = append(whereParts, fmt.Sprintf("p.is_featured = $%d", argCount))
		args = append(args, *filters.IsFeatured)
		argCount++
	}

	whereClause := strings.Join(whereParts, " AND ")

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.product_id)
		FROM products p
		WHERE %s`, whereClause)

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	// Data query with images
	query := fmt.Sprintf(`
		SELECT
			p.product_id, p.product_name, p.product_description, p.product_type, p.hashtag,
			p.sku, p.price, p.compare_at_price, p.quantity, p.is_featured, p.is_active,
			p.available_sizes, p.created_at, p.updated_at,
			COALESCE(pi.image_id, 0) as image_id,
			COALESCE(pi.image_path, '') as image_path,
			COALESCE(pi.thumbnail_path, '') as thumbnail_path,
			pi.alt_text,
			pi.size_type,
			COALESCE(pi.sort_order, 0) as sort_order,
			COALESCE(pi.is_primary, false) as is_primary
		FROM products p
		LEFT JOIN product_images pi ON p.product_id = pi.product_id
		WHERE %s
		ORDER BY p.created_at DESC, pi.sort_order ASC
		LIMIT $%d OFFSET $%d`,
		whereClause, argCount, argCount+1)

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	// Group products and their images
	productMap := make(map[int64]*models.ProductWithDetails)
	for rows.Next() {
		var p models.ProductWithDetails
		var imageID int64
		var imagePath string
		var thumbnailPath string
		var altText *string
		var sizeType *models.SizeType
		var sortOrder int
		var isPrimary bool
		var availableSizes []string

		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Type, &p.Hashtag,
			&p.SKU, &p.Price, &p.CompareAtPrice, &p.Quantity, &p.IsFeatured, &p.IsActive,
			&availableSizes, &p.CreatedAt, &p.UpdatedAt,
			&imageID, &imagePath, &thumbnailPath, &altText, &sizeType, &sortOrder, &isPrimary,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("scan product: %w", err)
		}

		// Convert availableSizes to SizeType slice
		p.AvailableSizes = make([]models.SizeType, len(availableSizes))
		for i, size := range availableSizes {
			p.AvailableSizes[i] = models.SizeType(size)
		}

		// Check if we've already seen this product
		if existingProduct, exists := productMap[p.ID]; exists {
			// Add image to existing product if it exists
			if imageID != 0 {
				altTextStr := ""
				if altText != nil {
					altTextStr = *altText
				}
				img := models.ProductImage{
					ID:           imageID,
					ProductID:    p.ID,
					URL:          imagePath,
					ThumbnailURL: thumbnailPath,
					AltText:      altTextStr,
					SizeType:     sizeType,
					SortOrder:    sortOrder,
					IsPrimary:    isPrimary,
				}
				existingProduct.Images = append(existingProduct.Images, img)
			}
		} else {
			// New product
			p.Categories = []models.Category{} // Empty for now
			if imageID != 0 {
				altTextStr := ""
				if altText != nil {
					altTextStr = *altText
				}
				img := models.ProductImage{
					ID:           imageID,
					ProductID:    p.ID,
					URL:          imagePath,
					ThumbnailURL: thumbnailPath,
					AltText:      altTextStr,
					SizeType:     sizeType,
					SortOrder:    sortOrder,
					IsPrimary:    isPrimary,
				}
				p.Images = []models.ProductImage{img}
			} else {
				p.Images = []models.ProductImage{}
			}
			productMap[p.ID] = &p
		}
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	// Convert map to slice
	products := make([]models.ProductWithDetails, 0, len(productMap))
	for _, product := range productMap {
		products = append(products, *product)
	}

	return products, total, nil
}

// GetProductByID retrieves a single product with details by ID.
func (r *ProductRepository) GetProductByID(ctx context.Context, productID int64) (*models.ProductWithDetails, error) {
	query := `
		SELECT
			p.product_id, p.product_name, p.product_description, p.product_type, p.hashtag,
			p.sku, p.price, p.compare_at_price, p.quantity, p.is_featured, p.is_active,
			p.available_sizes, p.created_at, p.updated_at,
			COALESCE(pi.image_id, 0) as image_id,
			COALESCE(pi.image_path, '') as image_path,
			COALESCE(pi.thumbnail_path, '') as thumbnail_path,
			pi.alt_text,
			pi.size_type,
			COALESCE(pi.sort_order, 0) as sort_order,
			COALESCE(pi.is_primary, false) as is_primary
		FROM products p
		LEFT JOIN product_images pi ON p.product_id = pi.product_id
		WHERE p.product_id = $1
		ORDER BY pi.sort_order ASC`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product: %w", err)
	}
	defer rows.Close()

	var p *models.ProductWithDetails
	var images []models.ProductImage

	for rows.Next() {
		if p == nil {
			p = &models.ProductWithDetails{}
			var imageID int64
			var imagePath string
			var thumbnailPath string
			var altText *string
			var sizeType *models.SizeType
			var sortOrder int
			var isPrimary bool
			var availableSizes string

			err := rows.Scan(
				&p.ID, &p.Name, &p.Description, &p.Type, &p.Hashtag,
				&p.SKU, &p.Price, &p.CompareAtPrice, &p.Quantity, &p.IsFeatured, &p.IsActive,
				&availableSizes, &p.CreatedAt, &p.UpdatedAt,
				&imageID, &imagePath, &thumbnailPath, &altText, &sizeType, &sortOrder, &isPrimary,
			)
			if err != nil {
				return nil, fmt.Errorf("scan product: %w", err)
			}

			// Convert availableSizes to SizeType slice
			if availableSizes != "" && availableSizes != "{}" {
				// Remove braces and split
				sizesStr := strings.Trim(availableSizes, "{}")
				if sizesStr != "" {
					sizeStrings := strings.Split(sizesStr, ",")
					p.AvailableSizes = make([]models.SizeType, len(sizeStrings))
					for i, size := range sizeStrings {
						p.AvailableSizes[i] = models.SizeType(strings.TrimSpace(size))
					}
				} else {
					p.AvailableSizes = []models.SizeType{}
				}
			} else {
				p.AvailableSizes = []models.SizeType{}
			}

			p.Categories = []models.Category{} // Empty for now

			if imageID != 0 {
				altTextStr := ""
				if altText != nil {
					altTextStr = *altText
				}
				img := models.ProductImage{
					ID:           imageID,
					ProductID:    p.ID,
					URL:          imagePath,
					ThumbnailURL: thumbnailPath,
					AltText:      altTextStr,
					SizeType:     sizeType,
					SortOrder:    sortOrder,
					IsPrimary:    isPrimary,
				}
				images = append(images, img)
			}
		} else {
			// Additional image rows
			var imageID int64
			var imagePath string
			var thumbnailPath string
			var altText *string
			var sizeType *models.SizeType
			var sortOrder int
			var isPrimary bool

			// Skip the product fields and scan only image fields
			err := rows.Scan(
				nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
				&imageID, &imagePath, &thumbnailPath, &altText, &sizeType, &sortOrder, &isPrimary,
			)
			if err != nil {
				return nil, fmt.Errorf("scan product image: %w", err)
			}

			if imageID != 0 {
				altTextStr := ""
				if altText != nil {
					altTextStr = *altText
				}
				img := models.ProductImage{
					ID:           imageID,
					ProductID:    p.ID,
					URL:          imagePath,
					ThumbnailURL: thumbnailPath,
					AltText:      altTextStr,
					SizeType:     sizeType,
					SortOrder:    sortOrder,
					IsPrimary:    isPrimary,
				}
				images = append(images, img)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	if p == nil {
		return nil, ErrNotFound
	}

	p.Images = images
	return p, nil
}

func (r *ProductRepository) GetProductImages(ctx context.Context, productID int64) ([]models.ProductImage, error) {
	query := `
		SELECT
			image_id, product_id, image_path, thumbnail_path, alt_text, sort_order, is_primary, created_at, updated_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY sort_order ASC, created_at ASC`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("query product images: %w", err)
	}
	defer rows.Close()

	var images []models.ProductImage
	for rows.Next() {
		var img models.ProductImage
		var altText *string
		var imagePath string
		var thumbnailPath *string

		err := rows.Scan(
			&img.ID, &img.ProductID, &imagePath, &thumbnailPath, &altText, &img.SortOrder, &img.IsPrimary, &img.CreatedAt, &img.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan product image: %w", err)
		}

		img.URL = imagePath
		if thumbnailPath != nil {
			img.ThumbnailURL = *thumbnailPath
		}
		if altText != nil {
			img.AltText = *altText
		} else {
			img.AltText = ""
		}

		images = append(images, img)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return images, nil
}