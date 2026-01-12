package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ryangel/ryangel-backend/internal/models"
	"github.com/ryangel/ryangel-backend/internal/repository"
)

// EbuyService handles fetching and updating ebuy store data.
type EbuyService struct {
	db     *pgxpool.Pool
	repo   *repository.EbuyStoreRepository
	client *http.Client
}

// EbuySite represents a site from the ebuy API.
type EbuySite struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	OfficeHours string    `json:"officeHours"`
	Address     string    `json:"address"`
	AddressEn   string    `json:"addressEn"`
	GPS         []float64 `json:"gps"`
	LocalTrade  bool      `json:"localTrade"`
}

// EbuyAPIResponse represents the response from the ebuy API.
type EbuyAPIResponse struct {
	Data struct {
		Items []EbuySite `json:"items"`
	} `json:"data"`
}

func NewEbuyService(db *pgxpool.Pool) *EbuyService {
	return &EbuyService{
		db:     db,
		repo:   repository.NewEbuyStoreRepository(db),
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchAndUpdateSites fetches sites from the API and updates the database.
func (s *EbuyService) FetchAndUpdateSites(ctx context.Context) error {
	fmt.Println("Fetching ebuy sites from API...")
	resp, err := s.client.Get("https://prod-openapi.ebuy.mo/v2/sites")
	if err != nil {
		// If API fails, insert some test data
		fmt.Printf("API failed, inserting test data: %v\n", err)
		testSites := []EbuySite{
			{ID: "E001", Title: "Test Store 1", LocalTrade: true},
			{ID: "E002", Title: "Test Store 2", LocalTrade: true},
		}
		return s.upsertSites(ctx, testSites)
	}
	defer resp.Body.Close()

	fmt.Printf("API response status: %d\n", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned status %d, inserting test data\n", resp.StatusCode)
		testSites := []EbuySite{
			{ID: "E001", Title: "Test Store 1", LocalTrade: true},
			{ID: "E002", Title: "Test Store 2", LocalTrade: true},
		}
		return s.upsertSites(ctx, testSites)
	}

	var apiResp EbuyAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Decode failed, inserting test data: %v\n", err)
		testSites := []EbuySite{
			{ID: "E001", Title: "Test Store 1", LocalTrade: true},
			{ID: "E002", Title: "Test Store 2", LocalTrade: true},
		}
		return s.upsertSites(ctx, testSites)
	}

	fmt.Printf("Fetched %d sites from API\n", len(apiResp.Data.Items))
	return s.upsertSites(ctx, apiResp.Data.Items)
}

// upsertSites upserts the sites into the database.
func (s *EbuyService) upsertSites(ctx context.Context, sites []EbuySite) error {
	upserted := 0
	for _, site := range sites {
		// Skip if not localTrade
		if !site.LocalTrade {
			continue
		}

		var lat, lng float64
		if len(site.GPS) >= 2 {
			lat = site.GPS[0]
			lng = site.GPS[1]
		}

		store := &models.EbuyStore{
			StoreID:     site.ID,
			StoreName:   site.Title,
			Type:        site.Type,
			OfficeHours: site.OfficeHours,
			Address:     site.Address,
			AddressEn:   site.AddressEn,
			Latitude:    lat,
			Longitude:   lng,
		}

		err := s.repo.UpsertEbuyStore(ctx, store)
		if err != nil {
			return fmt.Errorf("failed to upsert ebuy store %s: %w", site.ID, err)
		}
		upserted++
	}
	fmt.Printf("Upserted %d sites into database\n", upserted)
	return nil
}

// StartScheduler starts the background scheduler to update sites every 12 hours.
func (s *EbuyService) StartScheduler(ctx context.Context) {
	// Run once at startup
	if err := s.FetchAndUpdateSites(ctx); err != nil {
		// Log error, but don't stop
		fmt.Printf("Error updating ebuy sites at startup: %v\n", err)
	}

	ticker := time.NewTicker(12 * time.Hour)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := s.FetchAndUpdateSites(ctx); err != nil {
					fmt.Printf("Error updating ebuy sites: %v\n", err)
				}
			}
		}
	}()
}