package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ryangel/ryangel-backend/internal/models"
)

// EbuyStoreRepository handles database operations for ebuy stores.
type EbuyStoreRepository struct {
	db *pgxpool.Pool
}

func NewEbuyStoreRepository(db *pgxpool.Pool) *EbuyStoreRepository {
	return &EbuyStoreRepository{db: db}
}

// GetEbuyStores retrieves all ebuy stores.
func (r *EbuyStoreRepository) GetEbuyStores(ctx context.Context) ([]models.EbuyStore, error) {
	query := `SELECT store_id, store_name, type, office_hours, address, address_en, latitude, longitude FROM ebuy_store ORDER BY store_id`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stores []models.EbuyStore
	for rows.Next() {
		var store models.EbuyStore
		err := rows.Scan(&store.StoreID, &store.StoreName, &store.Type, &store.OfficeHours, &store.Address, &store.AddressEn, &store.Latitude, &store.Longitude)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}
	return stores, rows.Err()
}

// GetEbuyStoreByID retrieves a single ebuy store by ID.
func (r *EbuyStoreRepository) GetEbuyStoreByID(ctx context.Context, storeID string) (*models.EbuyStore, error) {
	query := `SELECT store_id, store_name, type, office_hours, address, address_en, latitude, longitude FROM ebuy_store WHERE store_id = $1`
	var store models.EbuyStore
	err := r.db.QueryRow(ctx, query, storeID).Scan(&store.StoreID, &store.StoreName, &store.Type, &store.OfficeHours, &store.Address, &store.AddressEn, &store.Latitude, &store.Longitude)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &store, nil
}

// UpsertEbuyStore inserts or updates an ebuy store.
func (r *EbuyStoreRepository) UpsertEbuyStore(ctx context.Context, store *models.EbuyStore) error {
	query := `
		INSERT INTO ebuy_store (store_id, store_name, type, office_hours, address, address_en, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (store_id) DO UPDATE SET
			store_name = EXCLUDED.store_name,
			type = EXCLUDED.type,
			office_hours = EXCLUDED.office_hours,
			address = EXCLUDED.address,
			address_en = EXCLUDED.address_en,
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude
	`
	_, err := r.db.Exec(ctx, query, store.StoreID, store.StoreName, store.Type, store.OfficeHours, store.Address, store.AddressEn, store.Latitude, store.Longitude)
	return err
}